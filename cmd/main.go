/*
Copyright paskal.maksim@gmail.com
Licensed under the Apache License, Version 2.0 (the "License")
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	gitVersion = "dev"
	kubeconfig = flag.String("kubeconfig", "", "kubernetes config")
	namespace  = flag.String("namespace", os.Getenv("MY_POD_NAMESPACE"), "pod namespace")
	pod        = flag.String("pod", os.Getenv("HOSTNAME"), "pod name")
	podLabel   = flag.String("podLabel", "nodeZone", "pod label to set")
	label      = flag.String("label", "topology.kubernetes.io/zone", "node label")
	logFlags   = flag.Int("logFlags", 0, "log flags")
	version    = flag.Bool("version", false, "version")
)

func main() {
	flag.Parse()

	log.SetFlags(*logFlags)

	if *version {
		log.Println(gitVersion)

		return
	}

	if len(*pod) == 0 {
		log.Fatal("no pod")
	}

	if len(*namespace) == 0 {
		log.Fatal("no namespace")
	}

	ctx := context.Background()

	var (
		err    error
		config *rest.Config
	)

	if len(*kubeconfig) > 0 {
		config, err = clientcmd.BuildConfigFromFlags("", *kubeconfig)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		config, err = rest.InClusterConfig()
		if err != nil {
			log.Fatal(err)
		}
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}

	podInfo, err := clientset.CoreV1().Pods(*namespace).Get(ctx, *pod, metav1.GetOptions{})
	if err != nil {
		log.Fatal(err)
	}

	nodeInfo, err := clientset.CoreV1().Nodes().Get(ctx, podInfo.Spec.NodeName, metav1.GetOptions{})
	if err != nil {
		log.Fatal(err)
	}

	labelPatch := fmt.Sprintf(`[{"op":"add","path":"/metadata/labels/%s","value":"%s" }]`,
		*podLabel,
		nodeInfo.Labels[*label],
	)

	_, err = clientset.
		CoreV1().
		Pods(podInfo.Namespace).
		Patch(ctx, podInfo.Name, types.JSONPatchType, []byte(labelPatch), metav1.PatchOptions{})

	if err != nil {
		log.Fatal(err)
	}

	log.Println("OK")
}

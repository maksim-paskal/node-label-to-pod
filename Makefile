test:
	./scripts/validate-license.sh
	go fmt ./cmd/
	go mod tidy
	go test -race ./cmd/
	golangci-lint run -v
	kubectl apply --dry-run=client --validate -f deployment.yaml
run:
	GOFLAGS="-trimpath" go run -v -race ./cmd -kubeconfig=kubeconfig -namespace=default $(args)
build:
	docker build . -t paskalmaksim/node-label-to-pod:dev
push:
	docker push paskalmaksim/node-label-to-pod:dev
testK8s:
	kubectl apply -f deployment.yaml
clean:
	kubectl delete -f deployment.yaml
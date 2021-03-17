test:
	./scripts/validate-license.sh
	go fmt ./cmd/
	go mod tidy
	go test -race ./cmd/
	golangci-lint run -v
run:
	GOFLAGS="-trimpath" go run -v -race ./cmd -kubeconfig=kubeconfig -namespace=default $(args)
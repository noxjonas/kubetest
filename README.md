# kind-test

Helper package to run e2e tests in go test framework.

## install

```shell
go install sigs.k8s.io/kind@v0.20.0

go mod download

go test ./... -v

kind create cluster --config ./cluster/kind.yaml --wait 30s
```
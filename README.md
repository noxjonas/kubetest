# kubetest

Helper package to run kubernetes e2e tests using go testing package.


## test

```shell
go install sigs.k8s.io/kind@v0.20.0

go mod download

kind create cluster --config ./cluster/kind.yaml --wait 30s

go test ./... -v

```
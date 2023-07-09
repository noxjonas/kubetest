# kind-test

## install

```shell
go install sigs.k8s.io/kind@v0.20.0

kind create cluster --config ./cluster/kind.yaml --wait 30s

```
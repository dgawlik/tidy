# Tidy

## What it is

When you are busy you want to reduce time speding reading kubernetes manifests. Unfortunately kubectl commands return manifests bloated with a lot of uimportant stuff. This tool does 2 things:

* removes nil/0/""/[]/{} values
* removes the defaults for fields in spec

I haven't tested it extensively but if everything works it should cut the manifest size in half.

## Installation
```bash
go build
```

Then add binary to the path.

Or download here:

https://github.com/dgawlik/tidy/actions/runs/10820426310

## Usage

```bash
kubectl get <resource> <name> -o json | tidy
```

Example:

```bash
kubectl get pod etcd-minikube -n kube-system -o json |  tidy > tidy.json
```

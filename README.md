# Provider Artifactory

`provider-jfrogartifactory` is a [Crossplane](https://crossplane.io/) provider that
is built using [Upjet](https://github.com/crossplane/upjet) code
generation tools and exposes XRM-conformant managed resources for the
Artifactory API.

## Getting Started

Install the provider by using the following command after changing the image tag
to the [latest release](https://marketplace.upbound.io/providers/myorg/provider-jfrogartifactory):
```
up ctp provider install myorg/provider-jfrogartifactory:v0.1.0
```

Alternatively, you can use declarative installation:
```
cat <<EOF | kubectl apply -f -
apiVersion: pkg.crossplane.io/v1
kind: Provider
metadata:
  name: provider-jfrogartifactory
spec:
  package: myorg/provider-jfrogartifactory:v0.1.0
EOF
```

Notice that in this example Provider resource is referencing ControllerConfig with debug enabled.

You can see the API reference [here](https://doc.crds.dev/github.com/myorg/provider-jfrogartifactory).

## Developing

Run code-generation pipeline:
```console
go run cmd/generator/main.go "$PWD"
```

Run against a Kubernetes cluster:

```console
make run
```

Build, push, and install:

```console
make all
```

Build binary:

```console
make build
```

## Report a Bug

For filing bugs, suggesting improvements, or requesting new features, please
open an [issue](https://github.com/myorg/provider-jfrogartifactory/issues).



# Testing

Cannot use OSS Artifactory because it does not support creating repositories through REST APIs.
Get a license for Artifactory: https://jfrog.com/start-free/#ft
Set the environment variable `ARTIFACTORY_LICENSE_KEY` in your local ~/.zshrc and restart your IDE.

in a terminal in the dev container:
```console
mage setupE2E
kubectl port-forward --namespace jfrog svc/artifactory-artifactory-nginx 8888:80
```

in new terminal:
```console
kubectl apply -f package/crds
make run
```

in new terminal:
```console
kubectl apply -f e2e/providerconfig.yaml
#kubectl apply -f examples/manifests/genericrepository.yaml
mage testE2E
```


# Manual testing by applying resources
## Steps to use this provider artifactory
- In a cluster ,Apply all the manifest files in package/crds
In one of you terminals checkout this repository https://github.com/suvaanshkumar/provider-jfrogartifactory amd run `make run`

In another terminal run the following

- Generate an identity token on artifactory to be used here.
- Create a file similar to creds.json present in examples/manifests/template folder and fill in the url and the and key
- Base64 encode this file and put it in the secret.yaml  file present in examples/manifests/template in data field and apply the secret
- Apply providerconfigartifactory.yaml present in examples/manifests
- Apply the genericrepository.yaml or any other resource you want
# Provider Artifactory

`provider-jfrogartifactory` is a [Crossplane](https://crossplane.io/) provider that
is built using [Upjet](https://github.com/crossplane/upjet) code
generation tools and exposes XRM-conformant managed resources for the
Artifactory API.

## Getting Started

Install the provider by using the following command after changing the image tag
to the [latest release](https://marketplace.upbound.io/providers/guidewire-oss/provider-jfrogartifactory):
```
up ctp provider install guidewire-oss/provider-jfrogartifactory:v0.1.0
```

Alternatively, you can use declarative installation:
```
cat <<EOF | kubectl apply -f -
apiVersion: pkg.crossplane.io/v1
kind: Provider
metadata:
  name: provider-jfrogartifactory
spec:
  package: guidewire-oss/provider-jfrogartifactory:v0.1.0
EOF
```

Notice that in this example Provider resource is referencing ControllerConfig with debug enabled.

You can see the API reference [here](https://doc.crds.dev/github.com/guidewire-oss/provider-jfrogartifactory).

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
open an [issue](https://github.com/guidewire-oss/provider-jfrogartifactory/issues).

# Running e2e tests using make run (using dev edge nodes)

in a terminal in the dev container:
```console
mage setupE2E
```

in new terminal:
```console
kubectl apply -f package/crds
make run
```

in new terminal:
```console
mage testE2E
```

# Get temporary artifactory license (alternative to using dev edge nodes for e2e testing)
- Cannot use OSS Artifactory because it does not support creating repositories through REST APIs.
- Get a license for Artifactory: https://jfrog.com/start-free/#ft
- Set the environment variable `ARTIFACTORY_LICENSE_KEY` in your local ~/.zshrc and restart your IDE.
- After running ```mage setupE2E```, run ```kubectl port-forward --namespace jfrog svc/artifactory-artifactory-nginx 8888:80```

# Testing using the provider
Note that ```...crossplane.yaml: No such file or directory``` can be ignored, and ```make build.all``` can be used to build the image for amd64 on arm64 machines
```console
make build
```
```console
up xpkg build \
--controller <IMAGE_NAME>:<IMAGE_TAG> \
--package-root ./package \
--output ./jfrogprovider.xpkg
```
Push the image to ECR:
```console
up xpkg push <ACCOUNT_ID>.dkr.ecr.us-west-2.amazonaws.com/jfrogprovider:<IMAGE_TAG> -f jfrogprovider.xpkg
```

To test on a scratch cluster using provider , do the following:

Create the following ```provider-artifactory.yaml``` on your local machine:
```console
apiVersion: pkg.crossplane.io/v1
kind: Provider
metadata:
    name: provider-artifactory
spec:
    package: <ACCOUNT_ID>.dkr.ecr.us-west-2.amazonaws.com/jfrogprovider:<IMAGE_TAG>
```
Then apply it on the scratch cluster:
```console
kubectl apply -f <FILE_PATH>/provider-artifactory.yaml
```
Apply the kubernetes secrets
```console
kubectl apply -f examples/manifests/templates/<FILE_NAME>.yaml
```
Apply the provider configs and repositories
```console
kubectl apply -f examples/manifests/<FILE_NAME>.yaml
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

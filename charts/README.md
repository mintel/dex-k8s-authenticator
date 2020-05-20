# Helm charts for installing 'dex' with 'dex-k8s-authenticator'

The charts in this folder install [`dex`](https://github.com/coreos/dex)
with [`dex-k8s-authenticator`](https://github.com/mintel/dex-k8s-authenticator)

`dex-k8s-authenticator` is a helper application for `dex`. `dex` lets you use external 
Identify Providers (like Google, Microsoft, GitHub, LDAP) to authenticate access to Kubernetes cluster
(e.g. for `kubectl`). This helper makes it easy to provide a web UI for one or more clusters.
It gives users the information and commands to configure `kubectl` to work with the credentials `dex` provides.

Each install of `dex` and/or `dex-k8s-authenticator` can support multiple Kubernetes clusters.
So you can install one of each for all your clusters, one in each cluster, or any combination.

```
git clone https://github.com/mintel/dex-k8s-authenticator.git
helm inspect values charts/dex > dex.yaml
helm inspect values charts/dex-k8s-authenticator > dex-k8s-authenticator.yaml
```
Edit the values files for your environment and requirements (`dex.yaml` and `dex-k8s-authenticator.yaml`). 

Create the DNS names for your `dex` (e.g. 'dex.example.com') and `dex-k8s-authenticator` (e.g. 'login.example.com')
pointed at the ingress controller you are using. Be sure to enable HTTPS. You can install 
[`cert-manager`](https://github.com/jetstack/cert-manager) to automatically issue Lets Encrypt certificates.

You also need to configure each Kubernetes cluster to use `dex` at e.g. 'dex.example.com' by [setting the OIDC parameters
for the Kubernetes API server](https://kubernetes.io/docs/admin/authentication/#openid-connect-tokens). 
This is easy using [`kube-aws` installer](https://github.com/kubernetes-incubator/kube-aws/tree/master/contrib/dex).

```
helm install --namespace dex --values dex.yaml charts/dex
helm install --namespace dex --values dex-k8s-authenticator charts/dex-k8s-authenticator
```
Navigate to https://login.example.com and follow the instructions to authenticate using `dex` and configure `kubectl`.

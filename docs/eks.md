# Running on AWS EKS

Running `dex-k8s-authenticator` on AWS EKS is a bit tricky. The reason is, AWS EKS does not allow you to set custom API server flags. [1] This is required to set DEX as an endpoint for OIDC provider for kubernetes API server.

The good news is, you can still get DEX working properly with EKS, however, this will make your auth-flow a bit more complicated than the one on not-managed regular kubernetes setups. This will also introduce an additional service to your auth-flow. kube-oidc-proxy. [2]

Here, I'll provide 2 diagrams to show how end-to-end auth works with `dex-k8s-authenticator` on both regular kubernetes clusters and on managed kubernetes clusters.

**Auth flow on self-setup kubernetes cluster**
![Authentication on a self-setup cluster](/docs/images/auth-regular.jpg)

**Auth flow on managed kubernetes cluster**
![Authentication on a managed cluster](/docs/images/auth-managed.jpg)


### How to Setup

Thankfully, setting all these up is not as difficult as it seems. Mostly because, all the required services have a proper helm-chart. You just need to pass correct values to these charts, depending on your environment.

Required Charts
- nginx-ingress-controller - https://github.com/helm/charts/tree/master/stable/nginx-ingress
- dex - https://github.com/helm/charts/tree/master/stable/dex
- kube-oidc-proxy - https://github.com/jetstack/kube-oidc-proxy/tree/master/deploy/charts/kube-oidc-proxy
- dex-k8s-authenticator - https://github.com/sl1pm4t/dex-k8s-authenticator/tree/master/charts

You should also setup a DNS record that points to your nginx controller (load-balancer), and setup an AWS ACM certificate.

**Dex Setup**

- Enable ingress on the chart values. Set your ingress host as "dex.example.com"
- Setup your connectors. (actual identity providers)
- Setup a static client, for redirect urls, provide "login.example.com"

**kube-oidc-proxy Setup**

- Enable ingress on the chart values. Set your ingress host as "oidc-proxy.example.com"
- Setup OIDC configuration. For issuer URL, provide ""dex.example.com""

**dex-k8s-authenticator Setup**

- Enable ingress on the chart values. Set your ingress host as "login.example.com"
- Setup your clusters.
  - for issuer; provide ""dex.example.com""
  - for k8s-master-uri; provide "oidc-proxy.example.com", instead of actual k8s-master-uri.
  - provide your client_secret as you set it up on your static_client during Dex setup.
  - for k8s-ca-pem; you need to provide the certificate for the ELB that serves "oidc-proxy.example.com", in PEM format.

### References

- [1] https://github.com/aws/containers-roadmap/issues/166
- [2] https://github.com/jetstack/kube-oidc-proxy

# Helm chart for dex-k8s-authenticator

This chart installs [`dex-k8s-authenticator`](https://github.com/mintel/dex-k8s-authenticator) in a Kubernetes cluster. 
`dex-k8s-authenticator` is a helper application for [`dex`](https://github.com/coreos/dex). `dex` lets you use external 
Identify Providers (like Google, Microsoft, GitHUb, LDAP) to authenticate access to Kubernetes cluster
(e.g. for `kubectl`). This helper makes it easy to provide a web UI for one or more clusters.
It give uses the information and commands to configure `kubectl` to work with the credentials `dex` provides.

You can configure one or more clusters in this chart configuration, for the one or more `dex` installs you may have.

```yaml
# Default values for dex-k8s-authenticator.

# Deploy environment label, e.g. dev, test, prod
global:
  deployEnv: dev

replicaCount: 1

image:
  repository: mintel/dex-k8s-authenticator
  tag: latest
  pullPolicy: Always

dexK8sAuthenticator:
  port: 5555
  debug: false
  web_path_prefix: /
  #logoUrl: http://<path-to-your-logo.png>
  #kubectl_version: v1.11.2
  #tlsCert: /path/to/dex-client.crt
  #tlsKey: /path/to/dex-client.key
  clusters:
  - name: my-cluster
    short_description: "My Cluster"
    description: "Example Cluster Long Description..."
    client_secret: pUBnBOY80SnXgjibTYM9ZWNzY2xreNGQok
    issuer: https://dex.example.com
    k8s_master_uri: https://my-cluster.example.com
    client_id: my-cluster
    redirect_uri: https://login.example.com/callback/my-cluster
    k8s_ca_uri: https://url-to-your-ca.crt

service:
  type: ClusterIP
  port: 5555

ingress:
  enabled: false
  annotations: {}
    # kubernetes.io/ingress.class: nginx
    # kubernetes.io/tls-acme: "true"
  path: /
  hosts:
    - chart-example.local
  tls: []
  #  - secretName: chart-example-tls
  #    hosts:
  #      - chart-example.local

resources: {}
  # We usually recommend not to specify default resources and to leave this as a conscious
  # choice for the user. This also increases chances charts run on environments with little
  # resources, such as Minikube. If you do want to specify resources, uncomment the following
  # lines, adjust them as necessary, and remove the curly braces after 'resources:'.
  # limits:
  #  cpu: 100m
  #  memory: 128Mi
  # requests:
  #  cpu: 100m
  #  memory: 128Mi

caCerts: 
  enabled: false
  secrets: {}
  # Array of Self Signed Certificates
  # cat CA.crt | base64 -w 0
  #
  #     name: The internal k8s name of the secret we create. It's also used in 
  #     the volumeMount name. It must respect the k8s naming convension (avoid 
  #     upper-case and '.' to be safe).
  #
  #     filename: The filename of the CA to be mounted. It must end in .crt for
  #     update-ca-certificates to work
  #
  #     value: The base64 encoded value of the CA
  #
  #secrets:
  #- name: ca-cert1
  #  filename: ca1.crt
  #  value: LS0tLS1......X2F
  #- name: ca-cert2
  #  filename: ca2.crt
  #  value: DS1tFA1......X2F


nodeSelector: {}

tolerations: []

affinity: {}
```

## SSL

Additional reading [here](./ssl.md)

### Adding Trusted Certs

Multiple trusted certs can be added using the `caCerts` option. Make sure to `base64` encode each CA.

### Service Requests on SSL

Define the filepath to your cert and key using the following options in your helm chart.
  - `dexK8sAuthenticator.tlsCert`
  - `dexK8sAuthenticator.tlsKey`

TODO: Requires more work as we don't have a way deploy self-signed certs or change to https scheme.
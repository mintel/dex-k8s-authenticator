# Configuration 

An example configuration is available [here](../examples/config.yaml)

## Configuration options


| Name              | Required | Context | Description                                                                           |
|-------------------|----------|---------|---------------------------------------------------------------------------------------|
| name              | yes      | cluster | Internal id of cluster                                                                |
| short_description | yes      | cluster | Short description of cluster                                                |
| description       | yes      | cluster | Extended description of cluster                                             |
| client_secret     | yes      | cluster | OAuth2 client-secret (shared between dex-k8s-auth and dex)                            |
| client_id         | yes      | cluster | OAuth2 client-id public identifier (shared between dex-k8s-auth and dex)              |
| issuer            | yes      | cluster | Dex issuer url                                                                        |
| redirect_uri      | yes      | cluster | Redirect uri called by dex (defines a callback on dex-k8s-auth)                       |
| k8s_master_uri    | no       | cluster | Kubernetes api-server endpoint (used in kubeconfig)                                   |
| k8s_ca_uri        | no       | cluster | A url pointing to the CA for generating 'certificate-authority' option in kubeconfig  |
| k8s_ca_pem        | no       | cluster | The CA for your k8s server (used in generating instructions)                          |
| tls_cert          | no       | root    | Path to TLS cert if SSL enabled                                                       |
| tls_key           | no       | root    | Path to TLS key if SSL enabled                                                        |
| idp_ca_uri        | no       | root    | A url pointing to the CA for generating 'idp-certificate-authority' in the kubeconfig |
| trusted_root_ca   | no       | root    | A list of trusted-root CA's to be loaded by dex-k8s-auth at runtime                   |
| listen            | yes      | root    | The listen address/port                                                                    |
| web_path_prefix   | no       | root    | A path-prefix to serve dex-k8s-auth at (defaults to '/')                              |
| kubectl_version   | no       | root    | A kubectl-version string that is used to provided a download path                     |
| logo_uri          | no       | root    | A url pointing to a logo image that is displayed in the header                        |
| debug             | no       | root    | Enable more debug by setting to true                                                  |

## Environment Variable Support

Environment variables can be included in the configuration. This enables you to pass in dynamic variables or secrets without having to hardcode them.
```
clusters:
  - name: example-cluster
    short_description: "Example Cluster"
    description: "Example Cluster Long Description..."
    client_secret: ${CLIENT_SECRET_EXAMPLE_CLUSTER}
```

In this example, `${CLIENT_SECRET_EXAMPLE_CLUSTER}` is substituted at runtime. Must be enclosed by `${...}`

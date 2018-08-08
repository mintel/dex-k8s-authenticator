# SSL

## Adding Trusted Certificates

There are a couple of ways to have `dex-k8s-authenticator` use trusted certificates.

### 1. Mounting `/certs/` volume

`entrypoint.sh` runs `update-ca-certificates` against certificates found in `/certs/`

They must end in the extension `.crt`

If using docker, you can mount a volume like so:

```
docker run --rm -t -i \
    -v /tmp/certs:/certs:ro \
    -v /tmp/config.yml:/tmp/config.yml:ro \ 
    mintel/dex-k8s-authenticator:latest --config /tmp/config.yml
```

### 2. `trusted_root_ca` config option

You can define multiple certificates via the configuration file:

```yaml
trusted_root_ca:
  - |
    -----BEGIN CERTIFICATE-----
    MIIGJDCCBAygAwI...
    -----END CERTIFICATE-----
```

## Serving requests on SSL

The configuration file requires the following:


```yaml
listen: https://127.0.0.1:5555
tls_cert: /path/to/dex-client.crt
tls_key: /path/to/dex-client.key
```

- Note, the `listen` option is using `https` not `http`
- You need to supply both `.crt` and the `.key` files

The `.crt` and `.key` file can be mounted as a volume.

## Helm Charts

### dex-k8a-authenticator

Our Helm chart providers options for both using trusted root certs, and serving requests on SSL.

For more information on SSL support in our Helm chart, please read [here](../charts/dex-k8s-authenticator)
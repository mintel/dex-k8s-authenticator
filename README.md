# Dex K8s Authenticator

[![golang-lint](https://github.com/mintel/dex-k8s-authenticator/workflows/golangci-lint/badge.svg)](https://github.com/mintel/dex-k8s-authenticator/actions/workflows/golangci-lint.yml)
[![k8s-lint](https://github.com/mintel/dex-k8s-authenticator/workflows/k8s-lint/badge.svg)](https://github.com/mintel/dex-k8s-authenticator/actions/workflows/k8s.yml)

A helper web-app which talks to one or more [Dex Identity services](https://github.com/dexidp/dex) to generate
`kubectl` commands for creating and modifying a `kubeconfig`.

The Web UI supports generating tokens against multiple cluster such as Dev / Staging / Production. 


## Also provides
* Helm Charts
* SSL Support
* Linux/Mac/Windows instructions

## Documentation

- [Developing and Running](docs/develop.md)
- [Configuration Options](docs/config.md)
- [Using the Helm Charts](docs/helm.md)
- [SSL Support](docs/ssl.md)

## Screen shots

![Index Page](examples/index-page.png)

![Kubeconfig Page](examples/kubeconfig-page.png)


## Contributing

Feel free to raise feature-requests and bugs. PR's are also very welcome.

## Alternatives

- https://github.com/heptiolabs/gangway
- https://github.com/micahhausler/k8s-oidc-helper
- https://github.com/negz/kuberos
- https://github.com/negz/kubehook
- https://github.com/fydrah/loginapp

This application is based on the original [example-app](https://github.com/coreos/dex/tree/master/cmd/example-app
) available in the CoreOS Dex repo.

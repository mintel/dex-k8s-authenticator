# Changelog
All notable changes to this project will be documented in this file.

## [v1.4.0]

### Changed

- Bump dex version to `v2.27.0` (security release)
- Switch to `dexidp/dex` container image registry

## [v1.3.0]

### Added

- Pass optional `connector_id` to cluster context (#146)
- Added `trusted_root_ca` to dex-k8s-authenticator helm chart (#143)
- Added `k8s_ca_pem_file` option (#136)
- Allow OIDC scopes per cluster (#129)
- Added namespace field to cluster-config (#124)
- Added HTTP Proxy support (#109)
- Added CircleCI tests

### Fixed

- Fix indentation for `nodeSelector` and `tolerations` in dex-k8s-authenticator helm chart (#137)
- Propgate SIGTERM for graceful shutdown (#110)

## [v1.2.0]

### Added

- Additional tab to display only the id-token
- Service Loadbalancer IP override capability in Helm chart
- Service annotations capability in Helm chart
- Options to specify `idp_ca_pem_file` and `trusted_root_ca_file`
- Support for fixed context name (instead of auto-generated)

### Changed

- Bump dex version to `v2.17.0`
- Bump to `golang:1.12-alpine3.10`
- Switch to Go Modules
- Minor update (skip cluster selection if only 1 cluster defined)

### Fixed

- Fixed Affinity indentation in Helm chart

## [v1.1.0]

### Added

- Documentation on `web_path_prefix`
- Helm charts now add a checksum annotation on the configmap to roll-deployments when configuration changes
- Added IDPCaPem option to support displaying of idp-ca inline


### Changed

- Bump dex version to `v2.13.0` and pull from new repo at quay.io/dexidp/dex
- Documentation improvements

### Fixed

- Fixes to some css to use relative paths

## [v1.0.0]
### Added

- New tabbed layout with clipboard copy options. Key driver for this is to
enable Windows specific instructions.
- Added envar substitutions. Can now generate a config based on values in the
environment (useful for the `client_secret`).
- Added `nodePort` support to Helm charts.
- Added `kubectl_version` option in config. Used to construct a download link to `kubectl` which may be useful.
- Added `web_path_prefix` config option to set url-prefix for serving requests and assets.
- Added `trusted_root_ca` config option to specifiy 1 more root CA's.
- Added `k8s_ca_pem` config option to provide abililty to specify the Kubernetes CA inline.

### Changed

- Use `ClusterName` in preference to `ClientID` when generating k8s context commands
- Bump dex version to `v2.10.0`
- Bump base image to `alpine 3.8`
- Documentation updates.
- Helm chart for dex-k8s-authenticator would fail when caCerts were specified due to breaking naming conventions on and Secret and Volume resources. Introduce a required `filename` option which lets us separate out the filename of the cert and the name of the k8s resource created.
- Slim down final docker container size.

### Fixed

- `update-ca-certificates` only accepts *.crt (only attempt to copy these)

## [v0.4.0]
### Added
- Abililty to provide K8s cert file content in configuration via k8s_ca_pem
cluster option.

### Fixed
- Explicitly define the CA certificate path using ${HOME}

## [v0.3.0]
### Added
- Allow self-signed certs to be used

### Changed
- Bump to golang:1.9.4-alpine3.7

### Fixed
- Fixed helm-chart ingress servicePort

## [v0.2.0]
### Added
- Helm chart serviceAccountName
- Documentation improvements

### Changed
- Helm chart RBAC (renamed some vars).

## [v0.1.0]
### Added
- Initial release

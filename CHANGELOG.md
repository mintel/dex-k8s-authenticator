# Changelog
All notable changes to this project will be documented in this file.

## [UNRELEASED]


### Added

- New tabbed layout with clipboard copy options. Key driver for this is to 
enable Windows specific instructions.

- Added `kubectl_version` option in config. Used to construct a download link
to `kubectl` which may be useful.

- Added envar substitutions. Can now generate a config based on values in the 
environment (useful for the `client_secret`).

### Changed
- Helm chart for dex-k8s-authenticator would fail when caCerts were specified
due to breaking naming conventions on and Secret and Volume resources.
Introduce a required `filename` option which lets us separate out the filename
of the cert and the name of the k8s resource created.
- Slim down final docker container size 

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

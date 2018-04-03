# Changelog
All notable changes to this project will be documented in this file.

## [v0.4.0]
### Added
- Abililty to provide K8s cert file content in configuration via k8s_ca_pem
cluster option.

### Fixed
- Explicitly define the CA certificate path using ${HOME}

## [v0.3.0]
### Added
- Allow self-signed certs to be used

# Changed
- Bump to golang:1.9.4-alpine3.7

# Fixed
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

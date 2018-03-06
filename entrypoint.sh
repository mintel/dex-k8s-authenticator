#!/bin/sh

if [ ! -z "$(ls -A /certs)" ]; then
  cp /certs/*.crt /certs/*.pem /usr/local/share/ca-certificates/
  update-ca-certificates
fi

# Execute dex-k8s-authenticator with any argument passed to docker run
/app/bin/dex-k8s-authenticator $@

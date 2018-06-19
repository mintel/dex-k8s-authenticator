#!/bin/sh

if [ ! -z "$(ls -A /certs)" ]; then
  find /certs/ -type f -exec sh -c 'cp -v -L '{}' /usr/local/share/ca-certificates/$(basename '{}').crt' \;
  update-ca-certificates
fi

# Execute dex-k8s-authenticator with any argument passed to docker run
/app/bin/dex-k8s-authenticator $@

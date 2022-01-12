#!/bin/sh

if [ ! -z "$(ls -A /certs)" ]; then
  cp -L /certs/*.crt /usr/local/share/ca-certificates/ 2>/dev/null
  update-ca-certificates
fi

credentials_file=/var/run/secrets/dex-k8s-authenticator/credentials/credentials
if [ -e $credentials_file ]
then
  echo "export client_secret" `cat $credentials_file` > /var/tmp/credential_export
  source /var/tmp/credential_export
fi

# Execute dex-k8s-authenticator with any argument passed to docker run
/app/bin/dex-k8s-authenticator $@
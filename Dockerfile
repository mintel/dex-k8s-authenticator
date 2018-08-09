FROM golang:1.10.3-alpine3.8

RUN apk add --no-cache --update alpine-sdk bash

COPY . /go/src/github.com/mintel/dex-k8s-authenticator
WORKDIR /go/src/github.com/mintel/dex-k8s-authenticator
RUN make get && make 

FROM alpine:3.8
# Dex connectors, such as GitHub and Google logins require root certificates.
# Proper installations should manage those certificates, but it's a bad user
# experience when this doesn't work out of the box.
#
# OpenSSL is required so wget can query HTTPS endpoints for health checking.
RUN apk add --update ca-certificates openssl curl

RUN mkdir -p /app/bin
COPY --from=0 /go/src/github.com/mintel/dex-k8s-authenticator/bin/dex-k8s-authenticator /app/bin/dex-k8s-authenticator
COPY --from=0 /go/src/github.com/mintel/dex-k8s-authenticator/html /app/html
COPY --from=0 /go/src/github.com/mintel/dex-k8s-authenticator/templates /app/templates

# Add any required certs/key by mounting a volume on /certs - Entrypoint will copy them and run update-ca-certificates at startup
RUN mkdir -p /certs

WORKDIR /app

COPY entrypoint.sh /
RUN chmod a+x /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]

CMD ["--help"]


FROM golang:1.9.2-alpine

RUN apk add --no-cache --update alpine-sdk bash

COPY . /go/src/github.com/mintel/dex-k8s-authenticator
WORKDIR /go/src/github.com/mintel/dex-k8s-authenticator
RUN make get && make 

FROM alpine:3.4
# Dex connectors, such as GitHub and Google logins require root certificates.
# Proper installations should manage those certificates, but it's a bad user
# experience when this doesn't work out of the box.
#
# OpenSSL is required so wget can query HTTPS endpoints for health checking.
RUN apk add --update ca-certificates openssl curl

COPY --from=0 /go/src/github.com/mintel/dex-k8s-authenticator /app

# Add any required certs here... and run update-ca-certificates
WORKDIR /app

ENTRYPOINT ["/app/bin/dex-k8s-authenticator"]

CMD ["--help"]


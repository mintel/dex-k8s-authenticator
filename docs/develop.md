# Developing and Running

## Building a binary

    make 
    
Creates ./bin/dex-k8s-authenticator

## Building a container

    make container

## Running it

### Start a Dex Server instance

You must have a `dex` instance running before starting `dex-k8s-authenticator`.

Follow the example here:
https://github.com/dexidp/website/blob/main/content/docs/getting-started.md

Start it with using the provided `./examples/dex-server-config-dev.yaml`

### Start Dex K8s Authenticator

    ./bin/dex-k8s-authenticator --config ./examples/config.yaml

* Browse to http://localhost:5555
* Click 'Example Cluster'
* Click 'Log in with Email'
* Login with `admin@example.com` followed by the password `password`
* You should be redirected back to the dex-k8s-authenticator

### Configuration Options

Additional configuration options are explained [here](config.md)

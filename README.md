# Temporal GRPC Proxy

## Set up
Set the following environment variables, e.g:
```sh
export TEMPORAL_ADDRESS=my-namespace.my-account.tmprl.cloud:7233
export TEMPORAL_TLS_CERT=/path/to/tls.crt
export TEMPORAL_TLS_KEY=/path/to/tls.key
```

## Run the proxy
```sh
go run main.go
```
The server will listen on `127.0.0.1:7233`

## Run a sample through the proxy
Go to your favorite Temporal sample from any of the samples repos.  In the `worker` and `starter` main functions, update the `client` parameters to set your namespace, e.g. `my-namespace.my-account`.

After the above change, run the sample as documented.  It will now run against Temporal Cloud, through the proxy!

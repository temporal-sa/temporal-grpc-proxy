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
The proxy server will listen on `127.0.0.1:7233`

## Run a sample through the proxy
Go to your favorite Temporal sample from any of the samples repos.  Update the `client` parameters to set your namespace, e.g. `my-namespace.my-account` (see SDK specific examples below).

After the above change, run the sample as documented.  It will now run against Temporal Cloud, through the proxy!

### Go client
```go
c, err := client.Dial(client.Options{
    Namespace: "my-namespace.my-account",
})
```

### Java client
```java
WorkflowClient client = WorkflowClient.newInstance(service, 
    WorkflowClientOptions.newBuilder().setNamespace("my-namespace.my-account").build());
```

### TypeScript client
```typescript
const client = new Client({
    connection,
    namespace: 'my-namespace.my-account',
});
```

### Python client
```python
client = await Client.connect("localhost:7233", 
    namespace="my-namespace.my-account")
```

### .NET client
```csharp
var client = await TemporalClient.ConnectAsync(new("localhost:7233")
{
    Namespace = "my-namespace.my-account",
});
```

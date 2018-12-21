# gRPC Protocol

Mesher support gRPC protocol

### Configurations
To enable gRPC proxy you must set the protocol config 
```yaml
cse:
  protocols:
    grpc:
      listenAddress: 127.0.0.1:40101 # or internalIP:port
```

### How to use mesher as sidecar proxy
Assume you original client is 
```go
	conn, err := grpc.Dial("10.0.1.1:50051",
		grpc.WithInsecure(),
		)
```
set http_proxy
```bash
export http_proxy=http://127.0.0.1:40100
```


## example
A gRPC example is [here](https://github.com/go-mesh/mesher-examples/tree/master/protocol/grpc-go)

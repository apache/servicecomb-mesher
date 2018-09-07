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
you must add custom dialer for you grpc client
Assume you original client is 
```go
	conn, err := grpc.Dial("10.0.1.1:50051",
		grpc.WithInsecure(),
		)
```
after modify 
```go
        //target address is consist of the provider name(in that case "Server") and provider port
	conn, err := grpc.Dial("Server:50051",
		grpc.WithInsecure(),
		grpc.WithDialer(func(addr string, time time.Duration) (net.Conn, error) {
			//127.0.0.1:40101 is local grpc proxy address
			return net.DialTimeout("tcp", "127.0.0.1:40101", time)
		}))
```

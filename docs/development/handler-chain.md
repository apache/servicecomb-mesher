# Handler chain
all the traffic will go through the handler chain.
A chain is composite of handlers, each handler has a particular logic.
Mesher also has a lots of feature working in chain, like route management, circuit breaking, rate-limiting.
In Summary, handler is the middle ware between client and servers, 
it is useful, when you want to add authorization to intercept illegal requests.

### How to write a handler
https://docs.go-chassis.com/dev-guides/how-to-implement-handler.html

### How to use it in handler chain
in chassis.yaml add your handler name in chain configuration.
as side car and API gateway, mesher's chain has different meaning.

For example, running as mesher-sidecar, service A call another service B, 
outgoing chain process all the service A requests before remote call, 
incoming chain process all the requests from service A, before access to service B API. 

In summary outgoing chain works when a service attempt to call other services, 
incoming chain works when other services call this service
```yaml
  handler:
    chain:
      Consumer:
        # if a service call other service, it go through this chain, loadbalance and transport is must 
        outgoing: router, bizkeeper-consumer, loadbalance, transport
      Provider:
        incoming: ratelimiter-provider
```

running as API gateway, 
incoming chain process all the requests from the external network, 
outgoing chain process all the the request between API gateway and backend services
```yaml
  handler:
    chain:
      Consumer:
        #loadbalance and transport is must 
        outgoing: router, bizkeeper-consumer, loadbalance, transport
      Provider:
        incoming: ratelimiter-provider
```
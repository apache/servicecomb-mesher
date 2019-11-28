# Handler chain
All the traffic will go through the handler chain.
A chain is composite of handlers, each handler has a particular logic.
Mesher also has lots of feature working in chain, like route management, circuit breaking  and rate-limiting.
In Summary, handler is the middle ware between clients and servers, 
it is useful when adding authorization to intercept illegal requests.

### How to write a handler
https://docs.go-chassis.com/dev-guides/how-to-implement-handler.html

### How to use it in handler chain
In chassis.yaml add your handler name in chain configuration.
As sidecar and API gateway, mesher's chain has different meanings.

For example, running as mesher-sidecar, service A call another service B, 
outgoing chain of service A processes all the service A requests before remote call, 
incoming chain of service B processes all the requests from service A, before access to service B API. 

In summary, outgoing chain works when a service attempt to call other services, 
incoming chain works when other services call this service.

```yaml
  handler:
    chain:
      Consumer:
        # if a service call other service, it go through this chain, loadbalance and transport is must 
        outgoing: router, bizkeeper-consumer, loadbalance, transport
      Provider:
        incoming: ratelimiter-provider
```

Running as API gateway, 
incoming chain processes all the requests from the external network, 
outgoing chain processes all the the requests between API gateway and back-end services.

```yaml
  handler:
    chain:
      Consumer:
        #loadbalance and transport is must 
        outgoing: router, bizkeeper-consumer, loadbalance, transport
      Provider:
        incoming: ratelimiter-provider
```
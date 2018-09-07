# Local Health check
you can use health checker to check local service health,
when service instance is not healthy, mesher will update the instance status in registry service to "DOWN" 
so that other service
can not discover this instance. If the service is healthy again, mesher will update status to "UP", 
then other instance can discover it again. 
currently this function works only when you use service center as registry

examples:

Check local http service
```yaml
localHealthCheck:
  - port: 8080
    protocol: rest
    uri: /health
    interval: 30s
    match:
      status: 200
      body: ok
```

### Options


**port**
>*(require, string)* must be a port number, mesher is only responsible to check local service, 
it use 127.0.0.1:{port} to check service

**protocol**
>*(optional, string)* mesher has a built-in checker "rest",for other protocol, 
will use default TCP checker unless you implement your own checker

**uri**
>*(optional, string)* uri start with /.


**interval**
>*(optional, string)* check interval, you can use number with unit: 1m, 10s. 

**match.status**
>*(optional, string)* the http response status must match status code

**match.body**
>*(optional, string)* the http response body must match body
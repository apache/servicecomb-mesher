# Local Health check
You can use health checker to check local service health. 
When service instance is unhealthy, mesher will update the instance status in registry service to "DOWN" 
so that other services
can not discover this instance. After the service becoming healthy again, mesher will update the status to "UP", 
then other instance can discover it again. 
Currently this function works only when using service center as registry.

Examples:

- Check local http service

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

>*(require, string)* Must be a port number, mesher is only responsible to check local services, 
it use 127.0.0.1:{port} to check services.

**protocol**

>*(optional, string)* Mesher has a built-in checker "rest",for other protocols, 
will use default TCP checker unless implementing your own checker.

**uri**

>*(optional, string)* Uri start with /.


**interval**
>*(optional, string)* Check interval, you can use number with unit: 1m, 10s. 

**match.status**
>*(optional, string)* The http response status must match status code.

**match.body**

>*(optional, string)* The http response body must match body.
# Destination Resolver

Destination Resolver is a module to parse each protocol request to get a target service name. 
You can write your own resolver implementation for different protocols.

## Configurations

Example:
```yaml
plugin:
  destinationResolver:
    http: host # host is a build-in and default resolver, it uses host name as service name
    grpc: ip
```



**plugin.destinationResolver**

>*(optional, map)* Define what kind of resolver, a protocol should use


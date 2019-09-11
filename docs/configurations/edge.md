# API gateway
mesher is able to work as a API gateway to mange traffic,to run mesher as an API gateway
```shell
mesher --config=mesher.yaml --mode edge
```
the ingress rule is in mesher.yaml

### Options

**mesher.ingress.type**
>*(optional, string)* default is servicecomb, it reads servicecomb ingress rule. 
>it is a plugin, you can custom your own implementation


**mesher.ingress.rule.http**
>*(optional, string)* rule about how to forward http traffic. it holds a yaml content as rule.

below explain the content, the rule list is like a filter, all the request will go through this rule list until match one rule.

**apiPath**
>*(required, string)* if request's url match this, it will use this rule

**host**
>*(optional, string)* if request HOST match this, mesher will use this rule, it can be empty, 
>if you set both host and apiPath, the request's host and api path must match them at the same time
>
**service.name**
>*(required, string)* target backend service name in registry service(like ServiceComb service center)
>
**service.redirectPath**
>*(optional, string)* by default, mesher use original request's url
>
**service.port.value**
>*(optional, string)* if you use java chassis or go chassis to develop backend service, no need to set it. 
>but if your backend service use mesher-sidecar, you must give your service port here.
>
### example
```yaml
mesher:
  ingress:
    type: servicecomb
    rule:
      http: |
        - host: example.com
          apiPath: /some/api
          service:
            name: example
            redirectPath: /another/api
            port:
              name: http-legacy
              value: 8080
        - apiPath: /some/api
          service:
            name: Server
            port:
              name: http
              value: 8080
```


### Enable TLS
generate private key
```sh
openssl genrsa -out server.key 2048
```
sign cert with private key
```shell script
openssl req -new -x509 -key server.key -out server.crt -days 3650
```
set file path in chassis.yaml
```yaml
ssl:
  mesher-edge.rest.Provider.certFile: server.crt
  mesher-edge.rest.Provider.keyFile: server.key
```

To know advanced feature about TLS configuration, check 
https://docs.go-chassis.com/user-guides/tls.html
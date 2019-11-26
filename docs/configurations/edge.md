# API gateway
Mesher is able to work as a API gateway to mange traffic, to run mesher as an API gateway:
```shell
mesher --config=mesher.yaml --mode edge
```
The ingress rule is in mesher.yaml.

### Options

**mesher.ingress.type**
>*(optional, string)* Default is servicecomb, it reads servicecomb ingress rule. 
>It is a plugin, you can custom your own implementation.

**mesher.ingress.rule.http**

>*(optional, string)* Rule about how to forward http traffic. It holds a yaml content as rule.



Below explaining the content, the rule list is like a filter, all the request will go through this rule list until matching one rule.

**apiPath**

>*(required, string)* If request's url matches this, it will use this rule.

**host**
>*(optional, string)* If request HOST matches this, mesher will use this rule. It can be empty. 
>If you set both host and apiPath, the request's host and api path must match them both.
>
**service.name**
>*(required, string)* Target back-end service name in registry service (like ServiceComb service center).
>
**service.redirectPath**
>*(optional, string)* By default, mesher uses original request's url.
>
**service.port.value**
>*(optional, string)* If using java chassis or go chassis to develop back-end service, no need to set it. 
>But if back-end service uses mesher-sidecar, service port must be given here.
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
Generate private key
```sh
openssl genrsa -out server.key 2048
```
Sign cert with private key
```shell script
openssl req -new -x509 -key server.key -out server.crt -days 3650
```
Set file path in chassis.yaml
```yaml
ssl:
  mesher-edge.rest.Provider.certFile: server.crt
  mesher-edge.rest.Provider.keyFile: server.key
```

To know advanced feature about TLS configuration, check 
https://docs.go-chassis.com/user-guides/tls.html
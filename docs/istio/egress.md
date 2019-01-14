# Egress
## Introduction

Mesher support Egress for your service, so that you can access any publicly accessible service from your microservice.
 
## Configuration

The egress related configurations is all in egress.yaml.

**infra**
> *(optional, string)* specifies from where the egress configuration need to be taken supports two values cse or pilot ,
      cse means the egress configuration from egress.yaml file,
      pilot means egress configuaration are taken from pilot of istio,
      default is *cse*
      
**address**
> *(optional, string)* The end point of pilot from which configuration need to be fetched.

**hosts**
> *(optional, []string)* host associated with external service, could be a DNS name with wildcard prefix

**ports.port**
> *(optional, int)* The port associated with the external service, default is *80*

**ports.protocol**
> *(optional, int)* The protocol associated with the external service,supports only http default is *HTTP*

## example
edit egress.yaml

```yaml
egress:
  infra: cse  # pilot or cse
  address: http://istio-pilot.istio-system:15010
egressRule:
  google-ext:
    - hosts:
        - "www.google.com"
        - "*.yahoo.com"
      ports:
        - port: 80
          protocol: HTTP
```

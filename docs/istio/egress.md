# Egress
## Introduction

Mesher support Egress for your service, so that you can access any publicly accessible services from your microservices.

## Configuration

The egress related configurations are all in egress.yaml.

**infra**

> *(optional, string)* Specifies from where the egress configuration need to be taken supports two values CSE or pilot ,
      CSE means the egress configurations from egress.yaml file,
      pilot means egress configurations are taken from pilot of Istio,
      default is *CSE*.

**address**

> *(optional, string)* The end point of pilot from which configuration need to be fetched.

**hosts**
> *(optional, []string)* Host associated with external service, could be a DNS name with wildcard prefix.

**ports.port**

> *(optional, int)* The port associated with the external service, default is *80*.

**ports.protocol**

> *(optional, int)* The protocol associated with the external service, supports only HTTP, default is *HTTP*.

## Example
Edit egress.yaml

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

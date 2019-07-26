# Mesher

[![Build Status](https://travis-ci.org/apache/servicecomb-mesher.svg?branch=master)](https://travis-ci.org/apache/servicecomb-mesher) [![Coverage Status](https://coveralls.io/repos/github/apache/servicecomb-mesher/badge.svg?branch=master)](https://coveralls.io/github/apache/servicecomb-mesher?branch=master) [![Go Report Card](https://goreportcard.com/badge/github.com/apache/servicecomb-mesher)](https://goreportcard.com/report/github.com/apache/servicecomb-mesher) [![GoDoc](https://godoc.org/github.com/apache/servicecomb-mesher?status.svg)](https://godoc.org/github.com/apache/servicecomb-mesher) 

A service mesh implementation based on [go chassis](https://github.com/go-chassis/go-chassis).

One big advantage of Mesher is it is able to 
work with go-chassis in same service mesh control plane like Istio. Without Istio they can work 
together with ServiceComb Service center and running on any infrastructure(docker, VM, baremetal). 
So if you choose go as your programing language, you can use go-chassis to gain better performance, and you can freely use 
other programing language which suit your service the most

Mesher support both linux and windows OS, 
that makes possible that .Net service can work with java, go, python service in one distributed system easily

# Features
- go-chassis: Mesher has all of features of [go chassis](https://github.com/go-chassis/go-chassis)
a go micro service framework
- Admin API：Listen on isolated port, let user to query runtime information 


# Get started
Refer to [mesher-examples](https://github.com/go-mesh/mesher-examples)

### How to build and run
#### Build from scratch
1. Install ServiceComb [service-center](https://github.com/ServiceComb/service-center/releases)

2. build and run, use go mod(go 1.11+, experimental but a recommended way)
```shell
cd mesher
GO111MODULE=on go mod download
#optional
GO111MODULE=on go mod vendor
go build mesher.go
./mesher
```
####Build script
```bash
cd build
./build_proxy.sh

```
it will build binary and docker image
- tar file: release/mesher-latest-linux-amd64.tar
- docker: servicecomb/mesher:latest

# Documentations

https://mesher.readthedocs.io/en/latest/

# Test with servicecomb-java-chassis
You can start with a easy use case which uses servicecomb-java-chassis sample [calculator](https://github.com/apache/servicecomb-java-chassis/tree/master/samples/bmi/calculator)

* First go build mesher.go and you can get a executable file ./mesher
## Create two mersher project
```bash
mkdir /usr/local/src/mersher-consumer
cp ./mesher /usr/local/src/mersher-consumer
cp -r ./conf /usr/local/src/mersher-consumer

mkdir  /usr/local/src/mersher-provider
cp ./mesher /usr/local/src/mersher-provider
cp -r ./conf /usr/local/src/mersher-provider
```
### Begin to edit conf
- First you need exec ifconfig in linux to get your intranet ip.
for example 192.168.88.99.
#### Edit  mersher-consumer/conf
```bash
vi /usr/local/src/mersher-consumer/conf/microservice.yaml
```
- Content as below:
```
## microservice property
service_description:
  name: mersher-consumer
  version: 0.0.1
  environment:  #microservice environment
  properties:
    allowCrossApp: true #whether to allow calls across applications
```
- Edit conf ip:
```bash
vi /usr/local/src/mersher-consumer/conf/chassis.yaml
```
- Just need to edit:
*127.0.0.1 to your intranet ip*
- For example:
```
grpc:
      listenAddress: 192.168.86.99:40101
    http:
      listenAddress: 192.168.86.99:30101
    rest-admin:
      listenAddress: 192.168.86.99:30102 # listen addr use to adminAPI
```
#### Edit  mersher-provider/conf
```bash
vi /usr/local/src/mersher-provider/conf/microservice.yaml
```
- Content as below:
```
## microservice property
service_description:
  name: mersher-provider
  version: 0.0.1
  environment:  #microservice environment
  properties:
    allowCrossApp: true #whether to allow calls across applications
```
- Edit conf ip:
```bash
vi /usr/local/src/mersher-provider/conf/chassis.yaml
```
- Just need to edit:
*127.0.0.1 to your intranet ip and make sure port does not conflict 
- For example:
```
grpc:
      listenAddress: 192.168.86.99:40102
    http:
      listenAddress: 192.168.86.99:30108
    rest-admin:
      listenAddress: 192.168.86.99:30109 # listen addr use to adminAPI
```
## Install and run Service-Center 
- You can easily do it by [ServiceComb Quick Start](http://servicecomb.apache.org/cn/docs/quick-start/)
- [service-center]http://servicecomb.apache.org/cn/release/
```
cd /usr/local/src/apache-servicecomb-service-center-1.2.0-linux-amd64;./start-service-center.sh
```
- This server is used for Service discovery

## Install and run calculator
- You can easily do it by [ServiceComb Quick Start](http://servicecomb.apache.org/cn/docs/quick-start/)
## Test 
```
export http_proxy=http://127.0.0.1:30101
curl -X GET -i -v http://mersher-provider:7777/bmi?height=180\&\&weight=100
```
- You will get result:
```
Note: Unnecessary use of -X or --request, GET is already inferred.
*   Trying 127.0.0.1...
* TCP_NODELAY set
* Connected to 127.0.0.1 (127.0.0.1) port 30101 (#0)
> GET http://mersher-provider:7777/bmi?height=180&&weight=100 HTTP/1.1
> Host: mersher-provider:7777
> User-Agent: curl/7.58.0
> Accept: */*
> Proxy-Connection: Keep-Alive
>
< HTTP/1.1 200 OK
HTTP/1.1 200 OK
< Content-Length: 85
Content-Length: 85
< Content-Type: application/json; charset=utf-8
Content-Type: application/json; charset=utf-8
< Date: Fri, 26 Jul 2019 08:47:46 GMT
Date: Fri, 26 Jul 2019 08:47:46 GMT

<
* Connection #0 to host 127.0.0.1 left intact
{"result":30.9,"instanceId":"07ef9721af8211e9965dfa163ef423d3","callTime":"16:47:46"}
```

- And you can get log by:
```
 grep "Create client" /usr/local/src/mersher-consumer/log/mesher.log
 {"level":"INFO","timestamp":"2019-07-26 16:48:50.497 +08:00","file":"client/client_manager.go:104","msg":"Create client for rest:mersher-provider:192.168.86.99:30108"}
```
```
grep "Create client" /usr/local/src/mersher-provider/log/mesher.log
{"level":"INFO","timestamp":"2019-07-26 16:48:50.497 +08:00","file":"client/client_manager.go:104","msg":"Create client for rest:mersher-provider:127.0.0.1:7777"}
```
- This result means that we are executing the following call chain:  
curl client -> mersher-consumer -> mersher-provider -> calculator







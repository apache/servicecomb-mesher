# Mesher

[![Build Status](https://travis-ci.org/apache/servicecomb-mesher.svg?branch=master)](https://travis-ci.org/apache/servicecomb-mesher) 
[![Coverage Status](https://coveralls.io/repos/github/apache/servicecomb-mesher/badge.svg?branch=master)](https://coveralls.io/github/apache/servicecomb-mesher?branch=master) 
[![Go Report Card](https://goreportcard.com/badge/github.com/apache/servicecomb-mesher)](https://goreportcard.com/report/github.com/apache/servicecomb-mesher) 
[![GoDoc](https://godoc.org/github.com/apache/servicecomb-mesher?status.svg)](https://godoc.org/github.com/apache/servicecomb-mesher) 

A service mesh implementation based on [go chassis](https://github.com/go-chassis/go-chassis).

# Why use mesher
- any infrastructure: if you use ServiceComb as control plane, you can run on any infrastructure(docker, kubernetes,VM, bare metal). 
- service mesh and frameworks: 
you can develop micro services with java chassis or go chassis frameworks 
and use mesher to make other service join to the same system.
- flexible: you can develop and customize your own service mesh
- OS: support both linux and windows OS, which means you can govern your services writen in .net with java, go etc.
- API gateway: mesher is able to run as a independent edge service and manage ingress traffic.
# Features
- Build on top of go micro service framework: so that mesher has all of features of 
[go chassis](https://github.com/go-chassis/go-chassis),a high flexible go micro service framework. 
you can custom your own service mesh and API gateway by extending lots of components.
- Admin API：Listen on an isolated port, expose useful runtime information and metrics.
- support protocols: http and grpc
- No IP tables forwarding: Mesher leverage 
[http_proxy](http://kaamka.blogspot.com/2009/06/httpproxy-environment-variable.html) 
and [grpc proxy dialer](https://godoc.org/google.golang.org/grpc#WithDialer), 
that makes better performance than using ip tables
- local health check: as a sidecar, mesher is able to check local service health by policy
and dynamically remove it from service registry if service is unavailable.

# Get started
Refer to [mesher-examples](https://github.com/go-mesh/mesher-examples)

### How to build
#### Build from scratch
1. Install ServiceComb [service-center](https://github.com/ServiceComb/service-center/releases)

2. build and run, use go mod
```shell
export GOPROXY=https://goproxy.io #if you are facing network issue
cd mesher
GO111MODULE=on go mod download
#optional
GO111MODULE=on go mod vendor
cd cmd/mesher
go build mesher.go
```
#### Build by script

```bash
cd build
export GOPATH=/path/to/gopath
export GOPROXY=https://goproxy.io #if you are facing network issue
./build_proxy.sh

```
it will build binary and docker image
- tar file: release/mesher-latest-linux-amd64.tar
- docker image name: servicecomb/mesher-sidecar:latest

#### Build mesher edge service
```shell script
sudo docker build -t servicecomb/mesher-edge -f docker/edge/Dockerfile .
```

### How to run
#### edge service
```shell script
sudo docker run -e PAAS_CSE_SC_ENDPOINT="" \
-e PAAS_CSE_CC_ENDPOINT="" \
-p 80:30101 \
--name mesher-edge servicecomb/mesher-edge
```
# Documentations
You can see more documentations in [here](https://mesher.readthedocs.io/en/latest/), 
this online doc is for latest version of mesher, if you want to see your version's doc,
follow [here](docs/README.md) to generate it in local

# RoadMap
## 1.7.0
- mesher ingress: mesher as a edge service to manage ingress traffic
- integrate with servicecomb kie: user can change mesher-sidecar and mesher-ingress in same config server, don't need to change local files
## in future
- support skywalking

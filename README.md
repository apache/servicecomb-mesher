# Mesher

[![Build Status](https://travis-ci.org/apache/servicecomb-mesher.svg?branch=master)](https://travis-ci.org/apache/servicecomb-mesher) 
[![Coverage Status](https://coveralls.io/repos/github/apache/servicecomb-mesher/badge.svg?branch=master)](https://coveralls.io/github/apache/servicecomb-mesher?branch=master) 
[![Go Report Card](https://goreportcard.com/badge/github.com/apache/servicecomb-mesher)](https://goreportcard.com/report/github.com/apache/servicecomb-mesher) 
[![GoDoc](https://godoc.org/github.com/apache/servicecomb-mesher?status.svg)](https://godoc.org/github.com/apache/servicecomb-mesher) 

A service mesh implementation based on [go chassis](https://github.com/go-chassis/go-chassis).

# Why use mesher
It leverages Istio or ServiceComb as control plane. 
if you use ServiceComb as control plane, you can run on any infrastructure(docker, kubernetes,VM, bare metal). 
Besides you can develop java and go services with java chassis or go chassis to gain better performance.

Mesher support both linux and windows OS, 
which means you can govern your services writen in .net with java, go etc.
# Features
- go-chassis: Mesher has all of features of [go chassis](https://github.com/go-chassis/go-chassis)
a go micro service framework
- Admin API：Listen on isolated port, expose useful runtime information 
- support protocols: http, grpc

# Get started
Refer to [mesher-examples](https://github.com/go-mesh/mesher-examples)

### How to build
#### Build from scratch
1. Install ServiceComb [service-center](https://github.com/ServiceComb/service-center/releases)

2. build and run, use go mod(go 1.11+, experimental but a recommended way)
```shell
export GOPROXY=https://goproxy.io #if you are facing network issue
cd mesher
GO111MODULE=on go mod download
#optional
GO111MODULE=on go mod vendor
cd cmd/mesher
go build mesher.go
```
####Build by script

```bash
cd build
export GOPATH=/path/to/gopath
export GOPROXY=https://goproxy.io #if you are facing network issue
./build_proxy.sh

```
it will build binary and docker image
- tar file: release/mesher-latest-linux-amd64.tar
- docker: servicecomb/mesher-sidecar:latest

# Documentations

https://mesher.readthedocs.io/en/latest/

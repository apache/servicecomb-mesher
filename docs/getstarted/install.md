# Quick start

### Local
In this case, you will launch one mesher sidecar proxy and 
one service developed based on go-chassis as provider
and use curl as a dummy consumer to access this service

the network traffic: curl->mesher->service


1.Install ServiceComb [service-center](https://github.com/apache/incubator-servicecomb-service-center/releases)

2.Install [go-chassis](https://go-chassis.readthedocs.io/en/latest/getstarted/install.html) and 
run [rest server](https://github.com/go-chassis/go-chassis/tree/master/examples/rest/server)

2. Build and run, use go mod(go 1.11+, experimental but a recommended way)
```shell
cd mesher
GO111MODULE=on go mod download
#optional
GO111MODULE=on go mod vendor
go build mesher.go
./mesher
```

 
4.verify, in this case curl command is the consumer, mesher is consumer's sidecar, 
and rest server is provider
```shell
export http_proxy=http://127.0.0.1:30101
curl http://RESTServer:8083/sayhello/peter
```

**Notice**:
>>You don't need to set service registry in chassis.yaml, 
because by default registry address is 127.0.0.1:30100, 
just same service center default listen address.

>>curl command read lower case http_proxy environment variable.

### Run on different infrastructure

Mesher does not bind to any platform or infrastructures, plz refer to 
https://github.com/go-mesh/mesher-examples/tree/master/Infrastructure
to know how to run mesher on different infra

### Sidecar injector
Mesher supply a way to automatically inject mesher configurations in kubernetes

See detail https://github.com/go-chassis/sidecar-injector

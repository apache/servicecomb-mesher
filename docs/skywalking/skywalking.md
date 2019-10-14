# SkyWalking

Skywalking-manager is a handler plugin of mesher, it reports tracing data to skywalking server

## Configurations
**In conf/mesher.conf**

**apm.tracing.enable**
>  *(optional, bool)* enable apm

**apm.tracing.serverUri**
>  *(optional, string)* server address of skywalking

## Example
```yaml
apm:
  tracing:
    enable: true
    serverUri: 127.0.0.1:11800
```

## SkyWawlking-Manager Init
**In file proxy/bootstrap/bootstrap.go**
```shell script
import "github.com/apache/servicecomb-mesher/proxy/pkg/skywalking"
```
**In function Start()**
```shell script
skywalking.Init()
```

**In function SetHandlers() add two Handlers**
```shell script
consumerChain := strings.Join([]string{
		...
		skywalking.SkyWalkingConsumer,
		chassisHandler.Transport,
}}

providerChain := strings.Join([]string{
		...
		skywalking.SkyWalkingProvider,
		chassisHandler.Transport,
}}

```
## SkyWalking-Handler Init
**In cmd/mesher/mesher.go**
```shell script
import _ "github.com/apache/servicecomb-mesher/proxy/handler"
```

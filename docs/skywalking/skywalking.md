# SkyWalking

Skywalking-manager is a handler plugin of mesher, it reports tracing data to skywalking server

## Configurations
**In conf/mesher.conf**

**appPerfManage.apmName**
>  *(optional, string)* apm server name, here is skywalking

**appPerfManage.serverUri**
>  *(optional, string)* server address of skywalking

## Example
```yaml
appPerfManage:
  apmName: skywalking
  serverUri: 192.168.88.64:11800
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

**In function SetHandlers() add Handler**
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

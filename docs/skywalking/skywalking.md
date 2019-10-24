# SkyWalking

Skywalking-manager is a handler plugin of mesher, it reports tracing data to skywalking server.

## Configurations
**In conf/mesher.conf**

**servicecomb.apm.tracing.enable**
>  *(optional, bool)* enable application performance manager

**servicecomb.apm.tracing.serverUri**
>  *(optional, string)* server address of skywalking

## Example
```yaml
servicecomb:
  apm:                              #application performance monitor
    tracing:
      enable: true                  #enable tracing ability
      serverUri: 127.0.0.1:11800    #url of skywalking 
```
## Step：

# 1. SkyWawlking-Manager Init
**You must init skywawlking manager pkg which will manage connection and report msg to skywalking**
- For example:
- [1] You can import skywalking manager proxy/pkg/skywalking in file proxy/bootstrap/bootstrap.go.
- [2] Calling function Init() in proxy/pkg/skywalking manually to init skywalking manager.
- [3] Adding skywalking's consumer handler name SkyWalkingConsumer defined in proxy/pkg/skywalking to consumerChain.
- [4] Adding skywalking's provider handler name SkyWalkingProvider defined in proxy/pkg/skywalking to providerChain.
- more details about handler chains in [go-chassis](https://github.com/go-chassis/go-chassis#readme)

# 2. SkyWalking-Handler Init
- You must import proxy/handler pkg to init skywalking handler. Not only skywalking handler, all the handlers which are customized for mesher are defined here.
- For example you can import handler pkg in file cmd/mesher/mesher.go


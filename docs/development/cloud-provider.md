# Cloud Provider
By default Mesher do not support any cloud provider.
But there is plugin that helps mesher do it.

### Huawei Cloud 
Mesher is able to use huawei cloud ServiceComb engine. 
#### Access ServiceComb Engine API
Import auth in cmd/mesher/mesher.go
```go
import _ "github.com/huaweicse/auth/adaptor/gochassis"
```

It will sign all requests from mesher to ServiceComb Engine.

#### Use Config Center to manage configuration
Mesher uses servicecomb-kie as config server, 
```go
_ "github.com/apache/servicecomb-kie"
```
When you need to use ServiceComb Engine, you must replace this line. 
Import config center in cmd/mesher/mesher.go.

```go
_ "github.com/go-chassis/go-chassis-config/configcenter"
```
Set the config center in chassis.yaml
```yaml
  config:
    client:
      serverUri: https://xxx #endpoint of servicecomb engine
      refreshMode: 1 # 1: only pull config.
      refreshInterval: 30 # unit is second
      type: config_center
```
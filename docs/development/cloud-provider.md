# Cloud Provider
By default Mesher do not support any cloud provider.
But there is plugin that helps mesher to do it.

### Huawei Cloud 
Mesher is able to use huawei cloud ServiceComb engine. 
#### Access ServiceComb Engine API
import auth in cmd/mesher/mesher.go
```go
import _ "github.com/huaweicse/auth/adaptor/gochassis"
```

it will sign all requests between mesher to ServiceComb Engine.

#### Use Config Center to manage configuration
Mesher use servicecomb-kie as config server, 
```go
_ "github.com/go-chassis/go-chassis-config/servicecombkie"
```
when you need to use ServiceComb Engine, you must replace this line. 
import config center in cmd/mesher/mesher.go.
```go
_ "github.com/go-chassis/go-chassis-config/configcenter"
```
set the config center in chassis.yaml
```yaml
  config:
    client:
      serverUri: https://xxx #endpoint of servicecomb engine
      refreshMode: 1 # 1: only pull config.
      refreshInterval: 30 # unit is second
      type: config_center
```
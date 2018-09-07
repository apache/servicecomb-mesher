# Admin API

### Configurations

admin api server leverage protocol server, it listens on isolated port, by default admin is enabled, and default value of goRuntimeMetrics is false.

To start api server, set protocol server config in chassis.yaml
```yaml
cse:
   protocols:
     rest-admin:
       listenAddress: 0.0.0.0:30102  # listen addr for adminAPI
```

tune admin api in mesher.yaml
```yaml
admin: 
  enable: true
  goRuntimeMetrics : true # enable metrics
```


**admin.enable**
>*(optional, bool)* default is false

**admin.goRuntimeMetrics**
>*(optional, bool)* default is false, enable to expose go runtime metrics in /v1/mesher/metrics


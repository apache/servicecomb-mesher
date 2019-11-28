# Admin API

### Configurations

Admin API server leverages protocol server, it listens on isolated port. By default admin is enabled, and default value of goRuntimeMetrics is false.

To start API server, set protocol server config in chassis.yaml:
```yaml
cse:
   protocols:
     rest-admin:
       listenAddress: 0.0.0.0:30102  # listen addr for admin API
```

Tune admin api in mesher.yaml:
```yaml
admin: 
  enable: true
```


**admin.enable**
>*(optional, bool)* Default is false



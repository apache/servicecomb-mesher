# Mesher command Line 
when you start mesher process, you can use mesher command line to specify configurations like below
```shell
mesher --config=mesher.yaml --service-ports=rest:8080
```


### Options


**--config**
>*(optional, string)* the path to mesher configuration file, default value is {current_bin_work_dir}/conf/mesher.yaml


**--mode**
>*(optional, string)* mesher has 2 work mode, sidecar and per-host, default is sidecar


**--service-ports**
>*(optional, string)* running as sidecar, mesher need to know local service ports, 
this is to tell mesher service port list, 
The value format is {protocol}-{suffix} or {protocol}
if service has multiple protocol, you can separate with comma "rest-admin:8080,grpc:9000". 
default is empty, in that case mesher will use header X-Forwarded-Port as local service port, 
if it is empty also mesher can not communicate to your local service

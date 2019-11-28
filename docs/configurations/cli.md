# Mesher command Line 
When you start a mesher process, you can use mesher command line to specify configurations as follows:
```shell
mesher --config=mesher.yaml --service-ports=rest:8080
```


### Options

**--config**

>*(optional, string)* The path to mesher configuration file, default value is {current_bin_work_dir}/conf/mesher.yaml

**--mode**

>*(optional, string)* Mesher has 2 work modes, sidecar and edge, default is sidecar

**--service-ports**

>*(optional, string)* Running as sidecar, mesher needs to know local service ports, 
this is to tell mesher service port list, 
The value format is {protocol}-{suffix} or {protocol}. 
If service has multiple protocols, you can separate with comma "rest-admin:8080, grpc:9000", 
default is empty. In that case mesher will use header X-Forwarded-Port as local service port, 
if header X-Forwarded-Port is also empty, mesher can not communicate to your local service.

# Before you start
Before you start, you must know what you gonna do if you use mesher as your sidecar proxy. 

Assume you launched 2 services, 
each of service has a dedicated mesher as sidecar proxy.

The network traffic will be: ServiceA->mesherA->mesherB->ServiceB.

To run mesher along with your services, you need to set minimum configurations as below:


1. Give mesher your service name in microservice.yaml file 
2. Set service discovery service(service center, Istio etc) configurations in chassis.yaml
3. export HTTP_PROXY=http://127.0.0.1:30101 as your service runtime environment
4. (optional)Give mesher your service port list by ENV SERVICE_PORTS or CLI --service-ports


After the configurations, assume you serviceB is listening at 127.0.0.1:8080

the serviceA must use http://ServiceB:8080/{api_path} to access ServiceB

Now you can launch as many as serviceA and serviceB to make this system become a distributed system

**Notice**:
 >> consumer need to use http://provider_name:provider_port/ to access provider,
 instead of http://provider_ip:provider_port/. 
 if you choose to set step4, then you can simply use http://provider_name/ to access provider
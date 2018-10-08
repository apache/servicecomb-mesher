Discovery
======================

----

Introduction
++++

Istio Pilot can be integrated with Mesher, working as the Service Discovery component.

Configuration
++++

edit chassis.yaml.

**registrator.disabled**

 Must disable registrator, because registrator is is used in client side discovery. mesher leverage server side discovery which is supported by kubernetes

**serviceDiscovery.type**

 specify the discovery plugin type to "pilot" or "pilotv2", since Istio removes the xDS v1 API support from version 0.8, if you use Istio 0.8 or higher, make sure to set type to pilotv2.

**serviceDiscovery.address**

 the pilot address, in a Istio environment, for xDS v1 API, pilot usually listens on the http port 8080, while for xDS v2 API, it becomes a grpc port 15010.


examples
++++

::

  cse: # Using xDS v1 API
    service:
      Registry:
        registrator:
          disabled: true
        serviceDiscovery:
          type: pilot
          address: http://istio-pilot.istio-system:8080

::

  cse: # Using xDS v2 API
    service:
      Registry:
        registrator:
          disabled: true
        serviceDiscovery:
          type: pilotv2
          address: grpc://istio-pilot.istio-system:15010
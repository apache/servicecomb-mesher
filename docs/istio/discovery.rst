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

 specify the discovery plugin type to "pilotv2", since Istio removes the xDS v1 API support from version 0.7.1, type "pilot" is deprecated.

**serviceDiscovery.address**

 the pilot address, in a typical Istio environment, pilot usually listens on the grpc port 15010.


examples
++++

::

  cse: # Using xDS v2 API
    service:
      Registry:
        registrator:
          disabled: true
        serviceDiscovery:
          type: pilotv2
          address: grpc://istio-pilot.istio-system:15010

# Get started

Istio Pilot can be configured as the service discovery component for mesher. By default the Pilot plugin is not compiled into mesher binary. To make mesher work with Pilot, import the plugin in mesher's entrypoint source code:

```go
import _ "github.com/apache/servicecomb-mesher/plugins/registry/istiov2"
```

Then the Pilot plugin will be installed when mesher starts. Next step, configure Pilot as service discovery in `chassis.yaml`:

```yaml
cse:
  service:
    registry:
      registrator:
        disabled: true
      serviceDiscovery:
        type: pilotv2
        address: grpc://istio-pilot.istio-system:15010
```

Since mesher doesn't have to register the service to Pilot, the registrator config item should be disabled. Make serviceDiscovery.type to be pilotv2, to get service information by xDS v2 API(the v1 API is deprecated).

### The routing tags in Istio

In the original mesher configuration, user can specify tag based route rules, as described below:

```yaml
## router.yaml
router:
  infra: cse
routeRule:
  targetService:
    - precedence: 2
      route:
      - tags:
          version: v1
        weight: 40
      - tags:
          version: v2
          debug: true
        weight: 40
      - tags:
          version: v3
        weight: 20
```

Then in a typical Istio environment, which is likely to be Kubernetes cluster, user can specify the DestinationRules for targetService with the same tags:

```yaml
apiVersion: networking.istio.io/v1alpha3
kind: DestinationRule
metadata:
  name: targetService
spec:
  host: targetService
  subsets:
  - name: v1
    labels:
      version: v1
  - name: v2
    labels:
      version: v2
      debug: "true"
  - name: v3
    labels:
      version: v3
```

Notice that the subsets' tags are the same with those in `router.yaml`, then mesher's tag based load balancing strategy works as it originally does.

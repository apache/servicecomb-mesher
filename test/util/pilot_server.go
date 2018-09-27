package util

import (
	"sync"
	"testing"

	iputil "github.com/go-chassis/go-chassis/pkg/util/iputil"
	testenv "istio.io/istio/mixer/test/client/env"
	"istio.io/istio/pilot/pkg/bootstrap"
	"istio.io/istio/pilot/pkg/model"
	istiotestutil "istio.io/istio/tests/util"
)

var (
	// mixer-style test environment, includes mixer and envoy configs.
	initMutex sync.Mutex

	// service1 and service2 are used by mixer tests. Use 'service3' and 'app3' for pilot
	// local tests.

	// 10.10.0.0/24 is service CIDR range

	// 10.0.0.0/9 is instance CIDR range
	app3Ip    = "10.2.0.1"
	gatewayIP = "10.3.0.1"
	ingressIP = "10.3.0.2"
	localIP   = "10.3.0.3"
)

// InitLocalPilotTestEnv creates a local, in process Pilot with XDSv2 support and a set
// of common test configs. This is a singleton server, reused for all tests in this package.
//
// The server will have a set of pre-defined instances and services, and read CRDs from the
// common tests/testdata directory.
func InitLocalPilotTestEnv(t *testing.T) *bootstrap.Server {
	initMutex.Lock()
	defer initMutex.Unlock()

	ports := testenv.NewPorts(testenv.XDSTest)

	server := istiotestutil.EnsureTestServer()

	localIP = iputil.GetLocalIP()

	// Service and endpoints for hello.default - used in v1 pilot tests
	hostname := model.Hostname("hello.default.svc.cluster.local")
	server.EnvoyXdsServer.MemRegistry.AddService(hostname, &model.Service{
		Hostname: hostname,
		Address:  "10.10.0.3",
		Ports:    testPorts(0),
	})
	server.EnvoyXdsServer.MemRegistry.AddInstance(hostname, &model.ServiceInstance{
		Endpoint: model.NetworkEndpoint{
			Address: "127.0.0.1",
			Port:    int(ports.BackendPort),
			ServicePort: &model.Port{
				Name:     "http",
				Port:     80,
				Protocol: model.ProtocolHTTP,
			},
		},
		AvailabilityZone: "az",
	})

	// "local" service points to the current host and the in-process mixer http test endpoint
	server.EnvoyXdsServer.MemRegistry.AddService("local.default.svc.cluster.local", &model.Service{
		Hostname: "local.default.svc.cluster.local",
		Address:  "10.10.0.4",
		Ports: []*model.Port{
			{
				Name:     "http",
				Port:     80,
				Protocol: model.ProtocolHTTP,
			}},
	})
	server.EnvoyXdsServer.MemRegistry.AddInstance("local.default.svc.cluster.local", &model.ServiceInstance{
		Endpoint: model.NetworkEndpoint{
			Address: localIP,
			Port:    int(ports.BackendPort),
			ServicePort: &model.Port{
				Name:     "http",
				Port:     80,
				Protocol: model.ProtocolHTTP,
			},
		},
		AvailabilityZone: "az",
	})

	// Explicit test service, in the v2 memory registry. Similar with mock.MakeService,
	// but easier to read.
	server.EnvoyXdsServer.MemRegistry.AddService("service3.default.svc.cluster.local", &model.Service{
		Hostname: "service3.default.svc.cluster.local",
		Address:  "10.10.0.1",
		Ports:    testPorts(0),
	})

	server.EnvoyXdsServer.MemRegistry.AddInstance("service3.default.svc.cluster.local", &model.ServiceInstance{
		Endpoint: model.NetworkEndpoint{
			Address: app3Ip,
			Port:    2080,
			ServicePort: &model.Port{
				Name:     "http-main",
				Port:     1080,
				Protocol: model.ProtocolHTTP,
			},
		},
		Labels:           map[string]string{"version": "v1"},
		AvailabilityZone: "az",
	})
	server.EnvoyXdsServer.MemRegistry.AddInstance("service3.default.svc.cluster.local", &model.ServiceInstance{
		Endpoint: model.NetworkEndpoint{
			Address: gatewayIP,
			Port:    2080,
			ServicePort: &model.Port{
				Name:     "http-main",
				Port:     1080,
				Protocol: model.ProtocolHTTP,
			},
		},
		Labels:           map[string]string{"version": "v2", "app": "my-gateway-controller"},
		AvailabilityZone: "az",
	})

	// Mock ingress service
	server.EnvoyXdsServer.MemRegistry.AddService("istio-ingress.istio-system.svc.cluster.local", &model.Service{
		Hostname: "istio-ingress.istio-system.svc.cluster.local",
		Address:  "10.10.0.2",
		Ports: []*model.Port{
			{
				Name:     "http",
				Port:     80,
				Protocol: model.ProtocolHTTP,
			},
			{
				Name:     "https",
				Port:     443,
				Protocol: model.ProtocolHTTPS,
			},
		},
	})
	server.EnvoyXdsServer.MemRegistry.AddInstance("istio-ingress.istio-system.svc.cluster.local", &model.ServiceInstance{
		Endpoint: model.NetworkEndpoint{
			Address: ingressIP,
			Port:    80,
			ServicePort: &model.Port{
				Name:     "http",
				Port:     80,
				Protocol: model.ProtocolHTTP,
			},
		},
		Labels:           model.IstioIngressWorkloadLabels,
		AvailabilityZone: "az",
	})
	server.EnvoyXdsServer.MemRegistry.AddInstance("istio-ingress.istio-system.svc.cluster.local", &model.ServiceInstance{
		Endpoint: model.NetworkEndpoint{
			Address: ingressIP,
			Port:    443,
			ServicePort: &model.Port{
				Name:     "https",
				Port:     443,
				Protocol: model.ProtocolHTTPS,
			},
		},
		Labels:           model.IstioIngressWorkloadLabels,
		AvailabilityZone: "az",
	})

	//RouteConf Service4 is using port 80, to test that we generate multiple clusters (regression)
	// service4 has no endpoints
	server.EnvoyXdsServer.MemRegistry.AddService("service4.default.svc.cluster.local", &model.Service{
		Hostname: "service4.default.svc.cluster.local",
		Address:  "10.1.0.4",
		Ports: []*model.Port{
			{
				Name:     "http-main",
				Port:     80,
				Protocol: model.ProtocolHTTP,
			},
		},
	})

	// Update cache
	server.EnvoyXdsServer.ClearCacheFunc()()

	return server
}

func testPorts(base int) []*model.Port {
	return []*model.Port{
		{
			Name:     "http",
			Port:     base + 80,
			Protocol: model.ProtocolHTTP,
		}, {
			Name:     "http-status",
			Port:     base + 81,
			Protocol: model.ProtocolHTTP,
		}, {
			Name:     "custom",
			Port:     base + 90,
			Protocol: model.ProtocolTCP,
		}, {
			Name:     "mongo",
			Port:     base + 100,
			Protocol: model.ProtocolMongo,
		},
		{
			Name:     "redis",
			Port:     base + 110,
			Protocol: model.ProtocolRedis,
		}, {
			Name:     "h2port",
			Port:     base + 66,
			Protocol: model.ProtocolGRPC,
		}}
}

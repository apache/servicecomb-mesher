package istiov2

import (
	"os"
	"strconv"
	"testing"

	apiv2 "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	apiv2core "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	apiv2endpoint "github.com/envoyproxy/go-control-plane/envoy/api/v2/endpoint"
	istioinfra "github.com/go-mesh/mesher/pkg/infras/istio"

	"github.com/go-chassis/go-chassis/core/registry"
	"github.com/go-chassis/go-chassis/pkg/util/tags"
)

var VaildServiceDiscovery registry.ServiceDiscovery
var AllServices []*registry.MicroService

func TestNewDiscoveryService(t *testing.T) {
	options := registry.Options{
		Addrs:      []string{ValidPilotAddr},
		ConfigPath: KubeConfig,
	}

	// Explicitly set the env vars, though this is checkd in the init of cache_test
	os.Setenv("POD_NAME", TEST_POD_NAME)
	os.Setenv("NAMESPACE", NAMESPACE_DEFAULT)
	os.Setenv("INSTANCE_IP", LocalIPAddress)

	// No panic should happen
	VaildServiceDiscovery = NewDiscoveryService(options)

}

// func TestAutoSync(t *testing.T) {
//     archaius.Init()
//     VaildServiceDiscovery.AutoSync()
// }

func TestEmptyPilotAddrs(t *testing.T) {
	defer func() {
		if err := recover(); err == nil {
			t.Errorf("Panic should be caught")
		}
	}()

	emptyAddrsOptions := registry.Options{
		Addrs:      []string{},
		ConfigPath: KubeConfig,
	}
	NewDiscoveryService(emptyAddrsOptions)
}

func TestGetAllMicroServices(t *testing.T) {
	services, err := VaildServiceDiscovery.GetAllMicroServices()
	if err != nil {
		t.Errorf("Failed to get all micro services: %s", err.Error())
	}

	if len(services) == 0 {
		t.Log("Warn: no micro services found")
	}

}

func TestGetMicroServiceID(t *testing.T) {
	serviceName := "pilotv2server"
	msID, err := VaildServiceDiscovery.GetMicroServiceID("default", serviceName, "v3", "")
	if err != nil {
		t.Errorf("Failed to get micro service id: %s", err.Error())
	}

	if msID != serviceName {
		t.Errorf("In pilotv2 discovery, msID should be equal to serviceName(%s != %s)", msID, serviceName)
	}
}

func TestGetMicroService(t *testing.T) {
	serviceName := "istio-pilot"
	svc, err := VaildServiceDiscovery.GetMicroService(serviceName)
	if err != nil {
		t.Errorf("Failed to get micro service: %s", err.Error())
	}
	if svc == nil {
		t.Errorf("istio-pilot service should not be nil")
	}
}

func TestGetMicroServiceInstance(t *testing.T) {
	// serviceName := "istio-pilot"
	serviceName := "hello"
	instances, err := VaildServiceDiscovery.GetMicroServiceInstances("pilotv2client", serviceName)
	if err != nil {
		t.Errorf("Failed to get micro service instances of istio-pilot: %s", err.Error())
	}
	if len(instances) == 0 {
		t.Errorf("istio-pilot's instances should not be empty")
	}
}

func TestFindMicroServiceInstances(t *testing.T) {
	discovery, ok := VaildServiceDiscovery.(*ServiceDiscovery)
	if !ok {
		t.Errorf("Failed to convert discovery into type istiov2.ServiceDiscovery")
		return
	}
	client := discovery.client

	clusters, err := client.CDS()
	if err != nil {
		t.Errorf("Failed to teset FindMicroServiceInstances, CDS failed: %s", err.Error())
	}

	var clusterWithSubset *istioinfra.XdsClusterInfo = nil
	for _, c := range clusters {
		if info := istioinfra.ParseClusterName(c.Name); info != nil && info.Subset != "" {
			clusterWithSubset = info
		}
	}

	if clusterWithSubset != nil {
		// an empty tags will make sure target tag always match
		emptyTags := utiltags.Tags{
			KV:    map[string]string{},
			Label: "",
		}
		instances, err := VaildServiceDiscovery.FindMicroServiceInstances("pilotv2client", clusterWithSubset.ServiceName, emptyTags)
		if err != nil {
			t.Errorf("Failed to FindMicroServiceInstances of %s: %s", clusterWithSubset.ServiceName, err.Error())
		}
		if len(instances) == 0 {
			t.Logf("%s's service instances is empty\n", clusterWithSubset.ServiceName)
			t.Logf("Pls check if the destinationrule and corresponding pod tags are matching")
		}
	} else if len(clusters) != 0 {
		t.Log("No clusters are with subsets")
		targetCluster := clusters[0]

		tags := utiltags.Tags{
			KV: map[string]string{
				"version": "v1",
			},
			Label: "version=v1",
		}
		_, err := VaildServiceDiscovery.FindMicroServiceInstances("pilotv2client", targetCluster.Name, tags)
		if err == nil {
			t.Errorf("Should caught error to get the endpoints of cluster without tags")
		}
	}

}

func TestToMicroService(t *testing.T) {
	cluster := &apiv2.Cluster{
		Name: "pilotv2server",
	}

	svc := toMicroService(cluster)

	if svc.ServiceID != cluster.Name {
		t.Errorf("service id should be equal to cluster name(%s != %s)", svc.ServiceID, cluster.Name)
	}
}

func TestToMicroServiceInstance(t *testing.T) {
	lbendpoint := &apiv2endpoint.LbEndpoint{
		Endpoint: &apiv2endpoint.Endpoint{
			Address: &apiv2core.Address{
				Address: &apiv2core.Address_SocketAddress{
					SocketAddress: &apiv2core.SocketAddress{
						Address: "192.168.0.10:8822",
					},
				},
			},
		},
	}
	clusterName := "pilotv2server"
	tags := map[string]string{
		"version": "v1",
	}
	msi := toMicroServiceInstance(clusterName, lbendpoint, tags)

	socketAddr := lbendpoint.Endpoint.Address.GetSocketAddress()
	addr := socketAddr.GetAddress()
	port := socketAddr.GetPortValue()

	if msi.InstanceID != addr+"_"+strconv.FormatUint(uint64(port), 10) {
		t.Errorf("Invalid msi.InstanceID: %s should be equal to %s_%d", msi.InstanceID, addr, port)
	}

	if msi.HostName != clusterName {
		t.Errorf("Invalid msi.HostName: %s should be equal to %s", msi.HostName, clusterName)
	}

	// Test if the tags match
	if !tagsMatch(tags, msi.Metadata) {
		t.Errorf("Tags not match, %v should be subset of %s", tags, msi.Metadata)
	}
}

func TestClose(t *testing.T) {
	if err := VaildServiceDiscovery.Close(); err != nil {
		t.Error(err)
	}
}

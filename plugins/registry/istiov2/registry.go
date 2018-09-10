package pilotv2

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-chassis/go-chassis/core/common"
	"github.com/go-chassis/go-chassis/core/lager"
	"github.com/go-chassis/go-chassis/core/metadata"
	"github.com/go-chassis/go-chassis/core/registry"
	"github.com/go-chassis/go-chassis/pkg/util/iputil"
	"github.com/go-chassis/go-chassis/pkg/util/tags"

	apiv2 "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	apiv2endpoint "github.com/envoyproxy/go-control-plane/envoy/api/v2/endpoint"
)

var (
	POD_NAME      string
	POD_NAMESPACE string
	INSTANCE_IP   string
)

const (
	PilotV2Registry = "pilotv2"
)

type ServiceDiscovery struct {
	Name    string
	client  *XdsClient
	options registry.Options
}

func (discovery *ServiceDiscovery) GetMicroServiceID(appID, microServiceName, version, env string) (string, error) {
	return microServiceName, nil
}

func (discovery *ServiceDiscovery) GetAllMicroServices() ([]*registry.MicroService, error) {
	clusters, err := discovery.client.CDS()
	if err != nil {
		return nil, err
	}
	microServices := []*registry.MicroService{}
	for _, cluster := range clusters {
		microServices = append(microServices, toMicroService(&cluster))
	}
	return microServices, nil
}

func toMicroService(cluster *apiv2.Cluster) *registry.MicroService {
	svc := &registry.MicroService{}
	svc.ServiceID = cluster.Name
	svc.ServiceName = cluster.Name
	svc.Version = common.DefaultVersion
	svc.AppID = common.DefaultApp
	svc.Level = "BACK"
	svc.Status = "UP"
	svc.Framework = &registry.Framework{
		Name:    "Istio",
		Version: common.LatestVersion,
	}
	svc.RegisterBy = metadata.PlatformRegistrationComponent

	return svc
}

func toMicroServiceInstance(clusterName string, lbendpoint *apiv2endpoint.LbEndpoint, tags map[string]string) *registry.MicroServiceInstance {
	socketAddress := lbendpoint.Endpoint.Address.GetSocketAddress()
	addr := socketAddress.Address
	port := socketAddress.GetPortValue()
	msi := &registry.MicroServiceInstance{}
	msi.InstanceID = fmt.Sprintf("%s_%d", addr, port)
	msi.HostName = clusterName
	msi.EndpointsMap = map[string]string{
		common.ProtocolRest: fmt.Sprintf("%s:%d", addr, port),
	}
	msi.DefaultEndpoint = fmt.Sprintf("%s:%d", addr, port)
	msi.DefaultProtocol = common.ProtocolRest
	msi.Metadata = tags

	return msi
}

func (discovery *ServiceDiscovery) GetMicroService(microServiceID string) (*registry.MicroService, error) {
	// If the service is in the clusters, return it, or nil

	clusters, err := discovery.client.CDS()
	if err != nil {
		return nil, err
	}

	var targetCluster apiv2.Cluster
	for _, cluster := range clusters {
		parts := strings.Split(cluster.Name, "|")
		if len(parts) < 4 {
			lager.Logger.Warnf("Invalid cluster name: %s", cluster.Name)
			continue
		}

		svcName := parts[3]
		if strings.Index(svcName, microServiceID+".") == 0 {
			targetCluster = cluster
			break
		}
	}

	if &targetCluster == nil {
		return nil, nil
	}

	return toMicroService(&targetCluster), nil
}

func (discovery *ServiceDiscovery) GetMicroServiceInstances(consumerID, providerID string) ([]*registry.MicroServiceInstance, error) {
	// TODO Handle the registry.MicroserviceIndex cache
	// TODO Handle the microServiceName
	service, err := discovery.GetMicroService(providerID)
	if err != nil {
		return nil, err
	}

	loadAssignment, err := discovery.client.EDS(service.ServiceName)
	if err != nil {
		return nil, err
	}

	instances := []*registry.MicroServiceInstance{}
	endpionts := loadAssignment.Endpoints
	for _, item := range endpionts {
		for _, lbendpoint := range item.LbEndpoints {
			msi := toMicroServiceInstance(loadAssignment.ClusterName, &lbendpoint, nil) // The cluster without subset doesn't have tags
			instances = append(instances, msi)
		}
	}

	return instances, nil
}

func (discovery *ServiceDiscovery) FindMicroServiceInstances(consumerID, microServiceName string, tags utiltags.Tags) ([]*registry.MicroServiceInstance, error) {
	instances := simpleCache.GetWithTags(microServiceName, tags.KV)
	if instances == nil {
		var lbendpoints []apiv2endpoint.LbEndpoint
		var err error
		lbendpoints, clusterName, err := discovery.client.GetEndpointsByTags(microServiceName, tags.KV)
		if err != nil {
			return nil, err
		}

		updateInstanceIndexCache(lbendpoints, clusterName, tags.KV)

		instances = simpleCache.GetWithTags(microServiceName, tags.KV)
		if instances == nil {
			return nil, fmt.Errorf("Failed to find microservice instances of %s from cache", microServiceName)
		}
	}
	return instances, nil
}

var cacheManager *CacheManager

func (discovery *ServiceDiscovery) AutoSync() {
	var err error
	cacheManager, err = NewCacheManager(discovery.client)
	if err != nil {
		lager.Logger.Errorf("Failed to create cache manager, indexing will not work: %s", err.Error())
	} else {
		cacheManager.AutoSync()
	}
}

func (discovery *ServiceDiscovery) Close() error {
	// TODO Should we explicitly recycle discovery's other resources?
	// discovery.client.ReqCaches = nil
	return discovery.client.GrpcConn.Close()
}

func NewDiscoveryService(options registry.Options) registry.ServiceDiscovery {
	if len(options.Addrs) == 0 {
		panic("Failed to create discovery service: Address not specified")
	}
	pilotAddr := options.Addrs[0]
	nodeInfo := &NodeInfo{
		PodName:    POD_NAME,
		Namespace:  POD_NAMESPACE,
		InstanceIP: INSTANCE_IP,
	}
	xdsClient, err := NewXdsClient(pilotAddr, options.TLSConfig, nodeInfo, options.ConfigPath)
	if err != nil {
		panic("Failed to create XDS client: " + err.Error())
	}

	discovery := &ServiceDiscovery{
		client:  xdsClient,
		Name:    PilotV2Registry,
		options: options,
	}

	return discovery
}

func init() {
	// Init the node info
	POD_NAME = os.Getenv("POD_NAME")
	POD_NAMESPACE = os.Getenv("POD_NAMESPACE")
	INSTANCE_IP = os.Getenv("INSTANCE_IP")

	// TODO Handle the default value
	if POD_NAME == "" {
		POD_NAME = "pod_name_default"
	}
	if POD_NAMESPACE == "" {
		POD_NAMESPACE = "default"
	}
	if INSTANCE_IP == "" {
		fmt.Println("[WARN] Env var INSTANCE_IP not set, try to get instance ip from local network, the service might not work properly.")
		INSTANCE_IP = iputil.GetLocalIP()
		if INSTANCE_IP == "" {
			// Won't work without instance ip
			panic("Failed to get instance ip")
		}
	}

	registry.InstallServiceDiscovery(PilotV2Registry, NewDiscoveryService)
}

/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package istiov2

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	istioinfra "github.com/apache/servicecomb-mesher/proxy/pkg/infras/istio"
	apiv2 "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	apiv2endpoint "github.com/envoyproxy/go-control-plane/envoy/api/v2/endpoint"

	"github.com/go-chassis/go-chassis/core/common"
	"github.com/go-chassis/go-chassis/core/metadata"
	"github.com/go-chassis/go-chassis/core/registry"
	"github.com/go-chassis/go-chassis/pkg/util/iputil"
	"github.com/go-chassis/go-chassis/pkg/util/tags"
	"github.com/go-mesh/openlogging"
)

var (
	//PodName is the name of the pod that mesher runs in
	PodName string
	//PodNamespace is the namespace which the pod belongs to
	PodNamespace string
	//InstanceIP is the IP of the pod(the IP of the first network adaptor)
	InstanceIP string
)

const (
	PilotV2Registry = "pilotv2"
)

//ServiceDiscovery is the discovery service for istio pilot with xDS v2 API
type ServiceDiscovery struct {
	Name    string
	client  *istioinfra.XdsClient
	options registry.Options
}

//GetMicroServiceID returns the id of the micro service
func (discovery *ServiceDiscovery) GetMicroServiceID(appID, microServiceName, version, env string) (string, error) {
	return microServiceName, nil
}

//GetAllMicroServices returns all the micro services, which is mapped from xDS clusters
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
	portStr := strconv.FormatUint(uint64(port), 10)
	msi := &registry.MicroServiceInstance{}
	msi.InstanceID = addr + "_" + portStr
	msi.HostName = clusterName
	msi.DefaultEndpoint = addr + ":" + portStr
	msi.EndpointsMap = map[string]*registry.Endpoint{
		common.ProtocolRest: {Address: msi.DefaultEndpoint},
	}
	msi.DefaultProtocol = common.ProtocolRest
	msi.Metadata = tags

	return msi
}

//GetMicroService returns the micro service info
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
			openlogging.GetLogger().Warnf("Invalid cluster name: %s", cluster.Name)
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

//GetMicroServiceInstances returns the instances of the micro service
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

//FindMicroServiceInstances returns the micro service's instances filtered with tags
func (discovery *ServiceDiscovery) FindMicroServiceInstances(consumerID, microServiceName string, tags utiltags.Tags) ([]*registry.MicroServiceInstance, error) {
	if tags.KV == nil || tags.Label == "" { // Chassis might pass an empty tags
		return discovery.GetMicroServiceInstances(consumerID, microServiceName)
	}

	instances := simpleCache.GetWithTags(microServiceName, tags.KV)
	if len(instances) == 0 {
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

//AutoSync updates the services' info periodically in the background
func (discovery *ServiceDiscovery) AutoSync() {
	var err error
	cacheManager, err = NewCacheManager(discovery.client)
	if err != nil {
		openlogging.GetLogger().Errorf("Failed to create cache manager, indexing will not work: %s", err.Error())
	} else {
		cacheManager.AutoSync()
	}
}

//Close closes the discovery service
func (discovery *ServiceDiscovery) Close() error {
	return nil
}

//NewDiscoveryService creates the new ServiceDiscovery instance
func NewDiscoveryService(options registry.Options) registry.ServiceDiscovery {
	if len(options.Addrs) == 0 {
		panic("Failed to create discovery service: Address not specified")
	}
	pilotAddr := options.Addrs[0]
	nodeInfo := &istioinfra.NodeInfo{
		PodName:    PodName,
		Namespace:  PodNamespace,
		InstanceIP: InstanceIP,
	}
	xdsClient, err := istioinfra.NewXdsClient(pilotAddr, options.TLSConfig, nodeInfo, options.ConfigPath)
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
	PodName = os.Getenv("POD_NAME")
	PodNamespace = os.Getenv("POD_NAMESPACE")
	InstanceIP = os.Getenv("INSTANCE_IP")

	// TODO Handle the default value
	if PodName == "" {
		PodName = "pod_name_default"
	}
	if PodNamespace == "" {
		PodNamespace = "default"
	}
	if InstanceIP == "" {
		log.Println("[WARN] Env var INSTANCE_IP not set, try to get instance ip from local network, the service might not work properly.")
		InstanceIP = iputil.GetLocalIP()
		if InstanceIP == "" {
			// Won't work without instance ip
			panic("Failed to get instance ip")
		}
	}

	registry.InstallServiceDiscovery(PilotV2Registry, NewDiscoveryService)
}

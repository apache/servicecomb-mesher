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
	"strconv"
	"strings"
	"sync"
	"time"

	istioinfra "github.com/apache/servicecomb-mesher/proxy/pkg/infras/istio"
	apiv2endpoint "github.com/envoyproxy/go-control-plane/envoy/api/v2/endpoint"

	"github.com/go-chassis/go-archaius"
	"github.com/go-chassis/go-chassis/core/common"
	"github.com/go-chassis/go-chassis/core/config"
	"github.com/go-chassis/go-chassis/core/lager"
	"github.com/go-chassis/go-chassis/core/registry"
)

const (
	DefaultRefreshInterval = time.Second * 30
)

var simpleCache *EndpointCache

func init() {
	simpleCache = &EndpointCache{
		cache: map[string]EndpointSubset{},
	}
}

//CacheManager manages the caches for istio pilot results.
type CacheManager struct {
	xdsClient *istioinfra.XdsClient
}

//AutoSync fetches and updates the cluster and endpoint info periodically
func (cm *CacheManager) AutoSync() {
	cm.refreshCache()

	var ticker *time.Ticker
	refreshInterval := config.GetServiceDiscoveryRefreshInterval()
	if refreshInterval == "" {
		ticker = time.NewTicker(DefaultRefreshInterval)
	} else {
		timeValue, err := time.ParseDuration(refreshInterval)
		if err != nil {
			lager.Logger.Errorf("refeshInterval is invalid. So use Default value: %s", err.Error())
			timeValue = DefaultRefreshInterval
		}

		ticker = time.NewTicker(timeValue)
	}
	go func() {
		for range ticker.C {
			cm.refreshCache()
		}
	}()
}

func (cm *CacheManager) refreshCache() {
	// TODO What is the design of autodiscovery
	if archaius.GetBool("cse.service.registry.autodiscovery", false) {
		lager.Logger.Errorf("SyncPilotEndpoints failed: not supported")
	}

	err := cm.pullMicroserviceInstance()
	if err != nil {
		lager.Logger.Errorf("AutoUpdateMicroserviceInstance failed: %s", err.Error())
	}

	if archaius.GetBool("cse.service.registry.autoSchemaIndex", false) {
		lager.Logger.Errorf("MakeSchemaIndex failed: Not support operation")
	}

	if archaius.GetBool("cse.service.registry.autoIPIndex", false) {
		err = cm.MakeIPIndex()
		if err != nil {
			lager.Logger.Errorf("Auto Update IP index failed: %s", err.Error())
		}
	}
}

func (cm *CacheManager) pullMicroserviceInstance() error {
	clusterInfos, err := cm.getClusterInfos()
	if err != nil {
		return err
	}

	for _, clusterInfo := range clusterInfos {
		if clusterInfo.Subset != "" {
			// Update the cache
			instances := []*registry.MicroServiceInstance{}
			for _, addr := range clusterInfo.Addrs {
				msi := &registry.MicroServiceInstance{}
				msi.InstanceID = strings.Replace(addr, ":", "_", 0)
				msi.HostName = clusterInfo.ClusterName
				msi.EndpointsMap = map[string]*registry.Endpoint{
					common.ProtocolRest: {Address: addr},
				}
				msi.DefaultEndpoint = addr
				msi.DefaultProtocol = common.ProtocolRest
				msi.Metadata = clusterInfo.Tags

				instances = append(instances, msi)
			}

			endpointSubset := EndpointSubset{
				tags:       clusterInfo.Tags,
				instances:  instances,
				subsetName: clusterInfo.Subset,
			}
			simpleCache.Set(clusterInfo.ClusterName, endpointSubset)
		}
	}

	return nil
}

//MakeIPIndex caches the cluster info with address as the key
func (cm *CacheManager) MakeIPIndex() error {
	// TODO Use getClusterInfo to replace the logic
	clusterInfos, err := cm.getClusterInfos()
	if err != nil {
		return err
	}

	for _, clusterInfo := range clusterInfos {
		for _, addr := range clusterInfo.Addrs {
			si := &registry.SourceInfo{}
			// TODO Get tags by subset and put them into si.Tags
			si.Name = clusterInfo.ClusterName
			si.Tags = clusterInfo.Tags
			registry.SetIPIndex(addr, si)
			// TODO Why don't we have to index every endpoint?
			// break
		}
	}

	return nil
}

//NewCacheManager creates the CacheManager instance.
func NewCacheManager(xdsClient *istioinfra.XdsClient) (*CacheManager, error) {
	cacheManager := &CacheManager{
		xdsClient: xdsClient,
	}

	return cacheManager, nil
}

func (cm *CacheManager) getClusterInfos() ([]istioinfra.XdsClusterInfo, error) {
	clusterInfos := []istioinfra.XdsClusterInfo{}

	clusters, err := cm.xdsClient.CDS()
	if err != nil {
		return nil, err
	}

	for _, cluster := range clusters {
		// xDS v2 API: CDS won't obtain the cluster's endpoints, call EDS to get the endpoints

		clusterInfo := istioinfra.ParseClusterName(cluster.Name)
		if clusterInfo == nil {
			continue
		}

		// Get Tags
		if clusterInfo.Subset != "" { // Only clusters with subset contain labels
			if tags, err := cm.xdsClient.GetSubsetTags(clusterInfo.Namespace, clusterInfo.ServiceName, clusterInfo.Subset); err == nil {
				clusterInfo.Tags = tags
			}
		}

		// Get cluster instances' addresses
		loadAssignment, err := cm.xdsClient.EDS(cluster.Name)
		if err != nil {
			return nil, err
		}
		endpoints := loadAssignment.Endpoints
		for _, endpoint := range endpoints {
			for _, lbendpoint := range endpoint.LbEndpoints {
				socketAddress := lbendpoint.Endpoint.Address.GetSocketAddress()
				addr := socketAddress.GetAddress()
				port := socketAddress.GetPortValue()
				ipAddr := addr + ":" + strconv.FormatUint(uint64(port), 10)
				clusterInfo.Addrs = append(clusterInfo.Addrs, ipAddr)
			}
		}

		clusterInfos = append(clusterInfos, *clusterInfo)
	}
	return clusterInfos, nil
}

// TODO Cache with registry index cache
func updateInstanceIndexCache(lbendpoints []apiv2endpoint.LbEndpoint, clusterName string, tags map[string]string) {
	if len(lbendpoints) == 0 {
		simpleCache.Delete(clusterName)
		return
	}

	instances := make([]*registry.MicroServiceInstance, 0, len(lbendpoints))
	for _, lbendpoint := range lbendpoints {
		msi := toMicroServiceInstance(clusterName, &lbendpoint, tags)
		instances = append(instances, msi)
	}
	subset := EndpointSubset{
		tags:      tags,
		instances: instances,
	}

	info := istioinfra.ParseClusterName(clusterName)
	if info != nil {
		subset.subsetName = info.Subset
	}

	simpleCache.Set(clusterName, subset)
}

//EndpointCache caches the clusters' endpoint and tags
type EndpointCache struct {
	mux   sync.Mutex
	cache map[string]EndpointSubset
}

//EndpointSubset stores the tags and instances of a service
type EndpointSubset struct {
	subsetName string
	tags       map[string]string
	instances  []*registry.MicroServiceInstance
}

//Delete removes the cached instances of the specified cluster
func (c *EndpointCache) Delete(clusterName string) {
	c.mux.Lock()
	delete(c.cache, clusterName)
	c.mux.Unlock()
}

//Set updates the cluster's instance info
func (c *EndpointCache) Set(clusterName string, subset EndpointSubset) {
	c.mux.Lock()
	c.cache[clusterName] = subset
	c.mux.Unlock()
}

//GetWithTags returns the instances of the service, filtered with tags
func (c *EndpointCache) GetWithTags(serviceName string, tags map[string]string) []*registry.MicroServiceInstance {
	// Get subsets whose clusterName matches the service name
	matchedSubsets := []EndpointSubset{}
	c.mux.Lock()
	for clusterName, subset := range c.cache {
		info := istioinfra.ParseClusterName(clusterName)
		if info != nil && info.ServiceName == serviceName {
			matchedSubsets = append(matchedSubsets, subset)
		}
	}
	c.mux.Unlock()

	if len(matchedSubsets) == 0 {
		return nil
	}

	var instances []*registry.MicroServiceInstance

	for _, subset := range matchedSubsets {
		if tagsMatch(subset.tags, tags) {
			instances = subset.instances
			break
		}

	}
	return instances
}

// TODO There might be some utils in go-chassis doing the same thing
func tagsMatch(tags, targetTags map[string]string) bool {
	matched := true
	for k, v := range targetTags {
		if val, exists := tags[k]; !exists || val != v {
			matched = false
			break
		}
	}
	return matched
}

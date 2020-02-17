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
	"os"
	"strconv"
	"testing"

	istioinfra "github.com/apache/servicecomb-mesher/proxy/pkg/infras/istio"
	apiv2 "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	apiv2core "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	apiv2endpoint "github.com/envoyproxy/go-control-plane/envoy/api/v2/endpoint"

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

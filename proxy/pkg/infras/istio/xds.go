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

package pilotv2

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"strings"

	k8sinfra "github.com/apache/servicecomb-mesher/proxy/pkg/infras/k8s"
	apiv2 "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	apiv2core "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	apiv2endpoint "github.com/envoyproxy/go-control-plane/envoy/api/v2/endpoint"
	apiv2route "github.com/envoyproxy/go-control-plane/envoy/api/v2/route"

	"github.com/envoyproxy/go-control-plane/envoy/service/discovery/v2"
	"github.com/go-mesh/openlogging"
	"github.com/gogo/protobuf/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"k8s.io/client-go/rest"
)

//XdsClient provides the XDS API calls.
type XdsClient struct {
	PilotAddr   string
	TlsConfig   *tls.Config
	ReqCaches   map[XdsType]*XdsReqCache
	nodeInfo    *NodeInfo
	NodeID      string
	NodeCluster string
	k8sClient   *rest.RESTClient
}

//XdsType is the wrapper of string, the wrapper type should be "cds", "eds", "lds" or "rds"
type XdsType string

const (
	TypeCds XdsType = "cds"
	TypeEds XdsType = "eds"
	TypeLds XdsType = "lds"
	TypeRds XdsType = "rds"
)

//XdsReqCache stores the VersionInfo and Nonce for the XDS calls
type XdsReqCache struct {
	Nonce       string
	VersionInfo string
}

//NodeInfo stores the info of the node, which will be used to make a
//XDS call
type NodeInfo struct {
	PodName    string
	Namespace  string
	InstanceIP string
}

//XdsClusterInfo stores all the infos from a cluster name, which is in
//the format direction|port|subset|hostName
type XdsClusterInfo struct {
	ClusterName  string
	Direction    string
	Port         string
	Subset       string
	HostName     string
	ServiceName  string
	Namespace    string
	DomainSuffix string // DomainSuffix might not be used
	Tags         map[string]string
	Addrs        []string // The accessible addresses of the endpoints
}

//NewXdsClient returns the new XDS client.
func NewXdsClient(pilotAddr string, tlsConfig *tls.Config, nodeInfo *NodeInfo, kubeconfigPath string) (*XdsClient, error) {
	// TODO Handle the array
	xdsClient := &XdsClient{
		PilotAddr: pilotAddr,
		nodeInfo:  nodeInfo,
	}
	xdsClient.NodeID = "sidecar~" + nodeInfo.InstanceIP + "~" + nodeInfo.PodName + "~" + nodeInfo.Namespace
	xdsClient.NodeCluster = nodeInfo.PodName

	xdsClient.ReqCaches = map[XdsType]*XdsReqCache{
		TypeCds: {},
		TypeEds: {},
		TypeLds: {},
		TypeRds: {},
	}

	if k8sClient, err := k8sinfra.CreateK8SRestClient(kubeconfigPath, "apis", "networking.istio.io", "v1alpha3"); err != nil {
		return nil, err
	} else {
		xdsClient.k8sClient = k8sClient
	}

	return xdsClient, nil
}

//GetSubsetTags returns the tags of the specified subset.
func (client *XdsClient) GetSubsetTags(namespace, hostName, subsetName string) (map[string]string, error) {
	req := client.k8sClient.Get()
	req.Resource("destinationrules")
	req.Namespace(namespace)

	result := req.Do()
	rawBody, err := result.Raw()
	if err != nil {
		return nil, err
	}

	var drResult k8sinfra.DestinationRuleResult
	if err := json.Unmarshal(rawBody, &drResult); err != nil {
		return nil, err
	}

	// Find the subset
	tags := map[string]string{}
	for _, dr := range drResult.Items {
		if dr.Spec.Host == hostName {
			for _, subset := range dr.Spec.Subsets {
				if subset.Name == subsetName {
					for k, v := range subset.Labels {
						tags[k] = v
					}
					break
				}
			}
			break
		}
	}

	return tags, nil
}

func (client *XdsClient) getGrpcConn() (*grpc.ClientConn, error) {
	var conn *grpc.ClientConn
	var err error
	if client.TlsConfig != nil {
		creds := credentials.NewTLS(client.TlsConfig)
		conn, err = grpc.Dial(client.PilotAddr, grpc.WithTransportCredentials(creds))
	} else {
		conn, err = grpc.Dial(client.PilotAddr, grpc.WithInsecure())
	}

	return conn, err
}

func getAdsResClient(client *XdsClient) (v2.AggregatedDiscoveryService_StreamAggregatedResourcesClient, *grpc.ClientConn, error) {
	conn, err := client.getGrpcConn()
	if err != nil {
		return nil, nil, err
	}

	adsClient := v2.NewAggregatedDiscoveryServiceClient(conn)
	adsResClient, err := adsClient.StreamAggregatedResources(context.Background())
	if err != nil {
		return nil, nil, err
	}

	return adsResClient, conn, nil
}

func (client *XdsClient) getRouterClusters(clusterName string) ([]string, error) {
	virtualHosts, err := client.RDS(clusterName)
	if err != nil {
		return nil, err
	}

	routerClusters := []string{}
	for _, h := range virtualHosts {
		for _, r := range h.Routes {
			routerClusters = append(routerClusters, r.GetRoute().GetCluster())
		}
	}

	return routerClusters, nil
}

func (client *XdsClient) getVersionInfo(resType XdsType) string {
	return client.ReqCaches[resType].VersionInfo
}
func (client *XdsClient) getNonce(resType XdsType) string {
	return client.ReqCaches[resType].Nonce
}

func (client *XdsClient) setVersionInfo(resType XdsType, versionInfo string) {
	client.ReqCaches[resType].VersionInfo = versionInfo
}

func (client *XdsClient) setNonce(resType XdsType, nonce string) {
	client.ReqCaches[resType].Nonce = nonce
}

//CDS s the Clsuter Discovery Service API, which fetches all the clusters from istio pilot
func (client *XdsClient) CDS() ([]apiv2.Cluster, error) {
	adsResClient, conn, err := getAdsResClient(client)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	req := &apiv2.DiscoveryRequest{
		TypeUrl:       "type.googleapis.com/envoy.api.v2.Cluster",
		VersionInfo:   client.getVersionInfo(TypeCds),
		ResponseNonce: client.getNonce(TypeCds),
	}
	req.Node = &apiv2core.Node{
		// Sample taken from istio: router~172.30.77.6~istio-egressgateway-84b4d947cd-rqt45.istio-system~istio-system.svc.cluster.local-2
		// The Node.Id should be in format {nodeType}~{ipAddr}~{serviceId~{domain}, splitted by '~'
		// The format is required by pilot
		Id:      client.NodeID,
		Cluster: client.NodeCluster,
	}

	if err := adsResClient.Send(req); err != nil {
		return nil, err
	}

	resp, err := adsResClient.Recv()
	if err != nil {
		return nil, err
	}

	client.setNonce(TypeCds, resp.GetNonce())
	client.setVersionInfo(TypeCds, resp.GetVersionInfo())
	resources := resp.GetResources()

	var cluster apiv2.Cluster
	clusters := []apiv2.Cluster{}
	for _, res := range resources {
		if err := proto.Unmarshal(res.GetValue(), &cluster); err != nil {
			openlogging.GetLogger().Warnf("Failed to unmarshal cluster resource: %s", err.Error())
		} else {
			clusters = append(clusters, cluster)
		}
	}
	return clusters, nil
}

//EDS is the Endpoint Discovery Service API, the API takes the cluster's name and return all its endpoints(which provide address and port)
func (client *XdsClient) EDS(clusterName string) (*apiv2.ClusterLoadAssignment, error) {
	adsResClient, conn, err := getAdsResClient(client)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	req := &apiv2.DiscoveryRequest{
		TypeUrl:       "type.googleapis.com/envoy.api.v2.ClusterLoadAssignment",
		VersionInfo:   client.getVersionInfo(TypeEds),
		ResponseNonce: client.getNonce(TypeEds),
	}

	req.Node = &apiv2core.Node{
		Id:      client.NodeID,
		Cluster: client.NodeCluster,
	}
	req.ResourceNames = []string{clusterName}
	if err := adsResClient.Send(req); err != nil {
		return nil, err
	}

	resp, err := adsResClient.Recv()
	if err != nil {
		return nil, err
	}

	resources := resp.GetResources()
	client.setNonce(TypeEds, resp.GetNonce())
	client.setVersionInfo(TypeEds, resp.GetVersionInfo())

	var loadAssignment apiv2.ClusterLoadAssignment
	var e error
	// endpoints := []apiv2.ClusterLoadAssignment{}

	for _, res := range resources {
		if err := proto.Unmarshal(res.GetValue(), &loadAssignment); err != nil {
			e = err
		} else {
			// The cluster's LoadAssignment will always be ONE, with Endpoints as its field
			break
		}
	}
	return &loadAssignment, e
}

//GetEndpointsByTags fetches the cluster's endpoints with tags. The tags is usually specified in a DestinationRule.
func (client *XdsClient) GetEndpointsByTags(serviceName string, tags map[string]string) ([]apiv2endpoint.LbEndpoint, string, error) {
	clusters, err := client.CDS()
	if err != nil {
		return nil, "", err
	}

	lbendpoints := []apiv2endpoint.LbEndpoint{}
	clusterName := ""
	for _, cluster := range clusters {
		clusterInfo := ParseClusterName(cluster.Name)
		if clusterInfo == nil || clusterInfo.Subset == "" || clusterInfo.ServiceName != serviceName {
			continue
		}
		// So clusterInfo is not nil and subset is not empty
		if subsetTags, err := client.GetSubsetTags(clusterInfo.Namespace, clusterInfo.ServiceName, clusterInfo.Subset); err == nil {
			// filter with tags
			matched := true
			for k, v := range tags {
				if subsetTagValue, exists := subsetTags[k]; exists == false || subsetTagValue != v {
					matched = false
					break
				}
			}

			if matched { // We got the cluster!
				clusterName = cluster.Name
				loadAssignment, err := client.EDS(cluster.Name)
				if err != nil {
					return nil, clusterName, err
				}

				for _, item := range loadAssignment.Endpoints {
					lbendpoints = append(lbendpoints, item.LbEndpoints...)
				}

				return lbendpoints, clusterName, nil
			}
		}
	}

	return lbendpoints, clusterName, nil
}

//RDS is the Router Discovery Service API, it returns the virtual hosts which contains Routes
func (client *XdsClient) RDS(clusterName string) ([]apiv2route.VirtualHost, error) {
	clusterInfo := ParseClusterName(clusterName)
	if clusterInfo == nil {
		return nil, fmt.Errorf("Invalid clusterName for routers: %s", clusterName)
	}

	adsResClient, conn, err := getAdsResClient(client)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	req := &apiv2.DiscoveryRequest{
		TypeUrl:       "type.googleapis.com/envoy.api.v2.RouteConfiguration",
		VersionInfo:   client.getVersionInfo(TypeRds),
		ResponseNonce: client.getNonce(TypeRds),
	}

	req.Node = &apiv2core.Node{
		Id:      client.NodeID,
		Cluster: client.NodeCluster,
	}
	req.ResourceNames = []string{clusterName}
	if err := adsResClient.Send(req); err != nil {
		return nil, err
	}

	resp, err := adsResClient.Recv()
	if err != nil {
		return nil, err
	}

	resources := resp.GetResources()
	client.setNonce(TypeRds, resp.GetNonce())
	client.setVersionInfo(TypeRds, resp.GetVersionInfo())

	var route apiv2.RouteConfiguration
	virtualHosts := []apiv2route.VirtualHost{}

	for _, res := range resources {
		if err := proto.Unmarshal(res.GetValue(), &route); err != nil {
			openlogging.GetLogger().Warnf("Failed to unmarshal router resource: ", err.Error())
		} else {
			vhosts := route.GetVirtualHosts()
			for _, vhost := range vhosts {
				if vhost.Name == clusterInfo.ServiceName+":"+clusterInfo.Port {
					virtualHosts = append(virtualHosts, vhost)
				}
			}
		}
	}
	return virtualHosts, nil
}

//LDS is the Listener Discovery Service API, which returns all the listerns
func (client *XdsClient) LDS() ([]apiv2.Listener, error) {
	adsResClient, conn, err := getAdsResClient(client)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	req := &apiv2.DiscoveryRequest{
		TypeUrl:       "type.googleapis.com/envoy.api.v2.Listener",
		VersionInfo:   client.getVersionInfo(TypeLds),
		ResponseNonce: client.getNonce(TypeLds),
	}

	req.Node = &apiv2core.Node{
		Id:      client.NodeID,
		Cluster: client.NodeCluster,
	}
	if err := adsResClient.Send(req); err != nil {
		return nil, err
	}

	resp, err := adsResClient.Recv()
	if err != nil {
		return nil, err
	}

	resources := resp.GetResources()
	client.setNonce(TypeLds, resp.GetNonce())
	client.setVersionInfo(TypeLds, resp.GetVersionInfo())

	var listener apiv2.Listener
	listeners := []apiv2.Listener{}

	for _, res := range resources {
		if err := proto.Unmarshal(res.GetValue(), &listener); err != nil {
			openlogging.GetLogger().Warnf("Failed to unmarshal listener resource: ", err.Error())
		} else {
			listeners = append(listeners, listener)
		}
	}
	return listeners, nil
}

//ParseClusterName parse the cluster's name, which is in the format direction|port|subset|hostName, the 4 items will be parsed into different fields. The hostName item will also be parsed into ServcieName, Namespace etc.
func ParseClusterName(clusterName string) *XdsClusterInfo {
	// clusterName format: direction|port|subset|hostName
	// hostName format: |svc.namespace.svc.cluster.local

	parts := strings.Split(clusterName, "|")
	if len(parts) != 4 {
		return nil
	}

	hostnameParts := strings.Split(parts[3], ".")
	if len(hostnameParts) < 2 {
		return nil
	}

	cluster := &XdsClusterInfo{
		Direction:    parts[0],
		Port:         parts[1],
		Subset:       parts[2],
		HostName:     parts[3],
		ServiceName:  hostnameParts[0],
		Namespace:    hostnameParts[1],
		DomainSuffix: strings.Join(hostnameParts[2:], "."),
		ClusterName:  clusterName,
	}

	return cluster
}

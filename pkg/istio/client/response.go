package client

import (
	"errors"

	xdsapi "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	"github.com/go-mesh/mesher/pkg/istio/util"
)

// GetClusterConfiguration returns cluster information from discovery response
func GetClusterConfiguration(res *xdsapi.DiscoveryResponse) ([]xdsapi.Cluster, error) {
	if res.TypeUrl != util.ClusterType {
		return nil, errors.New("Invalid typeURL" + res.TypeUrl)
	}

	var cluster []xdsapi.Cluster
	for _, value := range res.GetResources() {
		cla := &xdsapi.Cluster{}
		err := cla.Unmarshal(value.Value)
		if err != nil {
			return nil, errors.New("unmarshall error")

		}
		cluster = append(cluster, *cla)

	}
	return cluster, nil
}

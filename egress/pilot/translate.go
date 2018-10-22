package pilot

import (
	envoy_api "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	egressmodel "github.com/go-mesh/mesher/config/model"

	"strconv"
	"strings"
)

func clusterToEgressRule(clusters []envoy_api.Cluster) []*egressmodel.EgressRule {
	var egressRules []*egressmodel.EgressRule
	for _, cluster := range clusters {
		var rule egressmodel.EgressRule
		if cluster.Type == envoy_api.Cluster_ORIGINAL_DST {

			data := strings.Split(cluster.Name, "|")

			if len(data) > 0 && data[0] == OUTBOUND {
				newdata := strings.Split(data[1], "||")
				intport, err := strconv.Atoi(newdata[0])
				if err != nil {
					return nil
				}

				rule.Hosts = []string{newdata[1]}
				rule.Ports = []*egressmodel.EgressPort{
					{Port: int32(intport),
						Protocol: "http"}}
			}

		}
		egressRules = append(egressRules, &rule)

	}
	return egressRules

}

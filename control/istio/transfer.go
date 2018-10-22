package istio

import (
	"github.com/go-chassis/go-chassis/control"
	egressmodel "github.com/go-mesh/mesher/config/model"
)

func SaveToEgressCache(egressConfigFromPilot map[string][]*egressmodel.EgressRule) {
	{
		var egressconfig []control.EgressConfig
		for _, v := range egressConfigFromPilot {
			var Ports []*control.EgressPort
			for _, v1 := range v {
				for _, v2 := range v1.Ports {
					p := control.EgressPort{
						Port:     (*v2).Port,
						Protocol: (*v2).Protocol,
					}
					Ports = append(Ports, &p)
				}
				c := control.EgressConfig{
					Hosts: v1.Hosts,
					Ports: Ports,
				}
				egressconfig = append(egressconfig, c)
			}
		}
		EgressConfigCache.Set("", egressconfig, 0)
	}
}

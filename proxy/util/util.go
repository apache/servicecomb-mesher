package util

import (
	"fmt"
	"github.com/go-chassis/go-chassis/core/invocation"
	"github.com/go-mesh/mesher/proxy/cmd"
	"github.com/go-mesh/mesher/proxy/common"
)

//SetLocalServiceAddress assign invocation endpoint a local service address
//header "X-Forwarded-Port" has highest priority
// if it is empty
// it uses ports config in cmd param or env
func SetLocalServiceAddress(inv *invocation.Invocation, port string) error {
	inv.Endpoint = cmd.Configs.PortsMap[inv.Protocol]
	if port == "" {
		inv.Endpoint = cmd.Configs.PortsMap[inv.Protocol]
		if inv.Endpoint == "" {
			return fmt.Errorf("[%s] is not supported, [%s] didn't set env [%s] or cmd parameter --service-ports before mesher start",
				inv.Protocol, inv.MicroServiceName, common.EnvServicePorts)
		}
		return nil
	}
	inv.Endpoint = "127.0.0.1:" + port
	return nil
}

package main

import (
	_ "net/http/pprof"

	_ "github.com/go-mesh/mesher/proxy/resolver/authority"

	_ "github.com/go-mesh/mesher/proxy/handler"
	_ "github.com/go-mesh/mesher/proxy/protocol/dubbo/client/chassis"
	_ "github.com/go-mesh/mesher/proxy/protocol/dubbo/server"
	_ "github.com/go-mesh/mesher/proxy/protocol/dubbo/simpleRegistry"

	_ "github.com/go-chassis/go-chassis/configcenter" //use config center
	//protocols
	_ "github.com/go-mesh/mesher/proxy/protocol/grpc"
	_ "github.com/go-mesh/mesher/proxy/protocol/http"

	"github.com/go-mesh/mesher/proxy/server"

	_ "github.com/go-mesh/mesher/proxy/pkg/egress/archaius"
	_ "github.com/go-mesh/mesher/proxy/pkg/egress/pilot"

	_ "github.com/go-mesh/mesher/proxy/control/istio"
)

func main() {
	server.Run()
}

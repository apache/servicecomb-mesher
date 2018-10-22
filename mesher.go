package main

import (
	_ "net/http/pprof"

	_ "github.com/go-mesh/mesher/resolver/authority"

	_ "github.com/go-mesh/mesher/handler"
	_ "github.com/go-mesh/mesher/protocol/dubbo/client/chassis"
	_ "github.com/go-mesh/mesher/protocol/dubbo/server"
	_ "github.com/go-mesh/mesher/protocol/dubbo/simpleRegistry"

	_ "github.com/go-chassis/go-chassis/configcenter" //use config center
	//protocols
	_ "github.com/go-mesh/mesher/protocol/grpc"
	_ "github.com/go-mesh/mesher/protocol/http"

	"github.com/go-mesh/mesher/server"

	_ "github.com/go-mesh/mesher/egress/chassis"
	_ "github.com/go-mesh/mesher/egress/pilot"

	_ "github.com/go-mesh/mesher/control/archiaus"
	_ "github.com/go-mesh/mesher/control/istio"
)

func main() {
	server.Run()
}

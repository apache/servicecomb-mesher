package main

import (
	_ "net/http/pprof"

	_ "github.com/go-mesh/mesher/resolver/authority"

	_ "github.com/go-mesh/mesher/handler"
	_ "github.com/go-mesh/mesher/protocol/dubbo/client/chassis"
	_ "github.com/go-mesh/mesher/protocol/dubbo/server"
	_ "github.com/go-mesh/mesher/protocol/dubbo/simpleRegistry"

	_ "github.com/go-chassis/go-chassis/config-center" //use config center
	//protocols
	_ "github.com/go-mesh/mesher/protocol/grpc"
	_ "github.com/go-mesh/mesher/protocol/http"

	"github.com/go-mesh/mesher/server"
)

func main() {
	server.Run()
}

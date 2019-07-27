package main

import (
	_ "net/http/pprof"

	_ "github.com/apache/servicecomb-mesher/proxy/resolver/authority"

	_ "github.com/apache/servicecomb-mesher/proxy/handler"
	_ "github.com/apache/servicecomb-mesher/proxy/protocol/dubbo/client/chassis"
	_ "github.com/apache/servicecomb-mesher/proxy/protocol/dubbo/server"
	_ "github.com/apache/servicecomb-mesher/proxy/protocol/dubbo/simpleRegistry"

	_ "github.com/go-chassis/go-chassis/configcenter" //use config center
	//protocols
	_ "github.com/apache/servicecomb-mesher/proxy/protocol/grpc"
	_ "github.com/apache/servicecomb-mesher/proxy/protocol/http"

	"github.com/apache/servicecomb-mesher/proxy/server"

	_ "github.com/apache/servicecomb-mesher/proxy/pkg/egress/archaius"
	_ "github.com/apache/servicecomb-mesher/proxy/pkg/egress/pilot"

	_ "github.com/apache/servicecomb-mesher/proxy/control/istio"
)

func main() {
	server.Run()
}

module github.com/apache/servicecomb-mesher

require (
	github.com/apache/servicecomb-kie v0.1.1-0.20191119112752-d564a5ac0115
	github.com/envoyproxy/go-control-plane v0.6.0
	github.com/ghodss/yaml v1.0.0
	github.com/go-chassis/foundation v0.1.1-0.20191113114104-2b05871e9ec4
	github.com/go-chassis/go-archaius v1.0.0
	github.com/go-chassis/go-chassis v1.8.1-0.20191213091237-f7ba7d0c8fe5
	github.com/go-chassis/gohessian v0.0.0-20180702061429-e5130c25af55
	github.com/go-mesh/openlogging v1.0.1
	github.com/gogo/googleapis v1.3.0 // indirect
	github.com/gogo/protobuf v1.3.0
	github.com/gopherjs/gopherjs v0.0.0-20181103185306-d547d1d9531e // indirect
	github.com/imdario/mergo v0.3.7 // indirect
	github.com/lyft/protoc-gen-validate v0.1.0 // indirect
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/prometheus/client_golang v0.9.1
	github.com/prometheus/client_model v0.0.0-20190115171406-56726106282f
	github.com/stretchr/testify v1.4.0
	github.com/tetratelabs/go2sky v0.1.1-0.20190703154722-1eaab8035277
	github.com/urfave/cli v1.20.1-0.20181029213200-b67dcf995b6a
	golang.org/x/net v0.0.0-20191004110552-13f9640d40b9
	golang.org/x/oauth2 v0.0.0-20180821212333-d2e6202438be
	google.golang.org/grpc v1.19.1
	gopkg.in/yaml.v2 v2.2.4
	k8s.io/apimachinery v0.0.0-20191213205550-43dfa306b0f3
	k8s.io/client-go v11.0.0+incompatible
)

replace (
	github.com/openzipkin-contrib/zipkin-go-opentracing v0.3.5 => github.com/go-chassis/zipkin-go-opentracing v0.3.5-0.20190321072447-42cf74fc2a92
	github.com/tetratelabs/go2sky v0.1.1-0.20190703154722-1eaab8035277 => github.com/SkyAPM/go2sky v0.1.1-0.20190703154722-1eaab8035277
)

go 1.13

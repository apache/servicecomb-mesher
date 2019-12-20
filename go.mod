module github.com/apache/servicecomb-mesher

require (
	github.com/apache/servicecomb-kie v0.1.1-0.20191119112752-d564a5ac0115
	github.com/envoyproxy/go-control-plane v0.6.0
	github.com/ghodss/yaml v1.0.0
	github.com/go-chassis/foundation v0.1.1-0.20191113114104-2b05871e9ec4
	github.com/go-chassis/go-archaius v1.0.0
	github.com/go-chassis/go-chassis v1.8.1-0.20191220080831-5b58b80f31bc
	github.com/go-chassis/gohessian v0.0.0-20180702061429-e5130c25af55
	github.com/go-mesh/openlogging v1.0.1
	github.com/gogo/googleapis v1.3.1 // indirect
	github.com/gogo/protobuf v1.3.0
	github.com/lyft/protoc-gen-validate v0.1.0 // indirect
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/prometheus/client_golang v0.9.1
	github.com/prometheus/client_model v0.0.0-20190115171406-56726106282f
	github.com/stretchr/testify v1.4.0
	github.com/tetratelabs/go2sky v0.1.1-0.20190703154722-1eaab8035277
	github.com/urfave/cli v1.20.1-0.20181029213200-b67dcf995b6a
	golang.org/x/net v0.0.0-20191004110552-13f9640d40b9
	golang.org/x/oauth2 v0.0.0-20190604053449-0f29369cfe45
	google.golang.org/grpc v1.19.1
	gopkg.in/yaml.v2 v2.2.4
	k8s.io/apimachinery v0.17.0
	k8s.io/client-go v0.17.0
)

replace (
	github.com/openzipkin-contrib/zipkin-go-opentracing v0.3.5 => github.com/go-chassis/zipkin-go-opentracing v0.3.5-0.20190321072447-42cf74fc2a92
	github.com/tetratelabs/go2sky v0.1.1-0.20190703154722-1eaab8035277 => github.com/SkyAPM/go2sky v0.1.1-0.20190703154722-1eaab8035277
)

go 1.13

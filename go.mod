module github.com/apache/servicecomb-mesher

require (
	github.com/envoyproxy/go-control-plane v0.9.1-0.20191026205805-5f8ba28d4473
	github.com/ghodss/yaml v1.0.0
	github.com/go-chassis/foundation v0.1.1-0.20200825060850-b16bf420f7b3
	github.com/go-chassis/go-archaius v1.3.3
	github.com/go-chassis/go-chassis/v2 v2.0.3-0.20200916043058-7a753c9f1471
	github.com/go-chassis/gohessian v0.0.0-20180702061429-e5130c25af55
	github.com/go-chassis/openlog v1.1.2
	github.com/go-mesh/openlogging v1.0.1
	github.com/gogo/googleapis v1.3.1 // indirect
	github.com/gogo/protobuf v1.3.0
	github.com/lyft/protoc-gen-validate v0.1.0 // indirect
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/prometheus/client_golang v0.9.1
	github.com/prometheus/client_model v0.0.0-20190812154241-14fe0d1b01d4
	github.com/stretchr/testify v1.5.1
	github.com/tetratelabs/go2sky v0.1.1-0.20190703154722-1eaab8035277
	github.com/urfave/cli v1.20.1-0.20181029213200-b67dcf995b6a
	golang.org/x/net v0.0.0-20200520004742-59133d7f0dd7
	golang.org/x/oauth2 v0.0.0-20190604053449-0f29369cfe45
	google.golang.org/grpc v1.27.0
	gopkg.in/yaml.v2 v2.3.0
	k8s.io/apimachinery v0.17.0
	k8s.io/client-go v0.17.0
)

replace (
	github.com/openzipkin-contrib/zipkin-go-opentracing v0.3.5 => github.com/go-chassis/zipkin-go-opentracing v0.3.5-0.20190321072447-42cf74fc2a92
	github.com/tetratelabs/go2sky v0.1.1-0.20190703154722-1eaab8035277 => github.com/SkyAPM/go2sky v0.1.1-0.20190703154722-1eaab8035277
)

go 1.13

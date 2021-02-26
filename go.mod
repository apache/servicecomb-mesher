module github.com/apache/servicecomb-mesher

require (
	github.com/ghodss/yaml v1.0.0
	github.com/go-chassis/foundation v0.2.2
	github.com/go-chassis/go-archaius v1.3.6-0.20201210061741-7450779aaeb8
	github.com/go-chassis/go-chassis/v2 v2.1.1
	github.com/go-chassis/gohessian v0.0.0-20180702061429-e5130c25af55
	github.com/go-chassis/openlog v1.1.2
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/prometheus/client_golang v0.9.1
	github.com/prometheus/client_model v0.0.0-20190812154241-14fe0d1b01d4
	github.com/stretchr/testify v1.6.1
	github.com/tetratelabs/go2sky v0.1.1-0.20190703154722-1eaab8035277
	github.com/urfave/cli v1.20.1-0.20181029213200-b67dcf995b6a
	golang.org/x/net v0.0.0-20201209123823-ac852fbbde11
	golang.org/x/oauth2 v0.0.0-20190604053449-0f29369cfe45
	google.golang.org/grpc v1.27.0
	gopkg.in/yaml.v2 v2.3.0
)

replace (
	github.com/openzipkin-contrib/zipkin-go-opentracing v0.3.5 => github.com/go-chassis/zipkin-go-opentracing v0.3.5-0.20190321072447-42cf74fc2a92
	github.com/tetratelabs/go2sky v0.1.1-0.20190703154722-1eaab8035277 => github.com/SkyAPM/go2sky v0.1.1-0.20190703154722-1eaab8035277
)

go 1.13

module github.com/apache/servicecomb-mesher

require (
	github.com/apache/servicecomb-kie v0.0.0-20191021014716-b4dc01d3968b
	github.com/envoyproxy/go-control-plane v0.6.0
	github.com/ghodss/yaml v1.0.0
	github.com/go-chassis/foundation v0.0.0-20190621030543-c3b63f787f4c
	github.com/go-chassis/go-archaius v0.24.0
	github.com/go-chassis/go-chassis v1.7.3-0.20191018125535-1a99ab41f7ea
	github.com/go-chassis/go-chassis-config v0.14.0
	github.com/go-chassis/gohessian v0.0.0-20180702061429-e5130c25af55
	github.com/go-mesh/openlogging v1.0.1
	github.com/gogo/googleapis v1.3.0 // indirect
	github.com/gogo/protobuf v1.3.0
	github.com/google/gofuzz v1.0.0 // indirect
	github.com/gopherjs/gopherjs v0.0.0-20181103185306-d547d1d9531e // indirect
	github.com/imdario/mergo v0.3.7 // indirect
	github.com/lyft/protoc-gen-validate v0.1.0 // indirect
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/prometheus/client_golang v0.9.1
	github.com/prometheus/client_model v0.0.0-20190115171406-56726106282f
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/stretchr/testify v1.3.0
	github.com/urfave/cli v1.20.1-0.20181029213200-b67dcf995b6a
	golang.org/x/net v0.0.0-20190311183353-d8887717615a
	golang.org/x/time v0.0.0-20180412165947-fbb02b2291d2 // indirect
	google.golang.org/grpc v1.16.0
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/yaml.v2 v2.2.1
	k8s.io/apimachinery v0.0.0-20181022183627-f71dbbc36e12
	k8s.io/client-go v9.0.0+incompatible
)

replace github.com/openzipkin-contrib/zipkin-go-opentracing v0.3.5 => github.com/go-chassis/zipkin-go-opentracing v0.3.5-0.20190321072447-42cf74fc2a92

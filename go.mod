module github.com/go-mesh/mesher

require (
	github.com/Shopify/toxiproxy v2.1.3+incompatible // indirect
	github.com/beorn7/perks v0.0.0-20180321164747-3a771d992973 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/eapache/go-resiliency v1.1.0 // indirect
	github.com/eapache/go-xerial-snappy v0.0.0-20180814174437-776d5712da21 // indirect
	github.com/eapache/queue v1.1.0 // indirect
	github.com/emicklei/go-restful-swagger12 v0.0.0-20170926063155-7524189396c6 // indirect
	github.com/envoyproxy/go-control-plane v0.6.0
	github.com/ghodss/yaml v1.0.0 // indirect

	github.com/go-chassis/go-cc-client v0.0.0-20180831085349-c2bb6cef1640
	github.com/go-chassis/go-chassis v0.8.4-0.20180928015049-b4c551ac46e1
	github.com/go-chassis/gohessian v0.0.0-20180702061429-e5130c25af55
	github.com/go-logfmt/logfmt v0.3.0 // indirect
	github.com/go-mesh/openlogging v0.0.0-20180912071658-0fd4707a75ab
	github.com/gogo/googleapis v1.1.0 // indirect
	github.com/gogo/protobuf v1.1.1
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b // indirect
	github.com/golang/snappy v0.0.0-20180518054509-2e65f85255db // indirect
	github.com/google/gofuzz v0.0.0-20170612174753-24818f796faf // indirect
	github.com/imdario/mergo v0.3.6 // indirect
	github.com/jtolds/gls v4.2.1+incompatible // indirect
	github.com/kr/logfmt v0.0.0-20140226030751-b84e30acd515 // indirect
	github.com/lyft/protoc-gen-validate v0.0.10 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.1 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/onsi/gomega v1.4.2 // indirect
	github.com/opentracing-contrib/go-observer v0.0.0-20170622124052-a52f23424492 // indirect
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/prometheus/client_golang v0.8.0
	github.com/prometheus/client_model v0.0.0-20180712105110-5c3871d89910
	github.com/prometheus/common v0.0.0-20181020173914-7e9e6cabbd39 // indirect
	github.com/spf13/pflag v1.0.3 // indirect
	github.com/stretchr/objx v0.1.1 // indirect
	github.com/stretchr/testify v1.2.2
	github.com/uber-go/atomic v1.3.2 // indirect
	github.com/urfave/cli v0.0.0-20180821064027-934abfb2f102
	go.uber.org/atomic v1.3.2 // indirect
	golang.org/x/crypto v0.0.0-20181030102418-4d3f4d9ffa16 // indirect
	golang.org/x/net v0.0.0-20180906233101-161cd47e91fd
	golang.org/x/sys v0.0.0-20181031143558-9b800f95dbbc // indirect
	golang.org/x/time v0.0.0-20180412165947-fbb02b2291d2 // indirect
	google.golang.org/genproto v0.0.0-20181101192439-c830210a61df // indirect
	google.golang.org/grpc v1.14.0
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/yaml.v2 v2.2.1
	k8s.io/apimachinery v0.0.0-20180823151430-fda675fbe85280c4550452dae2a5ebf74e4a59b7
	k8s.io/client-go v8.0.0+incompatible
)

replace (
	cloud.google.com/go v0.28.0 => github.com/GoogleCloudPlatform/google-cloud-go v0.28.0
	github.com/envoyproxy/go-control-plane v0.6.0 => github.com/envoyproxy/go-control-plane v0.0.0-20180918192855-2137d919632883e52e7786f55f0f84e52a44fbf3
	github.com/kubernetes/client-go => ../k8s.io/client-go
	golang.org/x/crypto v0.0.0-20180820150726-614d502a4dac => github.com/golang/crypto v0.0.0-20180820150726-614d502a4dac
	golang.org/x/net v0.0.0-20180724234803-3673e40ba225 => github.com/golang/net v0.0.0-20180724234803-3673e40ba225
	golang.org/x/net v0.0.0-20180824152047-4bcd98cce591 => github.com/golang/net v0.0.0-20180824152047-4bcd98cce591

	golang.org/x/oauth2 v0.0.0-20180821212333-d2e6202438be => github.com/golang/oauth2 v0.0.0-20180821212333-d2e6202438be
	golang.org/x/sys v0.0.0-20180824143301-4910a1d54f87 => github.com/golang/sys v0.0.0-20180824143301-4910a1d54f87
	golang.org/x/text v0.3.0 => github.com/golang/text v0.3.0
	golang.org/x/time v0.0.0-20180412165947-fbb02b2291d2 => github.com/golang/time v0.0.0-20180412165947-fbb02b2291d2
	google.golang.org/appengine v1.2.0 => github.com/golang/appengine v1.2.0
	google.golang.org/genproto v0.0.0-20180817151627-c66870c02cf8 => github.com/google/go-genproto v0.0.0-20180817151627-c66870c02cf8
	google.golang.org/grpc v1.14.0 => github.com/grpc/grpc-go v1.14.0
)

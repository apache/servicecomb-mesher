module github.com/go-mesh/mesher

require (
	github.com/Shopify/toxiproxy v2.1.3+incompatible // indirect
	github.com/emicklei/go-restful-swagger12 v0.0.0-20170926063155-7524189396c6 // indirect
	github.com/envoyproxy/go-control-plane v0.6.0
	github.com/ghodss/yaml v1.0.0 // indirect
	github.com/go-chassis/go-archaius v0.0.0-20181108111652-ab19b4eae276

	github.com/go-chassis/go-cc-client v0.0.0-20181102101915-dea430061a34
	github.com/go-chassis/go-chassis v1.1.2-0.20181210100135-5e647f1a038b
	github.com/go-chassis/gohessian v0.0.0-20180702061429-e5130c25af55
	github.com/go-mesh/openlogging v0.0.0-20181122085847-3daf3ad8ed35
	github.com/gogo/googleapis v1.1.0 // indirect
	github.com/gogo/protobuf v1.1.1
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b // indirect
	github.com/google/gofuzz v0.0.0-20170612174753-24818f796faf // indirect
	github.com/imdario/mergo v0.3.6 // indirect
	github.com/lyft/protoc-gen-validate v0.0.11 // indirect
	github.com/onsi/gomega v1.4.2 // indirect
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/prometheus/client_golang v0.8.0
	github.com/prometheus/client_model v0.0.0-20180712105110-5c3871d89910
	github.com/stretchr/testify v1.2.2
	github.com/uber-go/atomic v1.3.2 // indirect
	github.com/urfave/cli v1.20.1-0.20181029213200-b67dcf995b6a
	go.uber.org/atomic v1.3.2 // indirect
	golang.org/x/crypto v0.0.0-20181030102418-4d3f4d9ffa16 // indirect
	golang.org/x/net v0.0.0-20180906233101-161cd47e91fd
	golang.org/x/time v0.0.0-20180412165947-fbb02b2291d2 // indirect
	google.golang.org/genproto v0.0.0-20181101192439-c830210a61df // indirect
	google.golang.org/grpc v1.14.0
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/yaml.v2 v2.2.1
	k8s.io/apimachinery v0.0.0-20181022183627-f71dbbc36e12
	k8s.io/client-go v9.0.0+incompatible
)

replace (
	cloud.google.com/go v0.28.0 => github.com/GoogleCloudPlatform/google-cloud-go v0.28.0
	github.com/envoyproxy/go-control-plane v0.6.0 => github.com/envoyproxy/go-control-plane v0.0.0-20180918192855-2137d919632883e52e7786f55f0f84e52a44fbf3
	github.com/kubernetes/client-go => ../k8s.io/client-go
	golang.org/x/crypto v0.0.0-20181030102418-4d3f4d9ffa16 => github.com/golang/crypto v0.0.0-20181030102418-4d3f4d9ffa16
	golang.org/x/net v0.0.0-20180906233101-161cd47e91fd => github.com/golang/net v0.0.0-20180906233101-161cd47e91fd
	golang.org/x/oauth2 v0.0.0-20180207181906-543e37812f10 => github.com/golang/oauth2 v0.0.0-20180207181906-543e37812f10
	golang.org/x/sync v0.0.0-20180314180146-1d60e4601c6f => github.com/golang/sync v0.0.0-20180314180146-1d60e4601c6f
	golang.org/x/sys v0.0.0-20180909124046-d0be0721c37e => github.com/golang/sys v0.0.0-20180909124046-d0be0721c37e
	golang.org/x/sys v0.0.0-20181031143558-9b800f95dbbc => github.com/golang/sys v0.0.0-20181031143558-9b800f95dbbc

	golang.org/x/text v0.3.0 => github.com/golang/text v0.3.0
	golang.org/x/time v0.0.0-20180412165947-fbb02b2291d2 => github.com/golang/time v0.0.0-20180412165947-fbb02b2291d2
	google.golang.org/appengine v1.2.0 => github.com/golang/appengine v1.2.0
	google.golang.org/genproto v0.0.0-20181101192439-c830210a61df => github.com/google/go-genproto v0.0.0-20181101192439-c830210a61df
	google.golang.org/grpc v1.14.0 => github.com/grpc/grpc-go v1.14.0
)

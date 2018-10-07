module github.com/go-mesh/mesher

require (
	cloud.google.com/go v0.28.0 // indirect
	code.cloudfoundry.org/copilot v0.0.0-20180928002835-76734bdb9045 // indirect
	fortio.org/fortio v1.3.0 // indirect
	github.com/emicklei/go-restful-swagger12 v0.0.0-20170926063155-7524189396c6 // indirect
	github.com/envoyproxy/go-control-plane v0.6.0

	github.com/go-chassis/go-cc-client v0.0.0-20180831085349-c2bb6cef1640
	github.com/go-chassis/go-chassis v0.8.3-0.20180914033538-0791a5cec8b4
	github.com/go-chassis/gohessian v0.0.0-20180702061429-e5130c25af55
	github.com/go-mesh/mesher-tools v0.0.0-20181006103649-cdc091b78a72
	github.com/go-mesh/openlogging v0.0.0-20180831021158-f5d1c4e7e506
	github.com/gogo/protobuf v1.1.1
	github.com/gogo/status v1.0.3 // indirect
	github.com/golang/groupcache v0.0.0-20180924190550-6f2cf27854a4 // indirect
	github.com/golang/sync v0.0.0-20180314180146-1d60e4601c6f // indirect
	github.com/google/go-github v17.0.0+incompatible // indirect
	github.com/google/go-querystring v1.0.0 // indirect
	github.com/google/uuid v1.0.0 // indirect
	github.com/gorilla/context v1.1.1 // indirect
	github.com/gorilla/mux v1.6.2 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.0.0 // indirect
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0 // indirect
	github.com/hashicorp/consul v1.2.3 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.0 // indirect
	github.com/hashicorp/go-multierror v1.0.0 // indirect
	github.com/hashicorp/go-rootcerts v0.0.0-20160503143440-6bb64b370b90 // indirect
	github.com/hashicorp/serf v0.8.1 // indirect
	github.com/howeyc/fsnotify v0.9.0 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/mitchellh/go-homedir v1.0.0 // indirect
	github.com/mitchellh/mapstructure v1.0.0 // indirect
	github.com/natefinch/lumberjack v2.0.0+incompatible // indirect
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/pkg/errors v0.8.0 // indirect
	github.com/prometheus/client_golang v0.8.0
	github.com/prometheus/client_model v0.0.0-20180712105110-5c3871d89910
	github.com/prometheus/prom2json v0.0.0-20180620215746-7b8ed2aed129 // indirect
	github.com/spf13/cobra v0.0.3 // indirect
	github.com/stretchr/testify v1.2.2
	github.com/urfave/cli v0.0.0-20180821064027-934abfb2f102
	go.uber.org/atomic v1.3.2 // indirect
	go.uber.org/multierr v1.1.0 // indirect
	go.uber.org/zap v1.9.1 // indirect
	golang.org/x/oauth2 v0.0.0-20180821212333-d2e6202438be // indirect
	google.golang.org/appengine v1.2.0 // indirect
	google.golang.org/grpc v1.14.0
	gopkg.in/validator.v2 v2.0.0-20180514200540-135c24b11c19 // indirect
	gopkg.in/yaml.v2 v2.2.1
	istio.io/api v0.0.0-20180926203357-0cf306e2fd19 // indirect
	istio.io/fortio v1.3.0 // indirect
	istio.io/istio v0.0.0-20180929031539-a8986e2dc2c3
	k8s.io/apiextensions-apiserver v0.0.0-20180925155151-ce69c54e57693220512104c84941e2ef1876449a // indirect
	k8s.io/apimachinery v0.0.0-20180823151430-fda675fbe85280c4550452dae2a5ebf74e4a59b7
	k8s.io/client-go v8.0.0+incompatible
	k8s.io/cluster-registry v0.0.6 // indirect

	k8s.io/ingress v0.0.0-20170803151325-fe19ebb09ee2 // indirect
	k8s.io/kube-openapi v0.0.0-20180928070517-c01ed926f124 // indirect
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

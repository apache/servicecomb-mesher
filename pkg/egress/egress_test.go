package egress_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/go-chassis/go-chassis/core/config"
	"github.com/go-chassis/go-chassis/core/lager"
	"github.com/go-mesh/mesher/cmd"
	mesherconfig "github.com/go-mesh/mesher/config"
	egressmodel "github.com/go-mesh/mesher/config/model"
	"github.com/go-mesh/mesher/control"
	"github.com/go-mesh/mesher/control/archiaus"
	_ "github.com/go-mesh/mesher/control/archiaus"
	_ "github.com/go-mesh/mesher/control/istio"
	"github.com/go-mesh/mesher/pkg/egress"
	"gopkg.in/yaml.v2"
)

func BenchmarkMatch(b *testing.B) {
	lager.Initialize("", "DEBUG", "",
		"size", true, 1, 10, 7)

	gopath := os.Getenv("GOPATH")
	os.Setenv("CHASSIS_HOME", gopath+"/src/github.com/go-mesh/mesher")
	cmd.Init()
	config.Init()
	mesherconfig.Init()
	control.Init()
	var yamlContent = `---
egress:
  infra: cse # pilot or cse
  address: http://istio-pilot.istio-system:15010
egressRule:
  google-ext:
    - hosts:
        - "www.google.com"
        - "*.yahoo.com"
      ports:
        - port: 80
          protocol: HTTP
  facebook-ext:
    - hosts:
        - "www.facebook.com"
      ports:
        - port: 80
          protocol: HTTP`

	ss := egressmodel.EgressConfig{}
	err := yaml.Unmarshal([]byte(yamlContent), &ss)
	if err != nil {
		fmt.Println("unmarshal failed")
	}
	archiaus.SaveToEgressCache(&ss)
	fmt.Println(archiaus.EgressConfigCache.Get(""))

	myString := "www.google.com"
	for i := 0; i < b.N; i++ {
		egress.Match(myString)
	}
}

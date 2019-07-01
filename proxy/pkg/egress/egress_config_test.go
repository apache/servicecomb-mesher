package egress_test

import (
	"fmt"
	mesherconfig "github.com/go-mesh/mesher/proxy/config"
	"github.com/go-mesh/mesher/proxy/pkg/egress"
	"gopkg.in/yaml.v2"
	"testing"
)

func TestValidateEgressRule(t *testing.T) {
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

	ss := mesherconfig.EgressConfig{}
	err := yaml.Unmarshal([]byte(yamlContent), &ss)
	if err != nil {
		fmt.Println("unmarshal failed")
	}

	bool, err := egress.ValidateEgressRule(ss.Destinations)
	if bool == false {
		t.Errorf("Expected true but got false")
	}
}

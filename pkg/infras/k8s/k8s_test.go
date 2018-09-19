package pilotv2

import (
	"os"
	"os/user"
	"testing"
)

var KubeConfig string

func init() {
	if KUBE_CONFIG := os.Getenv("KUBE_CONFIG"); KUBE_CONFIG != "" {
		KubeConfig = KUBE_CONFIG
	} else {
		usr, err := user.Current()
		if err != nil {
			panic("Failed to get current user info: " + err.Error())
		} else {
			KubeConfig = usr.HomeDir + "/" + ".kube/config"
		}
	}

}

func TestCreateK8sClient(t *testing.T) {
	_, err := CreateK8SRestClient(KubeConfig, "apis", "networking.istio.io", "v1alpha3")
	if err != nil {
		t.Errorf("Failed to create k8s rest client: %s", err.Error())
	}

	_, err = CreateK8SRestClient("*nonfile", "apis", "networking.istio.io", "v1alpha3")
	if err == nil {
		t.Errorf("Test failed, should return error with invalid kube config path")
	}
}

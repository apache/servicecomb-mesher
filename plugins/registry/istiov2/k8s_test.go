package pilotv2

import "testing"

var validkubeconfig string = "/home/lance/.kube/config"

func TestCreateK8sClient(t *testing.T) {
	_, err := CreateK8SRestClient(validkubeconfig, "apis", "networking.istio.io", "v1alpha3")
	if err != nil {
		t.Errorf("Failed to create k8s rest client: %s", err.Error())
	}

	_, err = CreateK8SRestClient("*nonfile", "apis", "networking.istio.io", "v1alpha3")
	if err == nil {
		t.Errorf("Test failed, should return error with invalid kube config path")
	}

}

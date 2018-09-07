package pilotv2

import (
	"testing"
)

var (
	kubeconfig     string = "/home/lance/.kube/config"
	ValidPilotAddr string = "192.168.0.70:13602"
	nodeInfo       *NodeInfo

	testXdsClient    *XdsClient
	testCacheManager *CacheManager
	err              error
)

func init() {
	nodeInfo = &NodeInfo{
		PodName:    "test-pod",
		Namespace:  "default",
		InstanceIP: "192.168.43.1",
	}

	testXdsClient, err = NewXdsClient(ValidPilotAddr, nil, nodeInfo, kubeconfig)
	if err != nil {
		panic("Failed to prepare test, xds client creation failed: " + err.Error())
	}
}

func TestNewCacheManager(t *testing.T) {
	testCacheManager, err = NewCacheManager(testXdsClient)
	if err != nil {
		t.Errorf("Failed to create CacheManager: %s", err.Error())
	}
}

// func TestAutoSync(t *testing.T) {
//     testCacheManager.AutoSync()
// }

func TestPullImcroserviceInstance(t *testing.T) {
	err = testCacheManager.pullMicroserviceInstance()
	if err != nil {
		t.Errorf("Failed to pull microservice instances: %s", err.Error())
	}
}

// func TestMakeIPIndex(t *testing.T) {
//     err := testCacheManager.MakeIPIndex()
//     if err != nil {
//         t.Errorf("Failed to make ip index: %s", err.Error())
//     }
// }

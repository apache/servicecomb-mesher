package pilotv2

import (
	"fmt"
	"os"
	"os/user"
	"testing"

	"github.com/go-chassis/go-chassis/pkg/util/iputil"
)

const (
	TEST_POD_NAME     = "testpod"
	NAMESPACE_DEFAULT = "default"
)

var (
	KubeConfig     string
	ValidPilotAddr string
	LocalIPAddress string
	nodeInfo       *NodeInfo

	testXdsClient    *XdsClient
	testCacheManager *CacheManager
	err              error
)

func init() {
	// Get kube config path and local ip
	if KUBE_CONFIG := os.Getenv("KUBE_CONFIG"); KUBE_CONFIG != "" {
		KubeConfig = KUBE_CONFIG
	} else {
		usr, err := user.Current()
		if err != nil {
			panic(fmt.Sprintf("Failed to get current user info: %s", err.Error()))
		} else {
			KubeConfig = fmt.Sprintf("%s/%s", usr.HomeDir, ".kube/config")
		}
	}

	if PILOT_ADDR := os.Getenv("PILOT_ADDR"); PILOT_ADDR != "" {
		ValidPilotAddr = PILOT_ADDR
	} else {
		panic("PILOT_ADDR should be specified to pass the pilot address")
	}

	if INSTANCE_IP := os.Getenv("INSTANCE_IP"); INSTANCE_IP != "" {
		LocalIPAddress = INSTANCE_IP
	} else if LocalIPAddress = iputil.GetLocalIP(); LocalIPAddress == "" {
		panic("Failed to get the local ip address, please check the network environment")
	}

	nodeInfo = &NodeInfo{
		PodName:    TEST_POD_NAME,
		Namespace:  NAMESPACE_DEFAULT,
		InstanceIP: LocalIPAddress,
	}

	testXdsClient, err = NewXdsClient(ValidPilotAddr, nil, nodeInfo, KubeConfig)
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
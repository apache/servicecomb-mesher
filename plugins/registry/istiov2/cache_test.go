package istiov2

import (
	"os"
	"os/user"
	"testing"

	"github.com/go-chassis/go-chassis/core/lager"
	"github.com/go-chassis/go-chassis/core/registry"
	"github.com/go-chassis/go-chassis/pkg/util/iputil"
	testutil "github.com/go-mesh/mesher-tools/test/util"
	istioinfra "github.com/go-mesh/mesher/pkg/infras/istio"
	"istio.io/istio/tests/util"
)

const (
	TEST_POD_NAME     = "testpod"
	NAMESPACE_DEFAULT = "default"
)

var (
	KubeConfig     string
	ValidPilotAddr string
	LocalIPAddress string
	nodeInfo       *istioinfra.NodeInfo

	testXdsClient    *istioinfra.XdsClient
	testCacheManager *CacheManager
	err              error
)

func TestMain(t *testing.T) {
	lager.Initialize("", "DEBUG", "", "size", true, 1, 10, 7)
	// Get kube config path and local ip
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

	if PILOT_ADDR := os.Getenv("PILOT_ADDR"); PILOT_ADDR != "" {
		ValidPilotAddr = PILOT_ADDR
	} else {
		// panic("PILOT_ADDR should be specified to pass the pilot address")
		testutil.InitLocalPilotTestEnv(t)
		ValidPilotAddr = util.MockPilotGrpcAddr
	}

	if INSTANCE_IP := os.Getenv("INSTANCE_IP"); INSTANCE_IP != "" {
		LocalIPAddress = INSTANCE_IP
	} else if LocalIPAddress = iputil.GetLocalIP(); LocalIPAddress == "" {
		panic("Failed to get the local ip address, please check the network environment")
	}

	nodeInfo = &istioinfra.NodeInfo{
		PodName:    TEST_POD_NAME,
		Namespace:  NAMESPACE_DEFAULT,
		InstanceIP: LocalIPAddress,
	}

	testXdsClient, err = istioinfra.NewXdsClient(ValidPilotAddr, nil, nodeInfo, KubeConfig)
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

func TestEndpointCache(t *testing.T) {
	ec := EndpointCache{
		cache: map[string]EndpointSubset{},
	}

	subset := EndpointSubset{
		subsetName: "foo",
		tags:       map[string]string{},
		instances:  []*registry.MicroServiceInstance{},
	}

	defer func() {
		if r := recover(); r != nil {
			t.Error("should not panic")
		}
	}()

	waitChannel := make(chan int)
	for i := 0; i < 1000; i++ {
		go func() {
			ec.Set("foo", subset)
			waitChannel <- 0

		}()
	}

	for i := 0; i < 1000; i++ {
		<-waitChannel
	}
}

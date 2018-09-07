package pilotv2

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type DestinationRuleResult struct {
	Items []MinDestinationRule
}

// MinDestinationRule is the minimum structure we need to get subsets
type MinDestinationRule struct {
	Metadata struct {
		Name      string `json:"name"`
		Namespace string `json:"namespace"`
	} `json:"metadata"`
	Spec struct {
		Host    string `json:"host"`
		Subsets []struct {
			Labels map[string]string `json:"labels"`
			Name   string            `json:"name"`
		} `json:"subsets"`
	} `json:"spec"`
}

func CreateK8SRestClient(kubeconfig, apiPath, group, version string) (*rest.RESTClient, error) {
	var config *rest.Config
	var err error
	if kubeconfig != "" {
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			config, err = rest.InClusterConfig()
		}
	} else {
		config, err = rest.InClusterConfig()
	}

	if err != nil {
		return nil, err
	}

	config.APIPath = apiPath
	config.GroupVersion = &schema.GroupVersion{
		Group:   group,
		Version: version,
	}
	config.NegotiatedSerializer = serializer.DirectCodecFactory{CodecFactory: serializer.NewCodecFactory(runtime.NewScheme())}

	k8sRestClient, err := rest.RESTClientFor(config)
	if err != nil {
		return nil, err
	}
	return k8sRestClient, nil
}

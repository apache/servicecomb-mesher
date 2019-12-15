/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package pilotv2

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

//DestinationRuleResult is the list of MinDestinationRules
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

//CreateK8SRestClient returns the kubernetes client for RESTful API calls
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
	config.NegotiatedSerializer = serializer.NewCodecFactory(runtime.NewScheme())

	k8sRestClient, err := rest.RESTClientFor(config)
	if err != nil {
		return nil, err
	}
	return k8sRestClient, nil
}

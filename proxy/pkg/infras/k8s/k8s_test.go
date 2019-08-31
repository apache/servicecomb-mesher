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

func TestCreateInvalidK8sClient(t *testing.T) {
	_, err := CreateK8SRestClient("", "apis", "networking.istio.io", "v1alpha3")
	if err == nil {
		t.Errorf("Passing a nil config for k8s client should return error")
	}
}

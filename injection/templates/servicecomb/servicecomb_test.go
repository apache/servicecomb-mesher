/*
 *  Licensed to the Apache Software Foundation (ASF) under one or more
 *  contributor license agreements.  See the NOTICE file distributed with
 *  this work for additional information regarding copyright ownership.
 *  The ASF licenses this file to You under the Apache License, Version 2.0
 *  (the "License"); you may not use this file except in compliance with
 *  the License.  You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 */

package servicecomb

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
)

var (
	testConfigFile   = "injectorconfig.yaml"
	testTemplateFile = "injectortempalet.yaml"

	configContent = []byte(`serviceCenter:
  tlsConfig:
    certFile: aaa`)

	templateContent = []byte(`apiVersion: v1
kind: Pod
spec:
  containers:
  - env:
    - name: http_proxy
      value: http://127.0.0.1:{{.AMesher.HTTPPort}}
    name: {{.App.AName}}`)
)

func init() {
	os.Remove(testConfigFile)
	os.Remove(testTemplateFile)
}

func TestNewServiceComb(t *testing.T) {
	t.Log("========New ServiceComb Template no configPath and no templatePath=")
	tmpl, err := NewServiceComb("", "")
	assert.Nil(t, err)
	assert.NotNil(t, tmpl)

	t.Log("========New ServiceComb Template wrong configPath and no templatePath=")
	tmpl, err = NewServiceComb(testConfigFile, "")
	assert.NotNil(t, err)

	t.Log("========New ServiceComb Template no configPath and wrong templatePath=")
	tmpl, err = NewServiceComb("", testTemplateFile)
	assert.NotNil(t, err)
}

func TestServiceComb_UpdateConfig(t *testing.T) {
	t.Log("========ServiceComb Template UpdateConfig with wrong configPath=")
	tmpl, err := NewServiceComb("", "")
	assert.Nil(t, err)

	err = tmpl.UpdateConfig(testConfigFile)
	assert.NotNil(t, err)

	t.Log("========ServiceComb Template UpdateConfig with configPath=")
	err = ioutil.WriteFile(testConfigFile, configContent, 0640)
	assert.Nil(t, err)
	defer os.Remove(testConfigFile)

	err = tmpl.UpdateConfig(testConfigFile)
	assert.Nil(t, err)

}

func TestServiceComb_UpdateTemplate(t *testing.T) {
	t.Log("========ServiceComb Template UpdateTemplate with wrong templatePath=")
	tmpl, err := NewServiceComb("", "")
	assert.Nil(t, err)

	err = tmpl.UpdateTemplate(testTemplateFile)
	assert.NotNil(t, err)

	t.Log("========ServiceComb Template UpdateTemplate with templatePath=")
	err = ioutil.WriteFile(testTemplateFile, templateContent, 0640)
	assert.Nil(t, err)
	defer os.Remove(testTemplateFile)

	err = tmpl.UpdateTemplate(testTemplateFile)
	assert.Nil(t, err)
}

func TestServiceComb_PodSpecFromTemplate(t *testing.T) {
	tmpl, err := NewServiceComb("", "")
	assert.Nil(t, err)

	t.Log("========ServiceComb Template PodSpecFromTemplate with empty pod=")
	pod := &corev1.Pod{}
	_, err = tmpl.PodSpecFromTemplate(pod)
	assert.NotNil(t, err)

	t.Log("========ServiceComb Template PodSpecFromTemplate with pod=")
	pod.GenerateName = "testPod-1"
	pod.Spec.Containers = append(pod.Spec.Containers, corev1.Container{
		Name:  "testPod",
		Ports: []corev1.ContainerPort{},
		Env: []corev1.EnvVar{
			{
				Name:  "VERSION",
				Value: "1.0.1",
			},
		},
	})
	_, err = tmpl.PodSpecFromTemplate(pod)
	assert.Nil(t, err)

	pod.Spec.Containers[0].Ports = append(pod.Spec.Containers[0].Ports, corev1.ContainerPort{
		Protocol:      corev1.ProtocolTCP,
		ContainerPort: 8080,
		Name:          "http",
	})
	_, err = tmpl.PodSpecFromTemplate(pod)
	assert.Nil(t, err)

	t.Log("========ServiceComb Template PodSpecFromTemplate with wrong template=")
	err = ioutil.WriteFile(testTemplateFile, templateContent, 0640)
	assert.Nil(t, err)
	defer os.Remove(testTemplateFile)

	err = tmpl.UpdateTemplate(testTemplateFile)
	assert.Nil(t, err)

	_, err = tmpl.PodSpecFromTemplate(pod)
	assert.NotNil(t, err)
}

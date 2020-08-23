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

package webhook

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/apache/servicecomb-mesher/injection/templates"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
	corev1 "k8s.io/api/core/v1"
)

var (
	mockTemplateName = "test-mock"
)

type mockTemplate struct{}

func init() {
	templates.Register(mockTemplateName, NewMockTemplate)
}

func NewMockTemplate(configPath, templatePath string) (templates.Templater, error) {
	return &mockTemplate{}, nil
}

func (t *mockTemplate) PodSpecFromTemplate(*corev1.Pod) (*corev1.PodSpec, error) {
	pod := &corev1.Pod{}
	err := yaml.Unmarshal(getTmplPodBytes(), pod)
	if err != nil {
		return nil, err
	}
	return &(pod.Spec), nil
}

func (t *mockTemplate) UpdateConfig(configPath string) error {
	return nil
}

func (t *mockTemplate) UpdateTemplate(tmplConfig string) error {
	return nil
}

func TestNewWebhook(t *testing.T) {
	t.Log("========Webhook Mock Template=")
	_, err := NewWebhook(WithTemplateName(mockTemplateName))
	assert.Nil(t, err)

	t.Log("========Webhook Template Notfound=")
	_, err = NewWebhook(WithTemplateName("test-notfound"))
	assert.NotNil(t, err)
}

func TestInjectHandler(t *testing.T) {
	t.Log("========Webhook InjectHandler Method Not Allowed=")
	wh, err := NewWebhook(WithTemplateName(mockTemplateName))
	assert.Nil(t, err)
	svr := httptest.NewServer(http.HandlerFunc(wh.injectHandler))
	defer svr.Close()

	resp, err := http.Get(svr.URL)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)

	t.Log("========Webhook InjectHandler Unsupported Media Type=")
	resp, err = http.Post(svr.URL, "", nil)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusUnsupportedMediaType, resp.StatusCode)

	t.Log("========Webhook InjectHandler Body is nil=")
	resp, err = http.Post(svr.URL, "application/json", nil)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	t.Log("========Webhook InjectHandler Body is wrong=")
	resp, err = http.Post(svr.URL, "application/json", bytes.NewBuffer([]byte("abc")))
	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	t.Log("========Webhook InjectHandler Success=")
	resp, err = http.Post(svr.URL, "application/json", bytes.NewBuffer(getMockBytes()))
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func getMockBytes() []byte {
	return []byte(`{
	"kind": "AdmissionReview",
	"apiVersion": "admission.k8s.io/v1beta1",
	"request": {
		"uid": "1caf247a-2076-11ea-bcbc-fa163eca30e0",
		"kind": {
			"group": "",
			"version": "v1",
			"kind": "Pod"
		},
		"resource": {
			"group": "",
			"version": "v1",
			"resource": "pods"
		},
		"namespace": "svccomb-test",
		"operation": "CREATE",
		"userInfo": {
			"username": "system:serviceaccount:kube-system:replicaset-controller",
			"uid": "57c75718-fa18-11e9-8271-fa163eca30e0",
			"groups": ["system:serviceaccounts", "system:serviceaccounts:kube-system", "system:authenticated"]
		},
		"object": {
			"metadata": {
				"generateName": "calculator-python-8dd449c6b-",
				"creationTimestamp": null,
				"labels": {
					"app": "calculator",
					"pod-template-hash": "8dd449c6b"
				},
				"ownerReferences": [{
					"apiVersion": "apps/v1",
					"kind": "ReplicaSet",
					"name": "calculator-python-8dd449c6b",
					"uid": "1cac4646-2076-11ea-bcbc-fa163eca30e0",
					"controller": true,
					"blockOwnerDeletion": true
				}]
			},
			"spec": {
				"volumes": [{
					"name": "default-token-9pvt7",
					"secret": {
						"secretName": "default-token-9pvt7"
					}
				}],
				"containers": [{
					"name": "calculator",
					"image": "servicecomb/calculator-python:latest",
					"ports": [{
						"containerPort": 4540,
						"protocol": "TCP"
					}],
					"resources": {},
					"volumeMounts": [{
						"name": "default-token-9pvt7",
						"readOnly": true,
						"mountPath": "/var/run/secrets/kubernetes.io/serviceaccount"
					}],
					"terminationMessagePath": "/dev/termination-log",
					"terminationMessagePolicy": "File",
					"imagePullPolicy": "IfNotPresent"
				}],
				"restartPolicy": "Always",
				"terminationGracePeriodSeconds": 30,
				"dnsPolicy": "ClusterFirst",
				"serviceAccountName": "default",
				"serviceAccount": "default",
				"securityContext": {},
				"schedulerName": "default-scheduler",
				"tolerations": [{
					"key": "node.kubernetes.io/not-ready",
					"operator": "Exists",
					"effect": "NoExecute",
					"tolerationSeconds": 300
				}, {
					"key": "node.kubernetes.io/unreachable",
					"operator": "Exists",
					"effect": "NoExecute",
					"tolerationSeconds": 300
				}],
				"priority": 0,
				"enableServiceLinks": true
			},
			"status": {}
		},
		"oldObject": null,
		"dryRun": false
	}
}`)
}

func getTmplPodBytes() []byte {
	return []byte(`apiVersion: v1
kind: Pod
spec:
  containers:
  - env:
    - name: http_proxy
      value: http://127.0.0.1:30101
    name: calculator
    resources: {}
    volumeMounts:
    - mountPath: /var/run/secrets/kubernetes.io/serviceaccount
      name: default-token-9pvt7
      readOnly: true
  - env:
    - name: SPECIFIC_ADDR
      value: 127.0.0.1:4540
    - name: SERVICE_NAME
      value: calculator
    - name: VERSION
      value: 1.0.0
    - name: CSE_REGISTRY_ADDR
      value: http://servicecenter.servicecomb.svc.cluster.local:30100
    image: servicecomb/mesher-sidecar:1.6.3
    name: svccomb-mesher
    ports:
    - containerPort: 40101
      name: grpc
      protocol: TCP
    - containerPort: 30101
      name: http
      protocol: TCP
    - containerPort: 30102
      name: rest-admin
      protocol: TCP
    resources: {}
    volumeMounts:
    - mountPath: /var/run/secrets/kubernetes.io/serviceaccount
      name: default-token-9pvt7
      readOnly: true`)
}

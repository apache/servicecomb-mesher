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
	"text/template"

	"github.com/go-mesh/openlogging"
)

var sidecarContainer = `apiVersion: v1
kind: Pod
spec:
  containers:
  - env:
    - name: http_proxy
      value: http://127.0.0.1:{{.Mesher.HTTPPort}}
    name: {{.App.Name}}
  - env:
    - name: SPECIFIC_ADDR
      value: 127.0.0.1:{{.App.Port}}
    - name: SERVICE_NAME
      value: {{.App.Name}}
    - name: VERSION
      value: {{.App.Version}}
    - name: CSE_REGISTRY_ADDR
      value: {{.ServiceCenter.Address}}
    image: {{.Mesher.Image}}:{{.Mesher.Tag}}
    imagePullPolicy: IfNotPresent
    name: {{.Mesher.Name}}
    ports:
    - containerport: {{.Mesher.GRPCPort}}
      name: grpc
      protocol: TCP
    - containerport: {{.Mesher.HTTPPort}}
      name: http
      protocol: TCP
    - containerport: {{.Mesher.AdminPort}}
      name: rest-admin
      protocol: TCP`

// DefaultTmpl returns default template
func DefaultTmpl() *template.Template {
	tmpl, err := template.New("sidecar").Parse(sidecarContainer)
	if err != nil {
		openlogging.Error("get default template failed: " + err.Error())
	}
	return tmpl
}

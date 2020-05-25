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

package egress_test

import (
	"fmt"
	"github.com/go-chassis/go-chassis/pkg/util/fileutil"
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/apache/servicecomb-mesher/proxy/cmd"
	mesherconfig "github.com/apache/servicecomb-mesher/proxy/config"
	_ "github.com/apache/servicecomb-mesher/proxy/control/istio"
	"github.com/apache/servicecomb-mesher/proxy/pkg/egress"
	"github.com/apache/servicecomb-mesher/proxy/pkg/egress/archaius"
	"github.com/go-chassis/go-chassis/control"
	_ "github.com/go-chassis/go-chassis/control/servicecomb"
	"github.com/go-chassis/go-chassis/core/config"
	"github.com/go-chassis/go-chassis/core/lager"
	"gopkg.in/yaml.v2"
)

func init() {
	lager.Init(&lager.Options{LoggerLevel: "DEBUG"})
}
func BenchmarkMatch(b *testing.B) {
	chassis := []byte(`
cse:
  service:
    registry:
      #disabled: false           optional:禁用注册发现选项，默认开始注册发现
      type: servicecenter           #optional:可选zookeeper/servicecenter，zookeeper供中软使用，不配置的情况下默认为servicecenter
      scope: full                   #optional:scope不为full时，只允许在本app间访问，不允许跨app访问；为full就是注册时允许跨app，并且发现本租户全部微服务
      address: http://127.0.0.1:30100
      #register: manual          optional：register不配置时默认为自动注册，可选参数有自动注册auto和手动注册manual
  
`)
	d, _ := os.Getwd()
	filename1 := filepath.Join(d, "chassis.yaml")
	f1, err := os.OpenFile(filename1, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0666)
	_, err = f1.Write(chassis)
	cmd.Init()
	config.Init()
	mesherconfig.Init()
	egress.Init()
	opts := control.Options{
		Infra:   config.GlobalDefinition.Panel.Infra,
		Address: config.GlobalDefinition.Panel.Settings["address"],
	}
	control.Init(opts)
	var yamlContent = `---
egress:
  infra: cse # pilot or cse
  address: http://istio-pilot.istio-system:15010
egressRule:
  google-ext:
    - hosts:
        - "www.google.com"
        - "*.yahoo.com"
      ports:
        - port: 80
          protocol: HTTP
  facebook-ext:
    - hosts:
        - "www.facebook.com"
      ports:
        - port: 80
          protocol: HTTP`

	ss := mesherconfig.EgressConfig{}
	err = yaml.Unmarshal([]byte(yamlContent), &ss)
	if err != nil {
		fmt.Println("unmarshal failed")
	}
	archaius.SetEgressRule(ss.Destinations)

	myString := "www.google.com"
	for i := 0; i < b.N; i++ {
		egress.Match(myString)
	}
}

func TestMatch(t *testing.T) {

	b := []byte(`
cse:
  service:
    registry:
      #disabled: false           optional:禁用注册发现选项，默认开始注册发现
      type: servicecenter           #optional:可选zookeeper/servicecenter，zookeeper供中软使用，不配置的情况下默认为servicecenter
      scope: full                   #optional:scope不为full时，只允许在本app间访问，不允许跨app访问；为full就是注册时允许跨app，并且发现本租户全部微服务
      address: http://127.0.0.1:30100
      #register: manual          optional：register不配置时默认为自动注册，可选参数有自动注册auto和手动注册manual
  
`)
	d, _ := os.Getwd()
	filename1 := filepath.Join(d, "chassis.yaml")
	f1, err := os.OpenFile(filename1, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0666)
	assert.NoError(t, err)
	_, err = f1.Write(b)
	b = []byte(`
---
#微服务的私有属性
#APPLICATION_ID: CSE #optional
service_description:
  name: Client
  #version: 0.1 #optional

`)
	d, _ = os.Getwd()
	filename1 = filepath.Join(d, "microservice.yaml")
	os.Remove(filename1)
	f1, err = os.Create(filename1)
	assert.NoError(t, err)
	defer f1.Close()
	_, err = io.WriteString(f1, string(b))
	assert.NoError(t, err)
	os.Setenv(fileutil.ChassisConfDir, d)
	cmd.Init()
	err = config.Init()
	err = mesherconfig.Init()
	err = egress.Init()
	//control.Init()
	var yamlContent = `---
egress:
  infra: cse # pilot or cse
  address: http://istio-pilot.istio-system:15010
egressRule:
  google-ext:
    - hosts:
        - "www.google.com"
        - "*.yahoo.com"
      ports:
        - port: 80
          protocol: HTTP
  facebook-ext:
    - hosts:
        - "www.facebook.com"
      ports:
        - port: 80
          protocol: HTTP`

	ss := mesherconfig.EgressConfig{}
	err = yaml.Unmarshal([]byte(yamlContent), &ss)
	if err != nil {
		fmt.Println("unmarshal failed")
	}
	archaius.SetEgressRule(ss.Destinations)

	myString := "www.google.com"
	c, _ := egress.Match(myString)
	if c == false {
		t.Errorf("Expected true but got false")
	}
	myString = "*.yahoo.com"
	c, _ = egress.Match(myString)
	if c == false {
		t.Errorf("Expected true but got false")
	}
}

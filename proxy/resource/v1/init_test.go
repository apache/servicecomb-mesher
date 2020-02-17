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

package v1

import (
	"testing"

	mesherconfig "github.com/apache/servicecomb-mesher/proxy/config"
	"github.com/go-chassis/go-chassis/core/lager"
	_ "github.com/go-chassis/go-chassis/core/router/servicecomb"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

var globalConfig = `
---
cse:
  loadbalance:
    strategyName: RoundRobin
  service:
    registry:
      address: http://10.162.197.14:30100
      scope: full
      watch: true
  protocols:
    http:
      listenAddress: 127.0.0.1:30101
  handler:
    chain:
      consumer:
        income:  ratelimiter-provider,local-selection
`
var mesherConf = `
routeRule:
  ShoppingCart:
    - precedence: 2
      route:
      - tags:
          version: 1.2
          app: HelloWorld
        weight: 80
      - tags:
          version: 1.3
          app: HelloWorld
        weight: 20
      match:
        refer: vmall-with-special-header
        source: vmall
        sourceTags:
            version: v2
        httpHeaders:
            cookie:
              regex: "^(.*?;)?(user=jason)(;.*)?$"
            X-Age:
              exact: "18"
    - precedence: 1
      route:
      - tags:
          version: 1.0
        weight: 100
`

func init() {
	lager.Init(&lager.Options{LoggerLevel: "DEBUG"})
}
func TestInit(t *testing.T) {
	t.Log("testing mesher admin protocol when protocol URI is valid")
	assert := assert.New(t)
	mesherConfig := new(mesherconfig.MesherConfig)
	yaml.Unmarshal([]byte(mesherConf), mesherConfig)
	mesherconfig.SetConfig(mesherConfig)
	err := Init()
	assert.Nil(err)
}

func TestInit2(t *testing.T) {
	t.Log("testing mesher admin protocol when protocol URI is not valid")
	assert := assert.New(t)
	mesherConfig := new(mesherconfig.MesherConfig)
	yaml.Unmarshal([]byte(mesherConf), mesherConfig)
	mesherConfig.Admin.ServerURI = "INVALID"
	mesherconfig.SetConfig(mesherConfig)
	err := Init()
	assert.Nil(err)
}

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

package config_test

import (
	//	"github.com/go-chassis/go-chassis/core/archaius"
	//	cConfig "github.com/go-chassis/go-chassis/core/config"
	//	"github.com/go-chassis/go-chassis/core/lager"
	//	"github.com/go-chassis/go-chassis/pkg/util/fileutil"
	"github.com/apache/servicecomb-mesher/proxy/cmd"
	"github.com/apache/servicecomb-mesher/proxy/config"
	"github.com/go-chassis/go-archaius"
	"github.com/go-chassis/go-chassis/core/lager"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
	//	"os"
	//	"path/filepath"
	"testing"
)

func init() {
	lager.Init(&lager.Options{LoggerLevel: "DEBUG"})
}
func TestGetConfigFilePath(t *testing.T) {
	var key = "mesher.yaml"
	archaius.Init(archaius.WithENVSource())
	cmd.Init()
	err := config.Init()
	assert.Error(t, err)
	f, _ := config.GetConfigFilePath(key)
	assert.Contains(t, f, key)
}

var file = []byte(`
localHealthCheck:
  - port: 8800
    protocol: rest
    uri: /health
    interval: 30s
    match:
      status: 200
      body: ok
pprof:
  enable: true
  listen: 0.0.0.0:6060
plugin:
  destinationResolver:
    http: host # how to turn host to destination name. default to service name，
    grpc: ip
  `)

func TestSetConfig(t *testing.T) {
	c := &config.MesherConfig{}
	if err := yaml.Unmarshal([]byte(file), c); err != nil {
		t.Error(err)
	}
	assert.Equal(t, "host", c.Plugin.DestinationResolver["http"])
	assert.Equal(t, "8800", c.HealthCheck[0].Port)
}

var egressFile = []byte(`
egress:
  infra: cse  # pilot or cse
  address: http://istio-pilot.istio-system:15010
  `)

func TestGetEgressEndpoints(t *testing.T) {
	config.Init()
	c := config.GetEgressConfig()
	if err := yaml.Unmarshal([]byte(egressFile), c); err != nil {
		t.Error(err)
	}

	assert.Equal(t, "http://istio-pilot.istio-system:15010", c.Egress.Address)
}

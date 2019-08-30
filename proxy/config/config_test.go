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
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
	//	"os"
	//	"path/filepath"
	"testing"
)

func TestGetConfigFilePath(t *testing.T) {
	var key = "mesher.yaml"
	cmd.Init()
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

// Testcase is trying to create files inside /tmp/build folder which is dynamic, so in travis it is not possible to create folder in prior, so can't test this case in travis
/*func TestInit(t *testing.T) {
	s, _ := fileutil.GetWorkDir()
	os.Setenv(fileutil.ChassisHome, s)
	chassisConf := filepath.Join(os.Getenv(fileutil.ChassisHome), "conf")
	os.MkdirAll(chassisConf, 0600)
	f, err := os.Create(filepath.Join(chassisConf, "chassis.yaml"))
	assert.NoError(t, err)
	t.Log(f.Name())

	f, err = os.Create(filepath.Join(chassisConf, "microservice.yaml"))
	t.Log(f.Name())
	assert.NoError(t, err)
	err = cConfig.Init()
	f, err = os.Create(filepath.Join(chassisConf, "mesher.yaml"))
	t.Log(f.Name())
	f.Write(file)
	f.Close()
	lager.Initialize("", "INFO", "", "size", true, 1, 10, 7)
	archaius.Init()

	err = config.Init()
	assert.NoError(t, err)
	t.Log(config.GetConfig())
	assert.Equal(t, "host", config.GetConfig().Plugin.DestinationResolver)
	assert.Equal(t, true, config.GetConfig().PProf.Enable)
	assert.Equal(t, "0.0.0.0:6060", config.GetConfig().PProf.Listen)
	assert.Equal(t, "rest", config.GetConfig().HealthCheck[0].Port)
}*/

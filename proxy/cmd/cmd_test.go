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

package cmd_test

import (
	"github.com/apache/servicecomb-mesher/proxy/cmd"
	"github.com/apache/servicecomb-mesher/proxy/common"
	"github.com/go-chassis/go-chassis/core/lager"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestParseConfigFromCmd(t *testing.T) {
	config := "/mesher.yaml"

	t.Log("========cmd --config=", config)
	os.Args = []string{"test", "--config", config}
	err := cmd.Init()
	configFromCmd := cmd.Configs
	assert.Equal(t, config, configFromCmd.ConfigFile)
	assert.Nil(t, err)

}

func TestConfigFromCmd_GeneratePortsMap(t *testing.T) {

	c := &cmd.ConfigFromCmd{
		LocalServicePorts: "rest:80,grpc:8000",
	}
	c.GeneratePortsMap()
	t.Log(c.PortsMap)
	assert.Equal(t, "127.0.0.1:80", c.PortsMap["rest"])
}
func TestConfigFromCmd_GeneratePortsMap2(t *testing.T) {

	c := &cmd.ConfigFromCmd{
		LocalServicePorts: "rest: 80,grpc",
	}
	err := c.GeneratePortsMap()
	t.Log(c.PortsMap)
	assert.Error(t, err)
}
func TestConfigFromCmd_GeneratePortsMap3(t *testing.T) {
	os.Setenv(common.EnvServicePorts, "rest:80,grpc:90")
	cmd.Init()
	_ = cmd.Configs.GeneratePortsMap()
	t.Log(cmd.Configs.PortsMap)
	assert.Equal(t, "127.0.0.1:80", cmd.Configs.PortsMap["rest"])
}
func TestConfigFromCmd_GeneratePortsMap4(t *testing.T) {
	os.Setenv(common.EnvSpecificAddr, "127.0.0.1:80")
	cmd.Init()
	_ = cmd.Configs.GeneratePortsMap()
	t.Log(cmd.Configs.PortsMap)
	assert.Equal(t, "127.0.0.1:80", cmd.Configs.PortsMap["rest"])
}
func init() {
	lager.Init(&lager.Options{LoggerLevel: "DEBUG"})
}

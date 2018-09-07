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

package health

import (
	"errors"
	"github.com/go-chassis/go-cc-client/configcenter-client"
	"github.com/go-chassis/go-chassis/core/config"
	"github.com/go-chassis/go-chassis/core/config/model"
	"github.com/go-chassis/go-chassis/core/lager"
	"github.com/go-chassis/go-chassis/core/registry"
	"github.com/go-chassis/go-chassis/core/registry/mock"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

var (
	mockError = errors.New("test mock error")
)

func TestGetMesherHealth(t *testing.T) {
	testGetServiceStatusSuccess(t)
	testGetServiceStatusFailed(t)

	testConfigCenterConnectSuccess(t)
	testConfigCenterConnectFailed(t)

	t.Log("mesher not connected to sc, not connected to configcenter")
	testGetServiceStatusFailed(t)
	resp := GetMesherHealth()
	assert.Equal(t, resp.ConnectedConfigCenterClient, false)
	assert.Equal(t, resp.ConnectedMonitoring, false)
	assert.Equal(t, resp.Status, Red)
	assert.NotEmpty(t, resp.Error)
}

func testInit() {
	lager.Initialize("", "INFO", "", "size", true, 1, 10, 7)
	p := filepath.Join(os.Getenv("GOPATH"), "src", "github.com", "go-chassis", "mesher", "conf")
	os.Setenv("CHASSIS_CONF_DIR", p)
	err := config.Init()
	if err != nil {
		config.GlobalDefinition = new(model.GlobalCfg)
	}
	config.MicroserviceDefinition = new(model.MicroserviceCfg)
}

func testGetServiceStatusSuccess(t *testing.T) {
	t.Log("mesher connected with SC")
	testInit()

	appId, microserviceName, version := "TestAPP", "TestService", "0.0.1"
	microserviceId, instanceId := "testMicroserviceId", "testInstanceId"
	mockinstances := []*registry.MicroServiceInstance{
		&registry.MicroServiceInstance{
			InstanceID: instanceId,
			ServiceID:  microserviceId,
		},
	}
	config.GlobalDefinition.AppID, config.SelfServiceName, config.SelfVersion = appId, microserviceName, version
	testRegistryObj := new(mock.RegistratorMock)
	registry.DefaultRegistrator = testRegistryObj

	testDiscoveryObj := new(mock.DiscoveryMock)
	registry.DefaultServiceDiscoveryService = testDiscoveryObj
	testDiscoveryObj.On("GetMicroServiceID", appId, microserviceName, version, "").Return(microserviceId, nil)
	testDiscoveryObj.On("GetMicroServiceInstances", microserviceId, microserviceId).Return(mockinstances, nil)
	testRegistryObj.On("Heartbeat", microserviceId, instanceId).Return(true, nil)

	respServiceName, respVersion, err := getServiceStatus()
	assert.Equal(t, respServiceName, microserviceName)
	assert.Equal(t, respVersion, version)
	assert.Nil(t, err)
}

func testGetServiceStatusFailed(t *testing.T) {
	t.Log("mesher not connected with SC")
	testInit()

	appId, microserviceName, version := "TestAPP", "TestService", "0.0.1"
	microserviceId := "testMicroserviceId"
	config.GlobalDefinition.AppID, config.SelfServiceName, config.SelfVersion = appId, microserviceName, version
	testDiscoveryObj := new(mock.DiscoveryMock)
	registry.DefaultServiceDiscoveryService = testDiscoveryObj
	testDiscoveryObj.On("GetMicroServiceID", appId, microserviceName, version, "").Return(microserviceId, mockError)

	respServiceName, respVersion, err := getServiceStatus()
	assert.Equal(t, respServiceName, microserviceName)
	assert.Equal(t, respVersion, version)
	assert.Equal(t, err, mockError)
}

func testConfigCenterConnectSuccess(t *testing.T) {
	testInit()
	t.Log("config center connected")

	testMemberDiscoverObj := new(MockMemberDiscovery)
	configcenterclient.MemberDiscoveryService = testMemberDiscoverObj
	testMemberDiscoverObj.On("RefreshMembers").Return(nil)
	testMemberDiscoverObj.On("GetConfigServer").Return([]string{"localhost:8080"}, nil)
	assert.True(t, isConfigCenterConnected())
}

func testConfigCenterConnectFailed(t *testing.T) {
	testInit()
	t.Log("config center not connected")

	testMemberDiscoverObj := new(MockMemberDiscovery)
	configcenterclient.MemberDiscoveryService = testMemberDiscoverObj
	testMemberDiscoverObj.On("GetConfigServer").Return([]string{"2.2.2.2:2222"}, mockError)
	assert.False(t, isConfigCenterConnected())
}

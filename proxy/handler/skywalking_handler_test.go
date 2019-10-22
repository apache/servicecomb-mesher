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

package handler_test

import (
	"context"
	"github.com/apache/servicecomb-mesher/proxy/config"
	mhandler "github.com/apache/servicecomb-mesher/proxy/handler"
	"github.com/apache/servicecomb-mesher/proxy/pkg/skywalking"
	gcconfig "github.com/go-chassis/go-chassis/core/config"
	"github.com/go-chassis/go-chassis/core/config/model"
	"github.com/go-chassis/go-chassis/core/handler"
	"github.com/go-chassis/go-chassis/core/invocation"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	Port      = ":49800"
	ServerUrl = "127.0.0.1:49800"
)

//initGcConfig
func initGcConfig() {
	var micCfg model.MicroserviceCfg
	micCfg.ServiceDescription.Name = "TEST"
	gcconfig.MicroserviceDefinition = &micCfg
}

//initMesherConfig
func initMesherConfig() {
	config.SetConfig(&config.MesherConfig{ServiceComb: &config.ServiceComb{config.APM{config.Tracing{Enable: true, ServerURI: "192.168.0.1:17289"}}}})
}

//initInv
func initInv() *invocation.Invocation {
	var i invocation.Invocation
	i.MicroServiceName = "test"
	i.Ctx = context.Background()
	i.Endpoint = "calculator"
	i.URLPathFormat = "/bmi"
	return &i
}

//TestProviderHandlerName
func TestProviderHandlerName(t *testing.T) {
	h := mhandler.SkyWalkingProviderHandler{}
	assert.Equal(t, h.Name(), skywalking.SkyWalkingProvider)
}

//TestNewProvier
func TestNewProvier(t *testing.T) {
	h := mhandler.NewSkyWalkingProvier()
	assert.NotEqual(t, h, nil)
	assert.Equal(t, h.Name(), skywalking.SkyWalkingProvider)
}

//TestProvierHandle
func TestProvierHandle(t *testing.T) {
	initGcConfig()
	initMesherConfig()
	skywalking.Init()
	c := handler.Chain{}
	c.AddHandler(mhandler.NewSkyWalkingProvier())

	gcconfig.GlobalDefinition = &model.GlobalCfg{}
	gcconfig.GlobalDefinition.Cse.Handler.Chain.Consumer = make(map[string]string)
	gcconfig.GlobalDefinition.Cse.Handler.Chain.Consumer["skywalking-provider"] = "skywalking-provider"

	c.Next(initInv(), func(r *invocation.Response) error {
		assert.Equal(t, r.Err, nil)
		return r.Err
	})
}

//TestConsumerHandlerName
func TestConsumerHandlerName(t *testing.T) {
	c := mhandler.SkyWalkingConsumerHandler{}
	assert.Equal(t, c.Name(), skywalking.SkyWalkingConsumer)
}

//TestNewConsumer
func TestNewConsumer(t *testing.T) {
	h := mhandler.NewSkyWalkingConsumer()
	assert.NotEqual(t, h, nil)
	assert.Equal(t, h.Name(), skywalking.SkyWalkingConsumer)
}

//TestConsumerHandle
func TestConsumerHandle(t *testing.T) {
	initGcConfig()
	initMesherConfig()
	skywalking.Init()
	c := handler.Chain{}
	c.AddHandler(mhandler.NewSkyWalkingConsumer())

	gcconfig.GlobalDefinition = &model.GlobalCfg{}
	gcconfig.GlobalDefinition.Cse.Handler.Chain.Consumer = make(map[string]string)
	gcconfig.GlobalDefinition.Cse.Handler.Chain.Consumer["skywalking-consumer"] = "skywalking-consumer"

	c.Next(initInv(), func(r *invocation.Response) error {
		assert.Equal(t, r.Err, nil)
		return r.Err
	})
}

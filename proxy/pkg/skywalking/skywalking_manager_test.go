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

package skywalking_test

import (
	"context"
	"github.com/apache/servicecomb-mesher/proxy/config"
	"github.com/apache/servicecomb-mesher/proxy/pkg/skywalking"
	gcconfig "github.com/go-chassis/go-chassis/core/config"
	"github.com/go-chassis/go-chassis/core/config/model"
	"github.com/go-chassis/go-chassis/core/invocation"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	Port      = ":49800"
	ServerUrl = "127.0.0.1:49800"
)

//initConfig
func initConfig() {
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

//TestInit init skywalking manager
func TestInit(t *testing.T) {
	initConfig()
	initMesherConfig()
	skywalking.Init()
	assert.NotEqual(t, gcconfig.MicroserviceDefinition, nil)
}

//TestCreateEntrySpan test skywalking manager creating entryspan
func TestCreateEntrySpan(t *testing.T) {
	initConfig()
	initMesherConfig()
	skywalking.Init()
	span, _, err := skywalking.CreateEntrySpan(initInv())
	assert.Equal(t, err, nil)
	assert.NotEqual(t, span, nil)
	span.End()
}

//TestCreateExitSpan test skywalking manager creating endspan
func TestCreateExitSpan(t *testing.T) {
	initConfig()
	initMesherConfig()
	skywalking.Init()
	inv := initInv()
	span, ctx, err := skywalking.CreateEntrySpan(inv)
	assert.Equal(t, err, nil)
	assert.NotEqual(t, span, nil)
	spanExit, err := skywalking.CreateExitSpan(ctx, inv)
	assert.Equal(t, err, nil)
	assert.NotEqual(t, spanExit, nil)
	spanExit.End()
	span.End()
}

//TestCreateLocalSpan test skywalking manager creating localspan
func TestCreateLocalSpan(t *testing.T) {
	initConfig()
	initMesherConfig()
	skywalking.Init()
	span, _, err := skywalking.CreateLocalSpan(context.Background())
	assert.Equal(t, err, nil)
	assert.NotEqual(t, span, nil)
}

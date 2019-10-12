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
	"github.com/apache/servicecomb-mesher/proxy/pkg/apm"
	"github.com/apache/servicecomb-mesher/proxy/pkg/apm/skywalking"
	"github.com/go-chassis/go-chassis/core/config"
	"github.com/go-chassis/go-chassis/core/config/model"
	"github.com/go-chassis/go-chassis/core/invocation"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	Port      = ":49800"
	ServerUrl = "127.0.0.1:49800"
)

func InitConfig() {
	var micCfg model.MicroserviceCfg
	micCfg.ServiceDescription.Name = "TEST"
	config.MicroserviceDefinition = &micCfg
}

func InitInv() *invocation.Invocation {
	var i invocation.Invocation
	i.MicroServiceName = "test"
	i.Ctx = context.Background()
	i.Endpoint = "calculator"
	i.URLPathFormat = "/bmi"
	return &i
}

func TestInit(t *testing.T) {
	InitConfig()
	_, err := skywalking.NewApmClient(apm.Options{ServerUri: ServerUrl})
	assert.Equal(t, err, nil)
}

func TestCreateEntrySpan(t *testing.T) {
	InitConfig()
	apmClient, err := skywalking.NewApmClient(apm.Options{ServerUri: ServerUrl})
	span, err := apmClient.CreateEntrySpan(InitInv())
	assert.Equal(t, err, nil)
	assert.NotEqual(t, span, nil)
	err = apmClient.EndSpan(span, 200)
}

func TestCreateExitSpan(t *testing.T) {
	InitConfig()
	apmClient, err := skywalking.NewApmClient(apm.Options{ServerUri: ServerUrl})
	inv := InitInv()
	span, err := apmClient.CreateEntrySpan(inv)
	assert.Equal(t, err, nil)
	assert.NotEqual(t, span, nil)
	spanExit, err := apmClient.CreateExitSpan(inv)
	assert.Equal(t, err, nil)
	assert.NotEqual(t, spanExit, nil)
	err = apmClient.EndSpan(spanExit, 200)
	assert.Equal(t, err, nil)
	err = apmClient.EndSpan(span, 200)
	assert.Equal(t, err, nil)
}

func TestEndSpan(t *testing.T) {
	InitConfig()
	apmClient, err := skywalking.NewApmClient(apm.Options{ServerUri: ServerUrl})
	span, err := apmClient.CreateEntrySpan(InitInv())
	assert.Equal(t, err, nil)
	assert.NotEqual(t, span, nil)
	err = apmClient.EndSpan(span, 200)
	assert.Equal(t, err, nil)
}

func TestCreateSpans(t *testing.T) {
	InitConfig()
	apmClient, err := skywalking.NewApmClient(apm.Options{ServerUri: ServerUrl})
	spans, err := apmClient.CreateSpans(InitInv())
	assert.Equal(t, err, nil)
	assert.NotEqual(t, spans, nil)
}

func TestEndSpans(t *testing.T) {
	InitConfig()
	apmClient, err := skywalking.NewApmClient(apm.Options{ServerUri: ServerUrl})
	spans, err := apmClient.CreateSpans(InitInv())
	assert.Equal(t, err, nil)
	assert.NotEqual(t, spans, nil)
	err = apmClient.EndSpans(spans, 200)
	assert.Equal(t, err, nil)
}

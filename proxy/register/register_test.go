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

package register

import (
	"testing"

	"github.com/apache/servicecomb-mesher/proxy/common"
	chassisCommon "github.com/go-chassis/go-chassis/core/common"
	"github.com/go-chassis/go-chassis/core/config"
	"github.com/go-chassis/go-chassis/core/config/model"
	"github.com/go-chassis/go-chassis/core/lager"
	"github.com/go-chassis/go-chassis/core/registry"
	"github.com/stretchr/testify/assert"
)

func init() {
	lager.Init(&lager.Options{LoggerLevel: "DEBUG"})
}
func TestAdaptEndpoints(t *testing.T) {
	protoMap := make(map[string]model.Protocol)
	config.GlobalDefinition = &model.GlobalCfg{
		Cse: model.CseStruct{
			Protocols: protoMap,
		},
	}

	AdaptEndpoints()
	assert.NotNil(t, registry.InstanceEndpoints)

	protoMap[chassisCommon.ProtocolRest] = model.Protocol{
		Advertise: "1.1.1.1:8080",
	}
	AdaptEndpoints()
	assert.NotNil(t, registry.InstanceEndpoints)

	protoMap[common.HTTPProtocol] = model.Protocol{
		Advertise: "1.1.1.1:8081",
	}
	delete(protoMap, chassisCommon.ProtocolRest)
	AdaptEndpoints()
	assert.Equal(t, 1, len(registry.InstanceEndpoints))
	_, ok := registry.InstanceEndpoints[common.HTTPProtocol]
	assert.False(t, ok)
	endpoint0 := registry.InstanceEndpoints[chassisCommon.ProtocolRest]
	endpoint1 := protoMap[common.HTTPProtocol].Advertise
	assert.Equal(t, endpoint0, endpoint1)
}

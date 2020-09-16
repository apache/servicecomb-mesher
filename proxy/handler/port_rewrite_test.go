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

package handler

import (
	"testing"

	"github.com/go-chassis/go-chassis/v2/core/common"
	"github.com/go-chassis/go-chassis/v2/core/config"
	"github.com/go-chassis/go-chassis/v2/core/config/model"
	"github.com/go-chassis/go-chassis/v2/core/handler"
	"github.com/go-chassis/go-chassis/v2/core/invocation"
	"github.com/stretchr/testify/assert"
)

func TestPortRewriteHandler_ValidEndpoint(t *testing.T) {
	t.Log("testing port rewrite handler with valid endpoint")

	c := handler.Chain{}
	c.AddHandler(&PortSelectionHandler{})

	config.GlobalDefinition = &model.GlobalCfg{}
	config.GlobalDefinition.ServiceComb.Handler.Chain.Consumer = make(map[string]string)
	config.GlobalDefinition.ServiceComb.Handler.Chain.Consumer["outgoing"] = PortMapForPilot
	i := &invocation.Invocation{
		MicroServiceName: "service1",
		SchemaID:         "schema1",
		OperationID:      "SayHello",
		Endpoint:         "127.0.0.1:5555",
	}

	c.Next(i, func(r *invocation.Response) {
		assert.NoError(t, r.Err)
	})
}

func TestPortRewriteHandler_InValidEndpoint(t *testing.T) {
	t.Log("testing port rewrite handler with empty endpoint")

	c := handler.Chain{}
	c.AddHandler(&PortSelectionHandler{})

	config.GlobalDefinition = &model.GlobalCfg{}
	config.GlobalDefinition.ServiceComb.Handler.Chain.Consumer = make(map[string]string)
	config.GlobalDefinition.ServiceComb.Handler.Chain.Consumer["outgoing"] = PortMapForPilot
	i := &invocation.Invocation{
		MicroServiceName: "service1",
		SchemaID:         "schema1",
		OperationID:      "SayHello",
		Endpoint:         "",
	}

	c.Next(i, func(r *invocation.Response) {
		assert.Error(t, r.Err)
	})

}

func TestPortRewriteHandler_Names(t *testing.T) {
	handlerObject := &PortSelectionHandler{}
	name := handlerObject.Name()
	assert.Equal(t, PortMapForPilot, name)
}

func TestReplacePort_InvalidEndpoint(t *testing.T) {
	output, err := replacePort("grpc", "")
	assert.Error(t, err)
	assert.Equal(t, "", output)
}

func TestReplacePort_ValidEndpoint(t *testing.T) {
	output, err := replacePort(common.ProtocolRest, "127.0.0.1:80")
	assert.Equal(t, "127.0.0.1:30101", output)
	assert.NoError(t, err)
}

func BenchmarkReplacePort(b *testing.B) {
	for i := 0; i < b.N; i++ {
		replacePort(common.ProtocolRest, "127.0.0.1:80")
	}
}

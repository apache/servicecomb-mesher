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

package schema

import (
	"github.com/go-chassis/go-chassis/core/config"
	"github.com/go-chassis/go-chassis/core/lager"
	"github.com/go-chassis/go-chassis/core/registry"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func init() {
	lager.Init(&lager.Options{LoggerLevel: "DEBUG"})
}

func Test_GetRspSchema(t *testing.T) {
	res := make(map[string]*MethRespond, 0)
	res["200"] = &MethRespond{
		Status: "200",
	}
	m := DefMethod{
		ownerSvc: "svc",
		Path:     "/test/StringArray",
		OperaID:  "StringArray",
		Responds: res,
	}

	assert.NotNil(t, m.GetRspSchema(200))

	assert.Nil(t, m.GetRspSchema(404))
}

func Test_GetParamNameAndWhere(t *testing.T) {
	res := make(map[string]*MethRespond, 0)
	res["200"] = &MethRespond{
		Status: "200",
	}

	mParams := []MethParam{{
		Name:  "para1",
		Where: "query",
		Indx:  0,
	}, {
		Name:  "para2",
		Where: "body",
		Indx:  1,
	},
	}

	m := DefMethod{
		ownerSvc: "svc",
		Path:     "/test/StringArray",
		OperaID:  "StringArray",
		Responds: res,
		Paras:    mParams,
	}

	// case in query
	name, indx := m.GetParamNameAndWhere(0)
	assert.Equal(t, "para1", name)
	assert.Equal(t, InQuery, indx)

	// case in body
	name, indx = m.GetParamNameAndWhere(1)
	assert.Equal(t, "para2", name)
	assert.Equal(t, InBody, indx)

	// other
	name, indx = m.GetParamNameAndWhere(2)
	assert.Equal(t, 0, len(name))
	assert.Equal(t, InQuery, indx)
}

func Test_GetParamSchema(t *testing.T) {
	res := make(map[string]*MethRespond, 0)
	res["200"] = &MethRespond{
		Status: "200",
	}

	mParams := []MethParam{{
		Name:  "para1",
		Where: "query",
		Indx:  0,
	}, {
		Name:  "para2",
		Where: "body",
		Indx:  1,
	},
	}

	m := DefMethod{
		ownerSvc: "svc",
		Path:     "/test/StringArray",
		OperaID:  "StringArray",
		Responds: res,
		Paras:    mParams,
	}

	// case in query
	param := m.GetParamSchema(0)
	assert.Equal(t, "para1", param.Name)

	// case in body
	param = m.GetParamSchema(1)
	assert.Equal(t, "para2", param.Name)

	// other
	param = m.GetParamSchema(2)
	assert.Nil(t, param)

}

// CovertSwaggerMethordToLocalMethord(&schema, &m, &meth)
func Test_CovertSwaggerMethordToLocalMethord(t *testing.T) {
	schema := &registry.SchemaContent{
		Definition: map[string]registry.Definition{
			"hello": {},
		},
	}
	paras := make([]registry.Parameter, 0)
	paras = append(paras, registry.Parameter{
		Name: "Hello",
		Type: "string",
		Schema: registry.SchemaValue{
			Type:      "string",
			Reference: "hello",
		},
	}, registry.Parameter{
		Name: "Hello1",
		Type: "",
		Schema: registry.SchemaValue{
			Type:      "string",
			Reference: "hello1",
		},
	}, registry.Parameter{
		Name: "Hello2",
		Type: "",
		Schema: registry.SchemaValue{
			Type:      "",
			Reference: "hello1",
		},
	})

	srcMethod := &registry.MethodInfo{
		Parameters: paras,
		Response: map[string]registry.Response{
			"200": {
				Schema: map[string]string{"type": "string"},
			},
			"201": {
				Schema: map[string]string{"$ref": "/v/hello"},
			},
		},
	}
	distMeth := &DefMethod{}
	CovertSwaggerMethordToLocalMethord(schema, srcMethod, distMeth)
}

func Test_GetSvcByInterface(t *testing.T) {
	config.Init()

	registry.DefaultContractDiscoveryService = new(MockContractDiscoveryService)
	v := GetSvcByInterface("hello")
	assert.NotNil(t, v)

	svcToInterfaceCache.Set("hello", &registry.MicroService{}, 0)
	// case has value
	v = GetSvcByInterface("hello")
	assert.NotNil(t, v)

}

func Test_GetMethodByInterface(t *testing.T) {
	registry.DefaultContractDiscoveryService = new(MockContractDiscoveryService)
	GetMethodByInterface("hello", "hello")
}

// ContractDiscoveryService struct for disco mock
type MockContractDiscoveryService struct {
	mock.Mock
}

func (m *MockContractDiscoveryService) GetMicroServicesByInterface(interfaceName string) (microservices []*registry.MicroService) {
	microservices = append(microservices, &registry.MicroService{})
	return
}

func (m *MockContractDiscoveryService) GetSchemaContentByInterface(interfaceName string) registry.SchemaContent {
	return registry.SchemaContent{}
}

func (m *MockContractDiscoveryService) GetSchemaContentByServiceName(svcName, version, appID, env string) []*registry.SchemaContent {
	var sc []*registry.SchemaContent
	sc = append(sc, &registry.SchemaContent{
		Paths: map[string]map[string]registry.MethodInfo{
			"hello": {},
		},
	})
	return nil
}

func (m *MockContractDiscoveryService) Close() error {
	return nil
}

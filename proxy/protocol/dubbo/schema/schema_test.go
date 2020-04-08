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
	"github.com/stretchr/testify/assert"
	"testing"
)

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

	mParams := []MethParam{MethParam{
		Name:  "para1",
		Where: "query",
		Indx:  0,
	}, MethParam{
		Name:  "para2",
		Where: "body",
		Indx:  1,
	}}

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

	mParams := []MethParam{MethParam{
		Name:  "para1",
		Where: "query",
		Indx:  0,
	}, MethParam{
		Name:  "para2",
		Where: "body",
		Indx:  1,
	}}

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

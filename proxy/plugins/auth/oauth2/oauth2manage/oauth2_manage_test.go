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

package oauth2manage

import (
	"github.com/go-chassis/go-chassis/core/invocation"
	"testing"
)

type SubjectTest struct{}

func (s *SubjectTest) GrandTypeProcess(inv *invocation.Invocation, cb invocation.ResponseCallBack) error {
	return nil
}

// Test RegisterType
func TestRegisterType(t *testing.T) {
	t.Log("testing register type")

	//Test register plugin name with oauth2
	grandType := "password"
	subjectTest := new(SubjectTest)
	RegisterType(grandType, subjectTest)

	// Test NewType with authorization_code
	grandType = "authorization_code"
	subjectTest = new(SubjectTest)
	RegisterType(grandType, subjectTest)
	_, err := NewType(grandType)
	if err != nil {
		t.Error("GrandType need to register")
		return
	}
}

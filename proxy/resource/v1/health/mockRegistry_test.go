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
	"github.com/stretchr/testify/mock"
)

type MockMemberDiscovery struct {
	mock.Mock
}

func (m *MockMemberDiscovery) ConfigurationInit(initConfigServer []string) error {
	args := m.Called(initConfigServer)
	return args.Error(0)
}
func (m *MockMemberDiscovery) GetConfigServer() ([]string, error) {
	args := m.Called()
	return args.Get(0).([]string), args.Error(1)
}
func (m *MockMemberDiscovery) RefreshMembers() error {
	args := m.Called()
	return args.Error(0)
}
func (m *MockMemberDiscovery) Shuffle() error {
	args := m.Called()
	return args.Error(0)
}
func (m *MockMemberDiscovery) GetWorkingConfigCenterIP(entryPoint []string) ([]string, error) {
	args := m.Called(entryPoint)
	return args.Get(0).([]string), args.Error(0)
}

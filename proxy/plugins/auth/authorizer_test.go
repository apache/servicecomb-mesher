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

package auth_test

import (
	"context"
	"github.com/apache/servicecomb-mesher/proxy/plugins/auth"
	"github.com/go-chassis/go-chassis/core/invocation"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

var inv *invocation.Invocation
var cb invocation.ResponseCallBack

// testAuth is a null struct
type testAuth struct{}

// Authorize is an empty implementation
func (t testAuth) Authorize(inv *invocation.Invocation, cb invocation.ResponseCallBack) error {
	return nil
}

// Authenticate is an empty implementation
func (t testAuth) Authenticate(inv *invocation.Invocation, cb invocation.ResponseCallBack) error {
	return nil
}

// Test RegisterAuthPlugin
func TestRegisterAuthPlugin(t *testing.T) {
	//Test register plugin name with oauth2
	AuthName := "oauth2"
	auth.RegisterAuthPlugin(AuthName, testAuth{})

	//Test auth plugin name already exist
	auth.RegisterAuthPlugin(AuthName, testAuth{})

	// Test NewAuth with register basicAuth
	AuthName = "basicAuth"
	auth.RegisterAuthPlugin(AuthName, testAuth{})
	_, err := auth.NewAuth(AuthName)
	assert.NoError(t, err)

	// Test NewAuth with not register openID
	AuthName = "openID"
	_, err = auth.NewAuth(AuthName)
	assert.Error(t, err)
}

func TestGetHTTPRequest(t *testing.T) {
	inv = invocation.New(context.TODO())
	req := httptest.NewRequest(http.MethodPost, "/api", nil)
	inv.Args = req
	_, err := auth.GetHTTPRequest(inv)
	assert.NoError(t, err)
}

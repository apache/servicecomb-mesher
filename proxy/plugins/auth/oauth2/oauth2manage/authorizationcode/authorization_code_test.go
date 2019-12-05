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

package authorizationcode

import (
	"context"
	"errors"
	"github.com/apache/servicecomb-mesher/proxy/config"
	"github.com/apache/servicecomb-mesher/proxy/plugins/auth/oauth2/oauth2manage"
	"github.com/go-chassis/go-chassis/core/invocation"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
	"net/http"
	"testing"
)

var content = []byte(`
authorization:
  authName: oauth2
  endpoint:
    authURL:    "your auth url"
    tokenURL:   "your token url"
  client:
    clientID:      "your client ID"
    clientSecret:  "your client secret"
`)

var inv *invocation.Invocation
var cb invocation.ResponseCallBack

func init() {
	inv = invocation.New(context.TODO())
	conf := &config.MesherConfig{}
	err := yaml.Unmarshal(content, conf)
	if err != nil {
		errors.New("parse the config content failed")
		return
	}
	config.SetConfig(conf)
}

// Test authorization_code model
func TestAuthorizationCode_GrandTypeProcess(t *testing.T) {
	grandType, err := oauth2manage.NewType(GrandType)
	if err != nil {
		t.Errorf("grand type error: %s", err.Error())
		return
	}

	req, err := http.NewRequest(http.MethodPost, "https://api/user?code=test", nil)
	if err != nil {
		t.Errorf("authorization failed: %s", err.Error())
		return
	}
	inv.Args = req

	err = grandType.GrandTypeProcess(inv, cb)
	assert.Error(t, err)
}

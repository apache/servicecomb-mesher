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

package oauth2

import (
	"context"
	"errors"
	"github.com/apache/servicecomb-mesher/proxy/config"
	"github.com/apache/servicecomb-mesher/proxy/plugins/auth"
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
	conf := &config.MesherConfig{}
	err := yaml.Unmarshal(content, conf)
	if err != nil {
		errors.New("parse the config content failed")
		return
	}
	conf.Authorization.Type = AuthName
	config.SetConfig(conf)
	inv = invocation.New(context.TODO())
}

func TestOAuth2_Authorize(t *testing.T) {
	authority, err := auth.NewAuth(AuthName)
	if err != nil {
		t.Errorf("new a authorization object failed: %s ", err.Error())
		return
	}

	t.Run("Normal", func(t *testing.T) {
		req, err := http.NewRequest("POST", "https://api/?grand_type=authorization_code&code=test", nil)
		assert.NoError(t, err)
		inv.Args = req

		err = authority.Authorize(inv, cb)
		assert.Error(t, err)
	})

	t.Run("Invalid", func(t *testing.T) {
		req, err := http.NewRequest("POST", "https://api/", nil)
		assert.NoError(t, err)
		inv.Args = req

		err = authority.Authorize(inv, cb)
		assert.Error(t, err)
	})
}

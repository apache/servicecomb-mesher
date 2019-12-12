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

package auth

import (
	"context"
	"github.com/apache/servicecomb-mesher/proxy/handler/oauth2/oauth2manage"
	"github.com/go-chassis/go-chassis/core/invocation"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
	"net/http"
	"testing"
)

var inv *invocation.Invocation

func init() {
	inv = invocation.New(context.TODO())
}

// Test authorization_code model
func TestAuthorizationCode_GrantTypeProcess(t *testing.T) {
	grantType, err := oauth2manage.NewType(GrantType)
	if err != nil {
		t.Errorf("grant type error: %s", err.Error())
		return
	}

	var cb = func(r *invocation.Response) error {
		assert.Error(t, r.Err)
		return r.Err
	}

	Use(&Config{
		UseConfig: func() *oauth2.Config {
			return &oauth2.Config{}
		},
	})

	t.Run("Invalid protocol", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, "https://api/?grant_type=authorization_code&code=test", nil)
		if err != nil {
			t.Errorf("authorization failed: %s", err.Error())
			return
		}
		inv.Args = req
		_, err = grantType.GrantTypeProcess(inv, cb)
		assert.Error(t, err)
	})
	t.Run("Invalid code", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, "https://api/?grant_type=authorization_code&code=", nil)
		if err != nil {
			t.Errorf("authorization failed: %s", err.Error())
			return
		}
		inv.Args = req
		_, err = grantType.GrantTypeProcess(inv, cb)
		assert.Error(t, err)
	})
}

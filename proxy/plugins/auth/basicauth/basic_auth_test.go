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

package basicauth

import (
	"context"
	"github.com/apache/servicecomb-mesher/proxy/plugins/auth"
	"github.com/go-chassis/go-chassis/core/invocation"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

var inv *invocation.Invocation

//Test NewBasicAuth with valid user
func TestNewBasicAuth_Valid(t *testing.T) {
	Use(&Subject{
		Realm: "test-realm",
		UserInfo: func(u, p string) error {
			return nil
		},
	})

	var cb = func(r *invocation.Response) error {
		assert.Error(t, r.Err)
		return r.Err
	}

	auth, err := auth.NewAuth(AuthName)
	assert.NoError(t, err)
	inv = invocation.New(context.TODO())

	t.Run("Normal", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api", nil)
		req.Header.Add("Authorization", "Basic dGVzdDp0ZXN0")
		inv.Args = req
		err = auth.Authorize(inv, cb)
		assert.NoError(t, err)
	})

	t.Run("Invalid", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api", nil)
		req.Header.Add("Authorization", "dGVzdDp0ZXN0")
		inv.Args = req
		err = auth.Authorize(inv, cb)
		assert.Error(t, err)
	})

}

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
	mconf "github.com/apache/servicecomb-mesher/proxy/config"
	"github.com/apache/servicecomb-mesher/proxy/plugins/auth/noop"
	"github.com/go-chassis/go-chassis/core/config"
	"github.com/go-chassis/go-chassis/core/config/model"
	"github.com/go-chassis/go-chassis/core/handler"
	"github.com/go-chassis/go-chassis/core/invocation"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func initHandler() handler.Chain {
	mconf.SetConfig(&mconf.MesherConfig{
		Authorization: &mconf.Authorization{
			Type: noop.AuthName,
		},
	})

	c := handler.Chain{}
	c.AddHandler(&AuthHandler{})
	return c
}

func initInv() *invocation.Invocation {

	config.GlobalDefinition = &model.GlobalCfg{}
	config.GlobalDefinition.Cse.Handler.Chain.Provider = make(map[string]string)
	config.GlobalDefinition.Cse.Handler.Chain.Provider["outgoing"] = Authorizer

	var i *invocation.Invocation

	i = invocation.New(nil)
	i.MicroServiceName = "service1"
	i.SchemaID = "schema1"
	i.OperationID = "SayHello"
	i.Endpoint = ""

	return i
}

func TestAuthorizationHandler(t *testing.T) {
	c := initHandler()
	i := initInv()

	t.Run("header", func(t *testing.T) {
		i.SetHeader("Authorization", "Basic dGVzdDp0ZXN0")
		c.Next(i, func(r *invocation.Response) error {
			assert.Nil(t, r.Err)
			return r.Err
		})
	})

	t.Run("request", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, "https://api/?grand_type=authorization_code&code=test", nil)
		if err != nil {
			t.Errorf("authorization failed: %s", err.Error())
			return
		}
		i.Args = req

		c.Next(i, func(r *invocation.Response) error {
			assert.Nil(t, r.Err)
			return r.Err
		})
	})
}

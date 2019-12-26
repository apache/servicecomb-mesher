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
	"github.com/go-chassis/go-chassis/core/config"
	"github.com/go-chassis/go-chassis/core/config/model"
	"github.com/go-chassis/go-chassis/core/handler"
	"github.com/go-chassis/go-chassis/core/invocation"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
	"net/http"
	"testing"
)

func initHandler() handler.Chain {
	c := handler.Chain{}
	c.AddHandler(&Handler{})
	return c
}

func initInv() *invocation.Invocation {

	config.GlobalDefinition = &model.GlobalCfg{}
	config.GlobalDefinition.Cse.Handler.Chain.Provider = make(map[string]string)
	config.GlobalDefinition.Cse.Handler.Chain.Provider["outgoing"] = AuthName

	var i *invocation.Invocation

	i = invocation.New(nil)
	i.MicroServiceName = "service1"
	i.SchemaID = "schema1"
	i.OperationID = "SayHello"
	i.Endpoint = ""

	return i
}
func TestOAuth2_Handle(t *testing.T) {
	c := initHandler()
	i := initInv()

	Use(&OAuth2{
		GrantType: "authorization_code",
		Authenticate: func(at string, req *http.Request) error {
			return nil
		},
		UseConfig: &oauth2.Config{
			ClientID:     "",           // (required, string) your client_ID
			ClientSecret: "",           // (required, string) your client_Secret
			Scopes:       []string{""}, // (optional, string) scope specifies requested permissions
			RedirectURL:  "",           // (required, string) URL to redirect users going through the OAuth2 flow
			Endpoint: oauth2.Endpoint{ // (required, string) your auth server endpoint
				AuthURL:  "",
				TokenURL: "",
			},
		},
	})

	t.Run("Invalid grant_type", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, "https://api/?grant_type=test&code=test", nil)
		if err != nil {
			t.Errorf("authorization failed: %s", err.Error())
			return
		}
		i.Args = req

		i.SetHeader("Authorization", "Basic dGVzdDp0ZXN0")
		c.Next(i, func(r *invocation.Response) error {
			assert.Error(t, r.Err)
			return r.Err
		})
	})
	t.Run("null grant_type", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, "https://api/?grant_type=&code=test", nil)
		if err != nil {
			t.Errorf("authorization failed: %s", err.Error())
			return
		}
		i.Args = req

		i.SetHeader("Authorization", "Basic dGVzdDp0ZXN0")
		c.Next(i, func(r *invocation.Response) error {
			assert.NoError(t, r.Err)
			return r.Err
		})
	})

	t.Run("normal grant_type", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, "https://api/?grant_type=authorization_code&code=test", nil)
		if err != nil {
			t.Errorf("authorization failed: %s", err.Error())
			return
		}
		i.Args = req

		c.Next(i, func(r *invocation.Response) error {
			assert.NoError(t, r.Err)
			return r.Err
		})
	})

	t.Run("normal state", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, "https://api/?grant_type=&code=test&state=random", nil)
		if err != nil {
			t.Errorf("authorization failed: %s", err.Error())
			return
		}
		i.Args = req

		i.SetHeader("Authorization", "Basic dGVzdDp0ZXN0")
		c.Next(i, func(r *invocation.Response) error {
			assert.NoError(t, r.Err)
			return r.Err
		})
	})
}

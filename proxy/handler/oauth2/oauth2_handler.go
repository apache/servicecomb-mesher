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
	"errors"
	"github.com/apache/servicecomb-mesher/proxy/handler/oauth2/oauth2manage"
	"github.com/go-chassis/go-chassis/core/handler"
	"github.com/go-chassis/go-chassis/core/invocation"
	"github.com/go-mesh/openlogging"
	"net/http"

	// authorization_code model
	_ "github.com/apache/servicecomb-mesher/proxy/handler/oauth2/oauth2manage/auth"
)

// AuthName is a constant
const AuthName = "oauth2"

//Handler is is a oauth2 pre process raw data in handler
type Handler struct {
}

// Handle is provider
func (oa *Handler) Handle(chain *handler.Chain, inv *invocation.Invocation, cb invocation.ResponseCallBack) {
	if req, ok := inv.Args.(*http.Request); ok {
		grantType := req.FormValue("grant_type")
		if grantType == "" {
			_ = errors.New("can not fetch the grant_type")
			return
		}

		gt, err := oauth2manage.NewType(grantType)
		if err != nil {
			openlogging.Error("grant_type error: " + err.Error())
			cb(oauth2manage.InvResponse(http.StatusUnauthorized, err))
			return
		}
		accessToken, err := gt.GrantTypeProcess(inv, cb)
		if err != nil {
			openlogging.Error("authorization error: " + err.Error())
			cb(oauth2manage.InvResponse(http.StatusUnauthorized, err))
			return
		}
		if auth != nil && auth.Authenticate != nil {
			err = auth.Authenticate(accessToken, req)
			if err != nil {
				_ = errors.New("invalid authentication")
				cb(oauth2manage.InvResponse(http.StatusUnauthorized, err))
				return
			}
		}
	}
	chain.Next(inv, func(r *invocation.Response) error {
		return cb(r)
	})
}

// Name returns router string
func (oa *Handler) Name() string {
	return AuthName
}

// NewOAuth2 returns new auth handler
func NewOAuth2() handler.Handler {
	return &Handler{}
}

func init() {
	err := handler.RegisterHandler(AuthName, NewOAuth2)
	if err != nil {
		openlogging.Error("register handler error: " + err.Error())
		return
	}
}

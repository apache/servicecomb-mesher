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
	"github.com/go-chassis/go-chassis/core/handler"
	"github.com/go-chassis/go-chassis/core/invocation"
	"github.com/go-mesh/openlogging"
	"net/http"
	"time"
)

// errors
var (
	ErrNoGrandType  = errors.New("no grant_type found")
	ErrInvalidCode  = errors.New("invalid code")
	ErrInvalidToken = errors.New("invalid authorization")
	ErrInvalidAuth  = errors.New("invalid authentication")
	ErrExpiredToken = errors.New("expired token")
)

// AuthName is a constant
const AuthName = "oauth2"

// Handler is is a oauth2 pre process raw data in handler
type Handler struct {
}

// Handle is provider
func (oa *Handler) Handle(chain *handler.Chain, inv *invocation.Invocation, cb invocation.ResponseCallBack) {
	if req, ok := inv.Args.(*http.Request); ok {
		grantType := req.FormValue("grant_type")
		if grantType == "" {
			WriteBackErr(ErrNoGrandType, http.StatusUnauthorized, cb)
			return
		}

		if auth != nil && auth.GrantType == "authorization_code" {
			if req, ok := inv.Args.(*http.Request); ok {
				code := req.FormValue("code")
				if code == "" {
					WriteBackErr(ErrInvalidCode, http.StatusUnauthorized, cb)
					return
				}

				accessToken, err := getToken(code, cb)
				if err != nil {
					openlogging.Error("authorization error: " + err.Error())
					WriteBackErr(ErrInvalidToken, http.StatusUnauthorized, cb)
					return
				}

				if auth.Authenticate != nil {
					err = auth.Authenticate(accessToken, req)
					if err != nil {
						openlogging.Error("authentication error: " + err.Error())
						WriteBackErr(ErrInvalidAuth, http.StatusUnauthorized, cb)
						return
					}
				}
			}
		}
		chain.Next(inv, func(r *invocation.Response) error {
			return cb(r)
		})
	}
}

// getToken deal with the authorization code and return the token
func getToken(code string, cb invocation.ResponseCallBack) (accessToken string, err error) {
	if auth.UseConfig != nil {
		config := auth.UseConfig
		token, err := config.Exchange(context.Background(), code)
		if err != nil {
			openlogging.Error("get token failed, errors: " + err.Error())
			WriteBackErr(ErrInvalidCode, http.StatusUnauthorized, cb)
			return "", err
		}

		// set the expiry token in 30 minutes
		token.Expiry = time.Now().Add(30 * 60 * time.Second)
		if time.Now().After(token.Expiry) {
			return "", ErrExpiredToken
		}
		accessToken = token.AccessToken
		return accessToken, nil
	}
	return "", nil
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

// WriteBackErr write err and callback
func WriteBackErr(err error, status int, cb invocation.ResponseCallBack) {
	r := &invocation.Response{
		Err:    err,
		Status: status,
	}
	err = cb(r)
	if err != nil {
		openlogging.Error("response error: " + err.Error())
		return
	}
}

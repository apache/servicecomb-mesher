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
	"github.com/apache/servicecomb-mesher/proxy/config"
	"github.com/apache/servicecomb-mesher/proxy/plugins/auth"
	"github.com/go-chassis/go-chassis/core/handler"
	"github.com/go-chassis/go-chassis/core/invocation"
	"github.com/go-mesh/openlogging"
)

// Authorizer is a constant
const Authorizer = "authorizer"

// AuthHandler struct contains config info
type AuthHandler struct {
	config *config.Authorization
}

// Handle is provider
func (au *AuthHandler) Handle(chain *handler.Chain, inv *invocation.Invocation, cb invocation.ResponseCallBack) {
	au.config = config.GetConfig().Authorization
	if au.config == nil {
		openlogging.Info("authorization config not found")
		return
	}

	authType, err := getType()
	if err != nil {
		openlogging.Error("parse the authorization kind error: " + err.Error())
		return
	}

	err = authType.Authorize(inv, cb)
	if err != nil {
		openlogging.Error("authorization failed: " + err.Error())
		return
	}

	chain.Next(inv, func(r *invocation.Response) error {
		return cb(r)
	})
}

// Name returns auth name
func (au *AuthHandler) Name() string {
	return Authorizer
}

// NewAuth returns new auth handler
func NewAuth() handler.Handler {
	return &AuthHandler{}
}

func init() {
	err := handler.RegisterHandler(Authorizer, NewAuth)
	if err != nil {
		openlogging.Error("register handler error: " + err.Error())
		return
	}
}

// getType fetch the authorization type
func getType() (authority auth.Auth, err error) {
	conf := config.GetConfig().Authorization
	authority, err = auth.NewAuth(conf.Type)
	if err != nil {
		openlogging.Error("parse the authorization kind failed: " + err.Error())
	}
	return
}

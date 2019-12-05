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
	"errors"
	"fmt"
	"github.com/apache/servicecomb-mesher/proxy/config"
	"github.com/emicklei/go-restful"
	"github.com/go-chassis/go-chassis/core/invocation"
	"github.com/go-mesh/openlogging"
	"golang.org/x/oauth2"
	"net/http"
)

var authMap = map[string]Auth{}

// Auth provide Authorize and Authenticate interface
type Auth interface {
	Authorize(inv *invocation.Invocation, cb invocation.ResponseCallBack) error
	Authenticate(inv *invocation.Invocation, cb invocation.ResponseCallBack) error
}

// RegisterAuthPlugin support the authorizer Plugins register
func RegisterAuthPlugin(kind string, auth Auth) {
	_, ok := authMap[kind]
	if ok {
		openlogging.Info("authorizer is already exit, name = " + kind)
		return
	}
	authMap[kind] = auth
}

// NewAuth parse the authorization kind
func NewAuth(kind string) (Auth, error) {
	a, ok := authMap[kind]
	if !ok {
		return nil, fmt.Errorf("authorization kind is not found, name = %s", kind)
	}
	return a, nil
}

// AuthConfig struct contain the config info
type AuthConfig struct {
	config *config.MesherConfig
}

// NewConfig reads the config file content
func NewConfig() *AuthConfig {
	conf := config.GetConfig()

	return &AuthConfig{config: conf}
}

// AuthConfig deal with the config file info
func (a *AuthConfig) AuthConfig() *oauth2.Config {
	var conf = oauth2.Config{
		ClientID:     a.config.Authorization.Client.ClientID,
		ClientSecret: a.config.Authorization.Client.ClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  a.config.Authorization.Endpoint.AuthURL,
			TokenURL: a.config.Authorization.Endpoint.TokenURL,
		},
		RedirectURL: a.config.Authorization.RedirectURL,
		Scopes:      a.config.Authorization.Scopes,
	}
	return &conf
}

// InvResponse returns the http status
func InvResponse(statusCode int, err error) *invocation.Response {
	return &invocation.Response{
		Status: statusCode,
		Err:    err,
	}
}

// GetHTTPRequest return the request
func GetHTTPRequest(inv *invocation.Invocation) (req *http.Request, err error) {
	switch r := inv.Args.(type) {
	case *restful.Request:
		req = r.Request
	case *http.Request:
		req = r
	default:
		openlogging.Error("http request not found: " + err.Error())
		err = errors.New("bad request")
	}
	return
}

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
	"encoding/base64"
	"errors"
	"github.com/apache/servicecomb-mesher/proxy/plugins/auth"
	"github.com/go-chassis/go-chassis/core/invocation"

	"github.com/go-mesh/openlogging"
	"net/http"
	"strings"
)

// errors
var (
	ErrInvalidBase64  = errors.New("invalid base64")
	ErrInvalidAuth    = errors.New("invalid authentication")
	ErrInvalidRequest = errors.New("invalid http request")
	ErrNoHeader       = errors.New("not authorized")
)

// AuthName is the authorization style
const AuthName = "basicAuth"

// HeaderAuth is common auth header
const HeaderAuth = "Authorization"

// BasicAuth struct
type BasicAuth struct{}

// Authorize is the method of basic authorization
func (b BasicAuth) Authorize(inv *invocation.Invocation, cb invocation.ResponseCallBack) error {
	req, err := auth.GetHTTPRequest(inv)
	if err != nil {
		cb(auth.InvResponse(http.StatusBadRequest, err))
		return ErrInvalidRequest
	}

	headerValue := req.Header.Get(HeaderAuth)
	if headerValue == "" {
		cb(auth.InvResponse(http.StatusUnauthorized, err))
		return ErrNoHeader
	}

	u, p, err := decode(headerValue)
	if err != nil {
		openlogging.Error("can not decode base 64:" + err.Error())
		cb(auth.InvResponse(http.StatusUnauthorized, err))
		return err
	}

	err = userInfo.UserInfo(u, p)
	if err != nil {
		openlogging.Error("authorization failed:" + err.Error())
		cb(auth.InvResponse(http.StatusUnauthorized, err))
		return err
	}

	return nil
}

// Authenticate is the implement of the interface
func (b BasicAuth) Authenticate(inv *invocation.Invocation, cb invocation.ResponseCallBack) error {
	return nil
}

func init() {
	auth.RegisterAuthPlugin(AuthName, &BasicAuth{})
}

// decode process the header info
func decode(headerValue string) (username string, password string, err error) {
	parts := strings.Split(headerValue, " ")
	if len(parts) != 2 {
		return "", "", ErrInvalidAuth

	}
	if parts[0] != "Basic" {
		return "", "", ErrInvalidAuth
	}
	s, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return "", "", ErrInvalidBase64
	}

	userPwd := strings.Split(string(s), ":")
	if len(userPwd) != 2 {
		return "", "", ErrInvalidAuth
	}

	return userPwd[0], userPwd[1], nil
}

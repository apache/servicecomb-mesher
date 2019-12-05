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
	"github.com/apache/servicecomb-mesher/proxy/plugins/auth"
	"github.com/apache/servicecomb-mesher/proxy/plugins/auth/oauth2/oauth2manage"
	"github.com/go-chassis/go-chassis/core/invocation"
	"github.com/go-mesh/openlogging"
	"net/http"

	// authorization_code model
	_ "github.com/apache/servicecomb-mesher/proxy/plugins/auth/oauth2/oauth2manage/authorizationcode"
)

// AuthName is the authorization style
const AuthName = "oauth2"

// OAuth2 struct contain the config info
type OAuth2 struct{}

// Authorize is the method of oauth2 authorization
func (oa *OAuth2) Authorize(inv *invocation.Invocation, cb invocation.ResponseCallBack) error {
	req, err := auth.GetHTTPRequest(inv)
	if err != nil {
		cb(auth.InvResponse(http.StatusBadRequest, err))
		return err
	}

	grandType := req.FormValue("grand_type")
	if grandType == "" {
		return errors.New("can not fetch the grand_type")
	}

	gt, err := oauth2manage.NewType(grandType)
	if err != nil {
		openlogging.Error("grand_type error: " + err.Error())
		return err
	}
	err = gt.GrandTypeProcess(inv, cb)
	if err != nil {
		openlogging.Error("authorization error: " + err.Error())
		return err
	}

	return nil
}

func init() {
	auth.RegisterAuthPlugin(AuthName, &OAuth2{})
}

// Authenticate is the implement of the interface
func (oa *OAuth2) Authenticate(inv *invocation.Invocation, cb invocation.ResponseCallBack) error {
	return nil
}

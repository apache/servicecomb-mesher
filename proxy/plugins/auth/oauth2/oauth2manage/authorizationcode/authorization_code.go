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

package authorizationcode

import (
	"context"
	"errors"
	"github.com/apache/servicecomb-mesher/proxy/plugins/auth"
	"github.com/apache/servicecomb-mesher/proxy/plugins/auth/oauth2/oauth2manage"
	"github.com/go-chassis/go-chassis/core/invocation"
	"github.com/go-mesh/openlogging"
	"net/http"
	"time"
)

// GrandType is the authorization model
const GrandType = "authorization_code"

// AuthorizationCode struct
type AuthorizationCode struct{}

func init() {
	oauth2manage.RegisterType(GrandType, &AuthorizationCode{})
}

// GrandTypeProcess is the method of authorization model
func (a *AuthorizationCode) GrandTypeProcess(inv *invocation.Invocation, cb invocation.ResponseCallBack) error {
	req, err := auth.GetHTTPRequest(inv)
	if err != nil {
		_ = cb(auth.InvResponse(http.StatusBadRequest, err))
		return err
	}
	code := req.FormValue("code")
	if code == "" {
		return errors.New("get code failed")
	}
	conf := auth.NewConfig()
	token, err := conf.AuthConfig().Exchange(context.Background(), code)
	if err != nil {
		openlogging.Error("get token failed, errors: " + err.Error())
		return err
	}
	if time.Now().After(token.Expiry) {
		return errors.New("expired token")
	}

	return nil
}

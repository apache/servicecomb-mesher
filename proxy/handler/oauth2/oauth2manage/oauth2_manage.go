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

package oauth2manage

import (
	"fmt"
	"github.com/go-chassis/go-chassis/core/invocation"
	"github.com/go-mesh/openlogging"
)

// grantTypeMap saves the authorization model
var grantTypeMap = map[string]Subject{}

// Subject interface process the grant_type
type Subject interface {
	GrantTypeProcess(inv *invocation.Invocation, cb invocation.ResponseCallBack) (string, error)
}

// RegisterType is handle the grant_type
func RegisterType(kind string, subject Subject) {
	_, ok := grantTypeMap[kind]
	if ok {
		openlogging.Info("grant type is already exit, name = " + kind)
		return
	}
	grantTypeMap[kind] = subject
}

// NewType is new a grant_type
func NewType(kind string) (Subject, error) {
	a, ok := grantTypeMap[kind]
	if !ok {
		return nil, fmt.Errorf("grant type is not found, name = %s", kind)
	}
	return a, nil
}

// InvResponse returns the http status
func InvResponse(statusCode int, err error) *invocation.Response {
	return &invocation.Response{
		Status: statusCode,
		Err:    err,
	}
}

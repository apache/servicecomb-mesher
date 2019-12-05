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

package noop

import (
	"github.com/apache/servicecomb-mesher/proxy/plugins/auth"
	"github.com/go-chassis/go-chassis/core/invocation"
)

// AuthName is the default name
const AuthName = "noop"

// Noop is an empty struct
type Noop struct {
}

// Authorize is an empty implementation
func (n *Noop) Authorize(inv *invocation.Invocation, cb invocation.ResponseCallBack) error {
	return nil
}

// Authenticate is an empty implementation
func (n *Noop) Authenticate(inv *invocation.Invocation, cb invocation.ResponseCallBack) error {
	return nil
}

func init() {
	auth.RegisterAuthPlugin(AuthName, &Noop{})
}

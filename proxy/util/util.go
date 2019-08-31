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

package util

import (
	"fmt"
	"github.com/apache/servicecomb-mesher/proxy/cmd"
	"github.com/apache/servicecomb-mesher/proxy/common"
	"github.com/go-chassis/go-chassis/core/invocation"
)

//SetLocalServiceAddress assign invocation endpoint a local service address
//header "X-Forwarded-Port" has highest priority
// if it is empty
// it uses ports config in cmd param or env
func SetLocalServiceAddress(inv *invocation.Invocation, port string) error {
	inv.Endpoint = cmd.Configs.PortsMap[inv.Protocol]
	if port == "" {
		inv.Endpoint = cmd.Configs.PortsMap[inv.Protocol]
		if inv.Endpoint == "" {
			return fmt.Errorf("[%s] is not supported, [%s] didn't set env [%s] or cmd parameter --service-ports before mesher start",
				inv.Protocol, inv.MicroServiceName, common.EnvServicePorts)
		}
		return nil
	}
	inv.Endpoint = "127.0.0.1:" + port
	return nil
}

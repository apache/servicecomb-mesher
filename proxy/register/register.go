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

package register

import (
	"github.com/apache/servicecomb-mesher/proxy/common"
	chassisCommon "github.com/go-chassis/go-chassis/core/common"
	"github.com/go-chassis/go-chassis/core/config"
	chassisModel "github.com/go-chassis/go-chassis/core/config/model"
	"github.com/go-chassis/go-chassis/core/registry"
	"github.com/go-chassis/go-chassis/pkg/util/iputil"
	"github.com/go-mesh/openlogging"
	"strings"
)

// AdaptEndpoints moves http endpoint to rest endpoint
func AdaptEndpoints() error {
	var err error
	// To be called by services based on CSE SDK,
	// mesher has to register endpoint with rest://ip:port
	oldProtoMap := config.GlobalDefinition.Cse.Protocols
	if _, ok := oldProtoMap[common.HTTPProtocol]; !ok {
		return nil
	}
	if _, ok := oldProtoMap[chassisCommon.ProtocolRest]; ok {
		return nil
	}

	newProtoMap := make(map[string]chassisModel.Protocol)
	for n, proto := range oldProtoMap {
		if n == common.HTTPProtocol {
			continue
		}
		newProtoMap[n] = proto
	}
	newProtoMap[chassisCommon.ProtocolRest] = oldProtoMap[common.HTTPProtocol]
	eps, err := registry.MakeEndpointMap(newProtoMap)
	if err != nil {
		return err
	}
	for protocol, ep := range eps {
		if ep.Address == "" {
			port := strings.Split(newProtoMap[protocol].Listen, ":")
			if len(port) == 2 { //check if port is not specified along with ip address, eventually in case port is not specified, server start will fail in subsequent processing.
				registry.InstanceEndpoints[protocol] = iputil.GetLocalIP() + ":" + port[1]
			}
		} else {
			registry.InstanceEndpoints[protocol] = ep.Address
		}
	}

	openlogging.Debug("adapt endpoints success")
	return nil
}

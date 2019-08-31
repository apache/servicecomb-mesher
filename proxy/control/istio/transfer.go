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

package istio

import (
	"github.com/apache/servicecomb-mesher/proxy/config"
	"github.com/go-chassis/go-chassis/control"
)

//SaveToEgressCache save the egress rules in the cache
func SaveToEgressCache(egressConfigFromPilot map[string][]*config.EgressRule) {
	{
		var egressconfig []control.EgressConfig
		for _, v := range egressConfigFromPilot {
			for _, v1 := range v {
				var Ports []*control.EgressPort
				for _, v2 := range v1.Ports {
					p := control.EgressPort{
						Port:     (*v2).Port,
						Protocol: (*v2).Protocol,
					}
					Ports = append(Ports, &p)
				}
				c := control.EgressConfig{
					Hosts: v1.Hosts,
					Ports: Ports,
				}

				egressconfig = append(egressconfig, c)
			}
		}
		EgressConfigCache.Set("", egressconfig, 0)
	}
}

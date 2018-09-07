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
	"github.com/go-chassis/go-chassis/core/invocation"
	"github.com/go-mesh/mesher/common"
	"github.com/go-mesh/mesher/config"
)

//EqualPolicy is a function
func EqualPolicy(inv *invocation.Invocation, p *config.Policy) bool {
	if inv.MicroServiceName != p.Destination {
		return false
	}
	for k, v := range p.Tags {
		if k == common.BuildInTagApp {
			if v == "" {
				v = common.DefaultApp
			}
			if v != inv.RouteTags.AppID() {
				return false
			}
			continue
		}
		if k == common.BuildInTagVersion {
			if v == "" {
				v = common.DefaultVersion
			}
			if v != inv.RouteTags.Version() {
				return false
			}
			continue
		}
		t, ok := inv.Metadata[k]
		if !ok {
			return false
		}
		if _, ok := t.(string); !ok {
			return false
		}
	}
	for k, v := range inv.Metadata {
		if v != p.Tags[k] {
			return false
		}
	}
	return true

}

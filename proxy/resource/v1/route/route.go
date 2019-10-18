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

package route

import (
	"github.com/go-chassis/go-chassis/core/config"
	"github.com/go-chassis/go-chassis/core/router"
)

//Rules is the struct for route rule
type Rules struct {
	Destinations map[string][]*config.RouteRule `yaml:"routeRule"`
}

var routeRules *Rules

//GetServiceRouteRule gets route rule for that service
func GetServiceRouteRule(serviceName string) []*config.RouteRule {
	if routeRules != nil {
		if v, ok := routeRules.Destinations[serviceName]; ok {
			return v
		}
	}

	rules := router.DefaultRouter.FetchRouteRuleByServiceName(serviceName)
	if routeRules == nil {
		routeRules = new(Rules)
	}
	if routeRules.Destinations == nil {
		routeRules.Destinations = make(map[string][]*config.RouteRule)
	}
	routeRules.Destinations[serviceName] = rules
	return rules
}

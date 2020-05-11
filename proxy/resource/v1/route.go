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

package v1

import (
	"fmt"
	"github.com/go-mesh/openlogging"
	"net/http"

	"github.com/go-chassis/go-chassis/core/common"
	"github.com/go-chassis/go-chassis/core/router"
	"github.com/go-chassis/go-chassis/server/restful"
)

// RouteResource is rest api to manage route rule
type RouteResource struct{}

//RouteRuleByService returns route config for particular service
func (a *RouteResource) RouteRuleByService(context *restful.Context) {
	serviceName := context.ReadPathParameter("serviceName")
	routeRule := router.DefaultRouter.FetchRouteRuleByServiceName(serviceName)
	if routeRule == nil {
		err := context.WriteHeaderAndJSON(http.StatusNotFound, fmt.Sprintf("%s routeRule not found", serviceName), common.JSON)
		if err != nil {
			openlogging.GetLogger().Errorf("Write HeaderAndJSON error %s: ", err.Error())
		}
		return
	}
	err := context.WriteHeaderAndJSON(http.StatusOK, routeRule, "text/vnd.yaml")
	if err != nil {
		openlogging.GetLogger().Errorf("Write HeaderAndJSON error %s: ", err.Error())
	}
}

//URLPatterns helps to respond for  Admin API calls
func (a *RouteResource) URLPatterns() []restful.Route {
	return []restful.Route{
		{Method: http.MethodGet, Path: "/v1/mesher/routeRule/{serviceName}", ResourceFuncName: "RouteRuleByService"},
	}
}

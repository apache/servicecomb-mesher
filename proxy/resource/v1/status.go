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
	"github.com/apache/servicecomb-mesher/proxy/resource/v1/health"
	"github.com/apache/servicecomb-mesher/proxy/resource/v1/version"
	"github.com/go-chassis/go-chassis/core/common"
	"github.com/go-chassis/go-chassis/pkg/metrics"
	"github.com/go-chassis/go-chassis/server/restful"
	"github.com/go-mesh/openlogging"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

//StatusResource is rest api
type StatusResource struct{}

//Health returns mesher health
func (a *StatusResource) Health(context *restful.Context) {
	healthResp := health.GetMesherHealth()
	if healthResp.Status == health.Red {
		err := context.WriteHeaderAndJSON(http.StatusInternalServerError, healthResp, common.JSON)
		if err != nil {
			openlogging.GetLogger().Errorf("Write HeaderAndJSON error %s: ", err.Error())
		}
		return
	}
	err := context.WriteHeaderAndJSON(http.StatusOK, healthResp, common.JSON)
	if err != nil {
		openlogging.GetLogger().Errorf("Write HeaderAndJSON error %s: ", err.Error())
	}
}

//GetMetrics returns metrics data
func (a *StatusResource) GetMetrics(context *restful.Context) {
	resp := context.ReadResponseWriter()
	req := context.ReadRequest()
	promhttp.HandlerFor(metrics.GetSystemPrometheusRegistry(), promhttp.HandlerOpts{}).ServeHTTP(resp, req)
}

//GetVersion writes version in response header
func (a *StatusResource) GetVersion(context *restful.Context) {
	versions := version.Ver()
	err := context.WriteHeaderAndJSON(http.StatusOK, versions, common.JSON)
	if err != nil {
		openlogging.GetLogger().Errorf("Write HeaderAndJSON error %s: ", err.Error())
	}
}

//URLPatterns helps to respond for  Admin API calls
func (a *StatusResource) URLPatterns() []restful.Route {
	return []restful.Route{
		{Method: http.MethodGet, Path: "/v1/mesher/version", ResourceFuncName: "GetVersion"},
		{Method: http.MethodGet, Path: "/v1/mesher/metrics", ResourceFuncName: "GetMetrics"},
		{Method: http.MethodGet, Path: "/v1/mesher/health", ResourceFuncName: "Health"},
	}
}

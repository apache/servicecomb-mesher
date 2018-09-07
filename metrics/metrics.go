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

package metrics

import (
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/go-chassis/go-chassis/core/config"
	"github.com/go-chassis/go-chassis/core/invocation"
	mesherconfig "github.com/go-mesh/mesher/config"
	"github.com/prometheus/client_golang/prometheus"
)

//Constants with attributes for metrics data
const (
	TotalRequest          = "requests_total"
	TotalSuccess          = "successes_total"
	TotalFailures         = "failures_total"
	RequestLatencySeconds = "request_latency_seconds"
	Error4XX              = "status_4xx"
	Error5XX              = "status_5xx"
	ServiceName           = "servicename"
	AppID                 = "appid"
	Version               = "version"
)

var (
	//LabelNames is a list with servicename, appID, version
	LabelNames = []string{ServiceName, AppID, Version}
	mutex      = sync.Mutex{}
)

var onceEnable sync.Once

//Init function initiates all config
func Init() {
	mesherLabelValues := map[string]string{ServiceName: config.SelfServiceName, AppID: config.GlobalDefinition.AppID, Version: config.SelfVersion}
	mesherStartTime := time.Now().Unix()
	DefaultPrometheusExporter.Gauge("start_time_seconds", float64(mesherStartTime), LabelNames, mesherLabelValues)
	mesherConfig := mesherconfig.GetConfig()
	promConfig := getPrometheusSinker(getSystemPrometheusRegistry())
	if mesherConfig.Admin.GoRuntimeMetrics == true {
		onceEnable.Do(func() {
			promConfig.PromRegistry.MustRegister(prometheus.NewProcessCollector(os.Getpid(), ""))
			promConfig.PromRegistry.MustRegister(prometheus.NewGoCollector())
		})
	}
}

//RecordResponse record the response
func RecordResponse(inv *invocation.Invocation, statusCode int) {
	mutex.Lock()
	defer mutex.Unlock()
	serviceLabelValues := map[string]string{ServiceName: inv.MicroServiceName, AppID: inv.RouteTags.AppID(), Version: inv.RouteTags.Version()}
	if statusCode >= http.StatusBadRequest && statusCode <= http.StatusUnavailableForLegalReasons {
		DefaultPrometheusExporter.Count(Error4XX, LabelNames, serviceLabelValues)
		DefaultPrometheusExporter.Count(TotalFailures, LabelNames, serviceLabelValues)
	} else if statusCode >= http.StatusInternalServerError && statusCode <= http.StatusNetworkAuthenticationRequired {
		DefaultPrometheusExporter.Count(Error5XX, LabelNames, serviceLabelValues)
		DefaultPrometheusExporter.Count(TotalFailures, LabelNames, serviceLabelValues)
	} else if statusCode >= http.StatusOK && statusCode <= http.StatusIMUsed {
		DefaultPrometheusExporter.Count(TotalSuccess, LabelNames, serviceLabelValues)
	}

	DefaultPrometheusExporter.Count(TotalRequest, LabelNames, serviceLabelValues)
}

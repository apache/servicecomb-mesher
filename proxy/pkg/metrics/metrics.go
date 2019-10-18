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

//Package metrics is a system-independent module
//it consider metrics key as first-class citizen
//each function is for recording one kind of metrics key and value
//it expose standard API to record runtime metrics for a service
//use prom as default metrics system
package metrics

import (
	"github.com/go-chassis/go-chassis/pkg/runtime"
	"sync"
	"time"
)

//Constants with attributes for metrics data
//Label start with word "L"
const (
	LTotalRequest          = "requests_total"
	LTotalSuccess          = "successes_total"
	LTotalFailures         = "failures_total"
	LRequestLatencySeconds = "request_latency_seconds"
	LError4XX              = "status_4xx"
	LError5XX              = "status_5xx"
	LServiceName           = "service_name"
	LApp                   = "app"
	LVersion               = "version"
	LStartTime             = "start_time_seconds"
)

var (
	//LabelNames is a fixed list with service name, appID, version
	LabelNames = []string{LServiceName, LApp, LVersion}
	mutex      = sync.Mutex{}
)

//Options define recorder options
type Options struct {
	LabelNames []string //default label names, if RecordOptions LabelNames is nil
}

var defaultRecorder *PromRecorder

//RecordStatus record an operation status
func RecordStatus(labelValues map[string]string, statusCode int) {
	defaultRecorder.RecordStatus(labelValues, statusCode)
}

//RecordLatency record an operation latency
func RecordLatency(labelValues map[string]string, latency float64) {
	defaultRecorder.RecordLatency(labelValues, latency)
}

//RecordStartTime record mesher start time
func RecordStartTime(labelValues map[string]string, start time.Time) {
	defaultRecorder.RecordStartTime(labelValues, start)
}

//Init initiate the recorder
func Init() error {
	var err error
	LabelValues := map[string]string{LServiceName: runtime.ServiceName, LApp: runtime.App, LVersion: runtime.Version}
	defaultRecorder, err = NewPromRecorder(&Options{
		LabelNames: LabelNames,
	})
	if err != nil {
		return err
	}
	RecordStartTime(LabelValues, time.Now())
	return nil
}

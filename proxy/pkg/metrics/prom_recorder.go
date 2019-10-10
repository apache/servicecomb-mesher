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
	"errors"
	"net/http"
	"sync"
	"time"
)

var onceEnable sync.Once

//PromRecorder record metrics
type PromRecorder struct {
	LabelNames []string
}

//NewPromRecorder return a prom recorder
func NewPromRecorder(opts *Options) (*PromRecorder, error) {
	if opts != nil {
		return &PromRecorder{LabelNames: opts.LabelNames}, nil
	}
	return nil, errors.New("options can not be nil")
}

//RecordStatus record different metrics based on status
func (e *PromRecorder) RecordStatus(LabelValues map[string]string, statusCode int) {
	if statusCode >= http.StatusBadRequest && statusCode <= http.StatusUnavailableForLegalReasons {
		DefaultPrometheusExporter.Count(LError4XX, e.LabelNames, LabelValues)
		DefaultPrometheusExporter.Count(LTotalFailures, e.LabelNames, LabelValues)
	} else if statusCode >= http.StatusInternalServerError && statusCode <= http.StatusNetworkAuthenticationRequired {
		DefaultPrometheusExporter.Count(LError5XX, e.LabelNames, LabelValues)
		DefaultPrometheusExporter.Count(LTotalFailures, e.LabelNames, LabelValues)
	} else if statusCode >= http.StatusOK && statusCode <= http.StatusIMUsed {
		DefaultPrometheusExporter.Count(LTotalSuccess, e.LabelNames, LabelValues)
	}
	DefaultPrometheusExporter.Count(LTotalRequest, e.LabelNames, LabelValues)
}

//RecordLatency record operation latency
func (e *PromRecorder) RecordLatency(LabelValues map[string]string, latency float64) {
	DefaultPrometheusExporter.Summary(LRequestLatencySeconds, latency, e.LabelNames, LabelValues)

}

//RecordStartTime save start time
func (e *PromRecorder) RecordStartTime(LabelValues map[string]string, start time.Time) {
	DefaultPrometheusExporter.Gauge(LStartTime, float64(start.Unix()), e.LabelNames, LabelValues)

}

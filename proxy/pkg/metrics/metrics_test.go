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

package metrics_test

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"net/http"
	"strings"
	"testing"

	"github.com/apache/servicecomb-mesher/proxy/pkg/metrics"
	"github.com/go-chassis/go-chassis/pkg/runtime"
)

func TestInit(t *testing.T) {
	err := metrics.Init()
	runtime.ServiceName = "A"
	runtime.Version = "v1.1"
	runtime.App = "app"
	assert.NoError(t, err)
}
func TestRecordStatus(t *testing.T) {
	assert := assert.New(t)
	var errorcount4xx float64
	var errorcount5xx float64
	lvs := map[string]string{
		metrics.LServiceName: "service",
		metrics.LVersion:     "",
		metrics.LApp:         "",
	}
	metrics.RecordStatus(lvs, http.StatusOK)
	metrics.RecordStatus(lvs, http.StatusNotFound)
	metrics.RecordStatus(lvs, http.StatusInternalServerError)
	metricFamilies, err := prometheus.DefaultGatherer.Gather()
	assert.Nil(err, "error should be nil while collecting metrics from prometheus")
	for _, metricFamily := range metricFamilies {
		if name := metricFamily.GetName(); strings.Contains(name, metrics.LError5XX) {
			errorcount4xx += *metricFamily.Metric[0].Counter.Value
		}
	}
	for _, metricFamily := range metricFamilies {
		if name := metricFamily.GetName(); strings.Contains(name, metrics.LError5XX) {
			errorcount5xx += *metricFamily.Metric[0].Counter.Value
		}
	}
	assert.Equal(errorcount4xx, float64(1))
	assert.Equal(errorcount5xx, float64(1))

	metrics.RecordLatency(lvs, 1000)
	metricFamilies, err = prometheus.DefaultGatherer.Gather()
	var c float64
	for _, metricFamily := range metricFamilies {
		if name := metricFamily.GetName(); strings.Contains(name, metrics.LRequestLatencySeconds) {
			c = *metricFamily.Metric[0].Summary.SampleSum
		}
	}
	assert.Equal(1000, int(c))
}

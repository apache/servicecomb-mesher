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
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/assert"
)

var (
	labelNames  = []string{"APPID", "VERSION"}
	labelValues = map[string]string{"APPID": "sockshop", "VERSION": "0.1"}
)

func TestPrometheusConfig_CounterFromNameAndLabelValues(t *testing.T) {
	assert := assert.New(t)
	var totalMetricCreated int
	DefaultPrometheusExporter.Count("total_request", labelNames, labelValues)
	metricFamilies, err := prometheus.DefaultGatherer.Gather()
	assert.Nil(err, "error should be nil while collecting metrics from prometheus")
	for _, metricFamily := range metricFamilies {
		if metricName := metricFamily.GetName(); strings.Contains(metricName, "total_request") {
			assert.Equal(metricFamily.GetType(), dto.MetricType_COUNTER)
			totalMetricCreated++
		}
	}
	assert.Equal(totalMetricCreated, 1)
}

func TestPrometheusConfig_GaugeFromNameAndLabelValues(t *testing.T) {
	assert := assert.New(t)
	var totalMetricCreated int
	var gaugeValue *float64
	DefaultPrometheusExporter.Gauge("memory_used", 12, labelNames, labelValues)
	metricFamilies, err := prometheus.DefaultGatherer.Gather()
	assert.Nil(err, "error should be nil while collecting metrics from prometheus")
	for _, metricFamily := range metricFamilies {
		if metricName := metricFamily.GetName(); strings.Contains(metricName, "memory_used") {
			assert.Equal(metricFamily.GetType(), dto.MetricType_GAUGE)
			totalMetricCreated++
			gaugeValue = metricFamily.Metric[0].Gauge.Value
		}
	}
	assert.Equal(totalMetricCreated, 1)
	assert.Equal(*gaugeValue, float64(12))
}

func TestPrometheusConfig_SummaryFromNameAndLabelValues(t *testing.T) {
	assert := assert.New(t)
	var totalMetricCreated int
	var sampleCount *uint64
	DefaultPrometheusExporter.Summary("request_latency", 12, labelNames, labelValues)
	metricFamilies, err := prometheus.DefaultGatherer.Gather()
	assert.Nil(err, "error should be nil while collecting metrics from prometheus")
	for _, metricFamily := range metricFamilies {
		if metricName := metricFamily.GetName(); strings.Contains(metricName, "request_latency") {
			assert.Equal(metricFamily.GetType(), dto.MetricType_SUMMARY)
			totalMetricCreated++
			sampleCount = metricFamily.Metric[0].Summary.SampleCount
		}
	}
	assert.Equal(totalMetricCreated, 1)
	assert.Equal(*sampleCount, uint64(1))
}

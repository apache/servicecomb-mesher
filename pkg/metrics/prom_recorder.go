package metrics

import (
	"errors"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"os"
	"sync"
	"time"
)

var onceEnable sync.Once

type promRecorder struct {
	LabelNames []string
}

//NewPromRecorder return a prom recorder
func NewPromRecorder(opts *Options) (Recorder, error) {
	promConfig := getPrometheusSinker(getSystemPrometheusRegistry())
	if opts != nil {
		if opts.EnableGoRuntimeMetrics {
			onceEnable.Do(func() {
				promConfig.PromRegistry.MustRegister(prometheus.NewProcessCollector(os.Getpid(), ""))
				promConfig.PromRegistry.MustRegister(prometheus.NewGoCollector())
			})
		}
		return &promRecorder{LabelNames: opts.LabelNames}, nil
	}
	return nil, errors.New("options can not be nil")
}

//GetLN return label names based on options
func (e *promRecorder) GetLN(opts *RecordOptions) (ln []string) {
	ln = e.LabelNames
	if opts != nil && len(opts.LabelNames) != 0 {
		ln = opts.LabelNames
	}
	return
}

//RecordStatus record different metrics based on status
func (e *promRecorder) RecordStatus(LabelValues map[string]string, statusCode int, opts *RecordOptions) {
	ln := e.GetLN(opts)
	if statusCode >= http.StatusBadRequest && statusCode <= http.StatusUnavailableForLegalReasons {
		DefaultPrometheusExporter.Count(LError4XX, ln, LabelValues)
		DefaultPrometheusExporter.Count(LTotalFailures, ln, LabelValues)
	} else if statusCode >= http.StatusInternalServerError && statusCode <= http.StatusNetworkAuthenticationRequired {
		DefaultPrometheusExporter.Count(LError5XX, ln, LabelValues)
		DefaultPrometheusExporter.Count(LTotalFailures, ln, LabelValues)
	} else if statusCode >= http.StatusOK && statusCode <= http.StatusIMUsed {
		DefaultPrometheusExporter.Count(LTotalSuccess, ln, LabelValues)
	}
	DefaultPrometheusExporter.Count(LTotalRequest, ln, LabelValues)
}

//RecordLatency record operation latency
func (e *promRecorder) RecordLatency(LabelValues map[string]string, latency float64, opts *RecordOptions) {
	ln := e.GetLN(opts)
	DefaultPrometheusExporter.Summary(LRequestLatencySeconds, latency, ln, LabelValues)

}

//RecordStartTime save start time
func (e *promRecorder) RecordStartTime(LabelValues map[string]string, start time.Time, opts *RecordOptions) {
	ln := e.GetLN(opts)
	DefaultPrometheusExporter.Gauge(LStartTime, float64(start.Unix()), ln, LabelValues)

}

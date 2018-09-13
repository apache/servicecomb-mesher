package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"os"
	"sync"
	"time"
)

var onceEnable sync.Once

type promRecorder struct {
}

func newPromRecorder(options *Options) Recorder {
	promConfig := getPrometheusSinker(getSystemPrometheusRegistry())
	if options.EnableGoRuntimeMetrics {
		onceEnable.Do(func() {
			promConfig.PromRegistry.MustRegister(prometheus.NewProcessCollector(os.Getpid(), ""))
			promConfig.PromRegistry.MustRegister(prometheus.NewGoCollector())
		})
	}
	return &promRecorder{}
}

//RecordStatus record different metrics based on status
func (e *promRecorder) RecordStatus(
	LabelNames []string, LabelValues map[string]string, statusCode int) {
	if statusCode >= http.StatusBadRequest && statusCode <= http.StatusUnavailableForLegalReasons {
		DefaultPrometheusExporter.Count(LError4XX, LabelNames, LabelValues)
		DefaultPrometheusExporter.Count(LTotalFailures, LabelNames, LabelValues)
	} else if statusCode >= http.StatusInternalServerError && statusCode <= http.StatusNetworkAuthenticationRequired {
		DefaultPrometheusExporter.Count(LError5XX, LabelNames, LabelValues)
		DefaultPrometheusExporter.Count(LTotalFailures, LabelNames, LabelValues)
	} else if statusCode >= http.StatusOK && statusCode <= http.StatusIMUsed {
		DefaultPrometheusExporter.Count(LTotalSuccess, LabelNames, LabelValues)
	}
	DefaultPrometheusExporter.Count(LTotalRequest, LabelNames, LabelValues)
}

//RecordLatency TODO
func (e *promRecorder) RecordLatency(
	LabelNames []string, LabelValues map[string]string, latency float64) {

}

//RecordStartTime save start time
func (e *promRecorder) RecordStartTime(LabelNames []string, LabelValues map[string]string, start time.Time) {
	DefaultPrometheusExporter.Gauge(LStartTime, float64(start.Unix()), LabelNames, LabelValues)

}

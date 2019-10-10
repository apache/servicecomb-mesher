package http

import (
	"github.com/apache/servicecomb-mesher/proxy/bootstrap"
	"github.com/apache/servicecomb-mesher/proxy/pkg/metrics"
	"github.com/go-chassis/go-chassis"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleIngressTraffic(t *testing.T) {
	err := metrics.Init()
	bootstrap.SetHandlers()
	chassis.Init()
	assert.NoError(t, err)
	req, _ := http.NewRequest(http.MethodGet, "/api", nil)
	w := httptest.NewRecorder()
	HandleIngressTraffic(w, req)
}

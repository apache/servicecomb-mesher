package http

import (
	"github.com/apache/servicecomb-mesher/proxy/ingress"
	_ "github.com/apache/servicecomb-mesher/proxy/ingress/servicecomb"
	"github.com/apache/servicecomb-mesher/proxy/pkg/metrics"
	"github.com/go-chassis/go-archaius"
	"github.com/go-chassis/go-chassis/control"
	_ "github.com/go-chassis/go-chassis/control/servicecomb"
	"github.com/go-chassis/go-chassis/core/common"
	"github.com/go-chassis/go-chassis/core/config"
	"github.com/go-chassis/go-chassis/core/config/model"
	"github.com/go-chassis/go-chassis/core/handler"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandleIngressTraffic(t *testing.T) {
	b := []byte(`
        - host: example.com
          limit: 30
          apiPath: /some/api
          service:
            name: example
            tags:
              version: 1.0.0
            redirectPath: /another/api
            port:
              name: http-legacy
              value: 8080
        - host: foo.com
          apiPath: /some/api
          service:
            name: foo
            tags:
              version: 1.0.0
            redirectPath: /another/api
            port:
              name: http
              value: 8080
`)
	err := metrics.Init()
	assert.NoError(t, err)
	err = handler.CreateChains(common.Provider, map[string]string{
		"incoming": strings.Join([]string{}, ","),
	},
	)
	err = handler.CreateChains(common.Consumer, map[string]string{
		"outgoing": strings.Join([]string{}, ","),
	},
	)
	archaius.Init(archaius.WithMemorySource())
	err = archaius.Set("mesher.ingress.rule.http", string(b))
	assert.NoError(t, err)
	config.GlobalDefinition = new(model.GlobalCfg)
	config.HystrixConfig = &model.HystrixConfigWrapper{}
	archaius.UnmarshalConfig(config.GlobalDefinition)
	assert.NoError(t, err)
	err = control.Init(control.Options{})
	assert.NoError(t, err)
	err = ingress.Init()
	assert.NoError(t, err)

	req, _ := http.NewRequest(http.MethodGet, "http://foo.com/some/api", nil)
	w := httptest.NewRecorder()
	HandleIngressTraffic(w, req)
}

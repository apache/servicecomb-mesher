package adminapi

import (
	"fmt"
	"net/http"

	"github.com/go-chassis/go-chassis/core/common"
	"github.com/go-chassis/go-chassis/core/config/model"
	"github.com/go-chassis/go-chassis/core/router"
	"github.com/go-chassis/go-chassis/pkg/metrics"
	"github.com/go-chassis/go-chassis/server/restful"
	"github.com/go-mesh/mesher/adminapi/health"
	"github.com/go-mesh/mesher/adminapi/version"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Admin is a struct used for implementation of rest admin program
type Admin struct{}

//GetVersion writes version in response header
func (a *Admin) GetVersion(context *restful.Context) {
	versions := version.Ver()
	context.WriteHeaderAndJSON(http.StatusOK, versions, common.JSON)
}

//GetMetrics returns metrics data
func (a *Admin) GetMetrics(context *restful.Context) {
	resp := context.ReadResponseWriter()
	req := context.ReadRequest()
	promhttp.HandlerFor(metrics.GetSystemPrometheusRegistry(), promhttp.HandlerOpts{}).ServeHTTP(resp, req)
}

//RouteRule returns all router configs
func (a *Admin) RouteRule(context *restful.Context) {
	routerConfig := &model.RouterConfig{
		Destinations: router.DefaultRouter.FetchRouteRule(),
	}
	context.WriteHeaderAndJSON(http.StatusOK, routerConfig, "text/vnd.yaml")
}

//RouteRuleByService returns route config for particular service
func (a *Admin) RouteRuleByService(context *restful.Context) {

	serviceName := context.ReadPathParameter("serviceName")
	routeRule := router.DefaultRouter.FetchRouteRuleByServiceName(serviceName)
	if routeRule == nil {
		context.WriteHeaderAndJSON(http.StatusNotFound, fmt.Sprintf("%s routeRule not found", serviceName), common.JSON)
		return
	}
	context.WriteHeaderAndJSON(http.StatusOK, routeRule, "text/vnd.yaml")
}

//MesherHealth returns mesher health
func (a *Admin) MesherHealth(context *restful.Context) {
	healthResp := health.GetMesherHealth()
	if healthResp.Status == health.Red {
		context.WriteHeaderAndJSON(http.StatusInternalServerError, healthResp, common.JSON)
		return
	}
	context.WriteHeaderAndJSON(http.StatusOK, healthResp, common.JSON)
}

//URLPatterns helps to respond for  Admin API calls
func (a *Admin) URLPatterns() []restful.Route {
	return []restful.Route{
		{Method: http.MethodGet, Path: "/v1/mesher/version", ResourceFuncName: "GetVersion"},
		{Method: http.MethodGet, Path: "/v1/mesher/metrics", ResourceFuncName: "GetMetrics"},
		{Method: http.MethodGet, Path: "/v1/mesher/routeRule", ResourceFuncName: "RouteRule"},
		{Method: http.MethodGet, Path: "/v1/mesher/routeRule/{serviceName}", ResourceFuncName: "RouteRuleByService"},
		{Method: http.MethodGet, Path: "/v1/mesher/health", ResourceFuncName: "MesherHealth"},
	}
}

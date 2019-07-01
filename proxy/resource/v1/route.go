package v1

import (
	"fmt"
	"net/http"

	"github.com/go-chassis/go-chassis/core/common"
	"github.com/go-chassis/go-chassis/core/router"
	"github.com/go-chassis/go-chassis/server/restful"
)

// RouteResource is rest api to manage route rule
type RouteResource struct{}

//RouteRuleByService returns route config for particular service
func (a *RouteResource) RouteRuleByService(context *restful.Context) {
	serviceName := context.ReadPathParameter("serviceName")
	routeRule := router.DefaultRouter.FetchRouteRuleByServiceName(serviceName)
	if routeRule == nil {
		context.WriteHeaderAndJSON(http.StatusNotFound, fmt.Sprintf("%s routeRule not found", serviceName), common.JSON)
		return
	}
	context.WriteHeaderAndJSON(http.StatusOK, routeRule, "text/vnd.yaml")
}

//URLPatterns helps to respond for  Admin API calls
func (a *RouteResource) URLPatterns() []restful.Route {
	return []restful.Route{
		{Method: http.MethodGet, Path: "/v1/mesher/routeRule/{serviceName}", ResourceFuncName: "RouteRuleByService"},
	}
}

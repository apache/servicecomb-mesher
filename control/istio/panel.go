package istio

import (
	"github.com/go-chassis/go-chassis/control"
	"github.com/go-chassis/go-chassis/core/config/model"
	"github.com/go-chassis/go-chassis/core/invocation"
	"github.com/go-chassis/go-chassis/third_party/forked/afex/hystrix-go/hystrix"
	meshercontrol "github.com/go-mesh/mesher/control"
	"github.com/go-mesh/mesher/egress"
)

func init() {
	meshercontrol.InstallPlugin("egresspilot", newPilotPanel)
}

//PilotPanel pull configs from istio pilot
type PilotPanel struct {
}

func newPilotPanel(options meshercontrol.Options) control.Panel {
	SaveToEgressCache(egress.DefaultEgress.FetchEgressRule())
	return &PilotPanel{}
}

//GetEgressRule get egress config
func (p *PilotPanel) GetEgressRule() []control.EgressConfig {
	c, ok := EgressConfigCache.Get("")
	if !ok {

		return nil
	}
	return c.([]control.EgressConfig)
}

//GetCircuitBreaker return command , and circuit breaker settings
func (p *PilotPanel) GetCircuitBreaker(inv invocation.Invocation, serviceType string) (string, hystrix.CommandConfig) {
	return "", hystrix.CommandConfig{}

}

//GetLoadBalancing get load balancing config
func (p *PilotPanel) GetLoadBalancing(inv invocation.Invocation) control.LoadBalancingConfig {
	return control.LoadBalancingConfig{}

}

//GetRateLimiting get rate limiting config
func (p *PilotPanel) GetRateLimiting(inv invocation.Invocation, serviceType string) control.RateLimitingConfig {
	return control.RateLimitingConfig{}
}

//GetFaultInjection get Fault injection config
func (p *PilotPanel) GetFaultInjection(inv invocation.Invocation) model.Fault {
	return model.Fault{}
}

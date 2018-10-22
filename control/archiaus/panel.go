package archiaus

import (
	"github.com/go-chassis/go-chassis/control"
	"github.com/go-chassis/go-chassis/core/config/model"
	"github.com/go-chassis/go-chassis/core/invocation"
	"github.com/go-chassis/go-chassis/third_party/forked/afex/hystrix-go/hystrix"
	"github.com/go-mesh/mesher/config"
	meshercontrol "github.com/go-mesh/mesher/control"
)

//Panel pull configs from archaius
type Panel struct {
}

func newPanel(options meshercontrol.Options) control.Panel {
	SaveToEgressCache(config.GetEgressConfig())
	return &Panel{}
}

//GetCircuitBreaker return command , and circuit breaker settings
func (p *Panel) GetCircuitBreaker(inv invocation.Invocation, serviceType string) (string, hystrix.CommandConfig) {
	return "", hystrix.CommandConfig{}
}

//GetLoadBalancing get load balancing config
func (p *Panel) GetLoadBalancing(inv invocation.Invocation) control.LoadBalancingConfig {
	return control.LoadBalancingConfig{}
}

//GetRateLimiting get rate limiting config
func (p *Panel) GetRateLimiting(inv invocation.Invocation, serviceType string) control.RateLimitingConfig {
	return control.RateLimitingConfig{}
}

//GetFaultInjection get Fault injection config
func (p *Panel) GetFaultInjection(inv invocation.Invocation) model.Fault {
	return model.Fault{}

}

//GetEgressRule get egress config
func (p *Panel) GetEgressRule() []control.EgressConfig {
	c, ok := EgressConfigCache.Get("")
	if !ok {

		return nil
	}
	return c.([]control.EgressConfig)
}

func init() {
	meshercontrol.InstallPlugin("egressarchaius", newPanel)
}

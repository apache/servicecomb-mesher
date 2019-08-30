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

package istio

import (
	meshercontrol "github.com/apache/servicecomb-mesher/proxy/control"
	"github.com/apache/servicecomb-mesher/proxy/pkg/egress"
	"github.com/go-chassis/go-chassis/control"
	"github.com/go-chassis/go-chassis/core/config/model"
	"github.com/go-chassis/go-chassis/core/invocation"
	"github.com/go-chassis/go-chassis/third_party/forked/afex/hystrix-go/hystrix"
)

func init() {
	meshercontrol.InstallPlugin("pilot", newPilotPanel)
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

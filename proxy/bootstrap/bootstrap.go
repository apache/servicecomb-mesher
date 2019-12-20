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

package bootstrap

import (
	"log"
	"strings"

	"github.com/apache/servicecomb-mesher/proxy/cmd"
	"github.com/apache/servicecomb-mesher/proxy/common"
	"github.com/apache/servicecomb-mesher/proxy/config"
	"github.com/apache/servicecomb-mesher/proxy/register"
	"github.com/apache/servicecomb-mesher/proxy/resolver"

	"github.com/apache/servicecomb-mesher/proxy/control"
	"github.com/apache/servicecomb-mesher/proxy/pkg/egress"
	"github.com/apache/servicecomb-mesher/proxy/pkg/metrics"
	"github.com/apache/servicecomb-mesher/proxy/pkg/runtime"
	"github.com/apache/servicecomb-mesher/proxy/resource/v1"
	"github.com/apache/servicecomb-mesher/proxy/resource/v1/version"
	"github.com/go-chassis/go-chassis"
	"github.com/go-chassis/go-chassis/core/handler"
	chassisHandler "github.com/go-chassis/go-chassis/core/handler"
	"github.com/go-chassis/go-chassis/core/lager"
	"github.com/go-chassis/go-chassis/core/metadata"
	"github.com/go-mesh/openlogging"
)

// Start initialize configs and components
func Start() error {
	if err := config.InitProtocols(); err != nil {
		return err
	}
	if err := config.Init(); err != nil {
		return err
	}
	if err := resolver.Init(); err != nil {
		return err
	}
	if err := DecideMode(); err != nil {
		return err
	}
	metrics.Init()
	if err := v1.Init(); err != nil {
		log.Println("Error occurred in starting admin server", err)
	}
	if err := register.AdaptEndpoints(); err != nil {
		return err
	}
	if cmd.Configs.LocalServicePorts == "" {
		lager.Logger.Warnf("local service ports is missing, service can not be called by mesher")
	} else {
		lager.Logger.Infof("local service ports is [%v]", cmd.Configs.PortsMap)
	}
	err := egress.Init()
	if err != nil {
		return err
	}

	if err := control.Init(); err != nil {
		return err
	}

	return nil

}

//DecideMode get config mode
func DecideMode() error {
	runtime.Role = cmd.Configs.Role
	openlogging.GetLogger().Info("Running as " + runtime.Role)
	return nil
}

//RegisterFramework registers framework
func RegisterFramework() {
	version := GetVersion()
	if framework := metadata.NewFramework(); cmd.Configs.Role == common.RoleSidecar {
		framework.SetName("Mesher")
		framework.SetVersion(version)
		framework.SetRegister("SIDECAR")
	} else {
		framework.SetName("Mesher")
		framework.SetVersion(version)
	}
}

//GetVersion returns version
func GetVersion() string {
	versionID := version.Ver().Version
	if len(versionID) == 0 {
		return version.DefaultVersion
	}
	return versionID
}

//SetHandlers leverage go-chassis API to set default handlers if there is no define in chassis.yaml
func SetHandlers() {
	consumerChain := strings.Join([]string{
		chassisHandler.Router,
		chassisHandler.RateLimiterConsumer,
		chassisHandler.BizkeeperConsumer,
		chassisHandler.Loadbalance,
		chassisHandler.Transport,
	}, ",")
	providerChain := strings.Join([]string{
		chassisHandler.RateLimiterProvider,
		chassisHandler.Transport,
	}, ",")
	consumerChainMap := map[string]string{
		common.ChainConsumerOutgoing: consumerChain,
	}
	providerChainMap := map[string]string{
		common.ChainProviderIncoming: providerChain,
		"default":                    chassisHandler.RateLimiterProvider,
	}
	chassis.SetDefaultConsumerChains(consumerChainMap)
	chassis.SetDefaultProviderChains(providerChainMap)
}

//InitEgressChain init the egress handler chain
func InitEgressChain() error {
	egresschain := strings.Join([]string{
		handler.RateLimiterConsumer,
		handler.BizkeeperConsumer,
		handler.Transport,
	}, ",")

	egressChainMap := map[string]string{
		common.ChainConsumerEgress: egresschain,
	}

	return handler.CreateChains(common.ConsumerEgress, egressChainMap)
}

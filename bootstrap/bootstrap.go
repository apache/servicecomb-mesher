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

	"github.com/go-mesh/mesher/adminapi"
	"github.com/go-mesh/mesher/adminapi/version"
	"github.com/go-mesh/mesher/cmd"
	"github.com/go-mesh/mesher/common"
	"github.com/go-mesh/mesher/config"
	"github.com/go-mesh/mesher/register"
	"github.com/go-mesh/mesher/resolver"

	"github.com/go-chassis/go-chassis"
	chassisHandler "github.com/go-chassis/go-chassis/core/handler"
	"github.com/go-chassis/go-chassis/core/lager"
	"github.com/go-chassis/go-chassis/core/metadata"
	"github.com/go-mesh/mesher/pkg/metrics"
	"github.com/go-mesh/mesher/pkg/runtime"
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
	if err := adminapi.Init(); err != nil {
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

	return nil

}

//DecideMode get config mode
func DecideMode() error {
	runtime.Mode = cmd.Configs.Mode
	openlogging.GetLogger().Info("Running as " + runtime.Mode)
	return nil
}

//RegisterFramework registers framework
func RegisterFramework() {
	if framework := metadata.NewFramework(); cmd.Configs.Mode == common.ModeSidecar {
		version := GetVersion()
		framework.SetName("Mesher")
		framework.SetVersion(version)
		framework.SetRegister("SIDECAR")
	} else if cmd.Configs.Mode == common.ModePerHost {
		framework.SetName("Mesher")
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
		chassisHandler.RatelimiterConsumer,
		chassisHandler.BizkeeperConsumer,
		chassisHandler.Loadbalance,
		chassisHandler.Transport,
	}, ",")
	providerChain := strings.Join([]string{
		chassisHandler.RatelimiterProvider,
		chassisHandler.Transport,
	}, ",")
	consumerChainMap := map[string]string{
		common.ChainConsumerOutgoing: consumerChain,
	}
	providerChainMap := map[string]string{
		common.ChainProviderIncoming: providerChain,
		"default":                    chassisHandler.RatelimiterProvider,
	}
	chassis.SetDefaultConsumerChains(consumerChainMap)
	chassis.SetDefaultProviderChains(providerChainMap)
}

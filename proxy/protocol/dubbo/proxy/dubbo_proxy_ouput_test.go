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

package dubboproxy

import (
	"fmt"
	"github.com/apache/servicecomb-mesher/proxy/bootstrap"
	"github.com/apache/servicecomb-mesher/proxy/cmd"
	"github.com/apache/servicecomb-mesher/proxy/common"
	mesherCommon "github.com/apache/servicecomb-mesher/proxy/common"
	mesherRuntime "github.com/apache/servicecomb-mesher/proxy/pkg/runtime"
	dubboclient "github.com/apache/servicecomb-mesher/proxy/protocol/dubbo/client"
	"github.com/apache/servicecomb-mesher/proxy/protocol/dubbo/dubbo"
	"github.com/go-chassis/go-chassis"
	chassisCommon "github.com/go-chassis/go-chassis/core/common"
	"github.com/go-chassis/go-chassis/core/config"
	"github.com/go-chassis/go-chassis/core/config/model"
	"github.com/go-chassis/go-chassis/core/handler"
	"github.com/go-chassis/go-chassis/core/invocation"
	"github.com/go-chassis/go-chassis/core/lager"
	"github.com/go-chassis/go-chassis/core/loadbalancer"
	"github.com/go-chassis/go-chassis/core/registry"
	"github.com/go-chassis/go-chassis/third_party/forked/afex/hystrix-go/hystrix"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"os"
	"strings"
	"testing"

	// rate limiter handler
	_ "github.com/go-chassis/go-chassis/middleware/ratelimiter"
)

func init() {
	lager.Init(&lager.Options{LoggerLevel: "DEBUG"})
}

func TestSetLocalServiceAddress(t *testing.T) {
	t.Run("Run not set env EnvServicePorts", func(t *testing.T) {
		inv := &invocation.Invocation{Protocol: "rest"}
		cmd.Init()
		_ = cmd.Configs.GeneratePortsMap()
		t.Log(cmd.Configs.PortsMap)
		err := SetLocalServiceAddress(inv)
		assert.Error(t, err)

		// case port
		inv.Port = "8080"
		err = SetLocalServiceAddress(inv)
		assert.NoError(t, err)
	})

	t.Run("Run with env EnvServicePorts", func(t *testing.T) {
		inv := &invocation.Invocation{Protocol: "rest"}
		os.Setenv(common.EnvServicePorts, "rest:8080,grpc:90")
		cmd.Init()
		_ = cmd.Configs.GeneratePortsMap()
		t.Log(cmd.Configs.PortsMap)
		err := SetLocalServiceAddress(inv)
		assert.NoError(t, err)
	})
}

func TestHandle(t *testing.T) {
	t.Run("Run as Provider", func(t *testing.T) {
		os.Setenv(common.EnvServicePorts, "dubbo:8081,rest:8080")
		cmd.Init()
		_ = cmd.Configs.GeneratePortsMap()

		protoMap := make(map[string]model.Protocol)
		config.GlobalDefinition = &model.GlobalCfg{
			Cse: model.CseStruct{
				Protocols: protoMap,
			},
		}

		bootstrap.RegisterFramework()
		bootstrap.SetHandlers()
		chassis.Init()

		consumerChain := strings.Join([]string{
			handler.Router,
			//"ratelimiter-consumer",
			//"bizkeeper-consumer",
			//handler.Loadbalance,
			//handler.Transport,
		}, ",")
		providerChain := strings.Join([]string{
			//handler.RateLimiterProvider,
			//handler.Transport,
		}, ",")
		consumerChainMap := map[string]string{
			common.ChainConsumerOutgoing: consumerChain,
		}
		providerChainMap := map[string]string{
			common.ChainProviderIncoming: providerChain,
			"default":                    handler.RateLimiterProvider,
		}

		registry.DefaultContractDiscoveryService = new(MockContractDiscoveryService)
		mesherRuntime.Role = mesherCommon.RoleSidecar

		req := dubbo.NewDubboRequest()
		req.SetAttachment(dubbo.PathKey, "hello")
		ctx := &dubbo.InvokeContext{req, &dubbo.DubboRsp{}, nil, "", "127.0.0.1:9090"}
		ctx.Rsp.Init()

		// case get chain error
		err := Handle(ctx)
		assert.Error(t, err)

		// case get chain ok
		handler.CreateChains(chassisCommon.Provider, providerChainMap)
		handler.CreateChains(chassisCommon.Consumer, consumerChainMap)
		err = Handle(ctx)
		assert.NoError(t, err)

		handler.ChainMap = make(map[string]*handler.Chain)
	})

	t.Run("Run as Consumer", func(t *testing.T) {
		//os.Setenv(common.EnvServicePorts, "rest:8080,grpc:90")
		os.Unsetenv(common.EnvServicePorts)
		cmd.Init()
		_ = cmd.Configs.GeneratePortsMap()

		protoMap := make(map[string]model.Protocol)
		config.GlobalDefinition = &model.GlobalCfg{
			Cse: model.CseStruct{
				Protocols: protoMap,
			},
		}

		bootstrap.RegisterFramework()
		bootstrap.SetHandlers()
		chassis.Init()

		consumerChain := strings.Join([]string{
			handler.Router,
			//"ratelimiter-consumer",
			//"bizkeeper-consumer",
			//handler.Loadbalance,
			//handler.Transport,
		}, ",")
		providerChain := strings.Join([]string{
			//handler.RateLimiterProvider,
			//handler.Transport,
		}, ",")
		consumerChainMap := map[string]string{
			common.ChainConsumerOutgoing: consumerChain,
		}
		providerChainMap := map[string]string{
			common.ChainProviderIncoming: providerChain,
			"default":                    handler.RateLimiterProvider,
		}

		registry.DefaultContractDiscoveryService = new(MockContractDiscoveryService)
		mesherRuntime.Role = mesherCommon.RoleSidecar

		req := dubbo.NewDubboRequest()
		req.SetAttachment(dubbo.PathKey, "hello")
		ctx := &dubbo.InvokeContext{req, &dubbo.DubboRsp{}, nil, "", "127.0.0.1:9090"}
		ctx.Rsp.Init()

		// case get chain error
		err := Handle(ctx)
		assert.Error(t, err)

		// case get chain ok
		handler.CreateChains(chassisCommon.Provider, providerChainMap)
		handler.CreateChains(chassisCommon.Consumer, consumerChainMap)
		err = Handle(ctx)
		assert.NoError(t, err)

		handler.ChainMap = make(map[string]*handler.Chain)
	})
}

func Test_handleDubboRequest(t *testing.T) {
	req := dubbo.NewDubboRequest()
	req.SetAttachment(dubbo.PathKey, "hello")
	ctx := &dubbo.InvokeContext{req, &dubbo.DubboRsp{}, nil, "", "127.0.0.1:9090"}
	ctx.Rsp.Init()

	inv := &invocation.Invocation{Protocol: "rest"}
	ir := &invocation.Response{}
	inv.Endpoint = "127.0.0.1:8080"

	// case responese ir.Result = nil
	handleDubboRequest(inv, ctx, ir)

	// case responese ir.Result != nil
	ir.Result = &dubboclient.WrapResponse{Resp: &dubbo.DubboRsp{}}
	handleDubboRequest(inv, ctx, ir)

	// Case ir.Err == hystrix.FallbackNullError
	ir.Err = hystrix.FallbackNullError{"Error."}
	handleDubboRequest(inv, ctx, ir)
	// Case ir.Err == hystrix.CircuitError:
	ir.Err = hystrix.CircuitError{"Error."}
	handleDubboRequest(inv, ctx, ir)
	// Case ir.Err == loadbalancer.LBError
	ir.Err = loadbalancer.LBError{"Error."}
	handleDubboRequest(inv, ctx, ir)
	// Case ir.Err == other
	ir.Err = fmt.Errorf("Other error.")
	handleDubboRequest(inv, ctx, ir)

	// case ir == nil
	handleDubboRequest(inv, ctx, nil)

}

// ContractDiscoveryService struct for disco mock
type MockContractDiscoveryService struct {
	mock.Mock
}

func (m *MockContractDiscoveryService) GetMicroServicesByInterface(interfaceName string) (microservices []*registry.MicroService) {
	microservices = append(microservices, &registry.MicroService{})
	return
}

func (m *MockContractDiscoveryService) GetSchemaContentByInterface(interfaceName string) registry.SchemaContent {
	return registry.SchemaContent{}
}

func (m *MockContractDiscoveryService) GetSchemaContentByServiceName(svcName, version, appID, env string) []*registry.SchemaContent {
	var sc []*registry.SchemaContent
	sc = append(sc, &registry.SchemaContent{
		Paths: map[string]map[string]registry.MethodInfo{
			"hello": {},
		},
	})
	return nil
}

func (m *MockContractDiscoveryService) Close() error {
	return nil
}

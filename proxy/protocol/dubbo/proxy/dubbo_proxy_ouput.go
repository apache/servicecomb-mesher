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
	"context"
	"fmt"
	"github.com/apache/servicecomb-mesher/proxy/cmd"
	mesherCommon "github.com/apache/servicecomb-mesher/proxy/common"
	mesherRuntime "github.com/apache/servicecomb-mesher/proxy/pkg/runtime"
	"github.com/apache/servicecomb-mesher/proxy/protocol"
	"github.com/apache/servicecomb-mesher/proxy/protocol/dubbo/client"
	"github.com/apache/servicecomb-mesher/proxy/protocol/dubbo/dubbo"
	"github.com/apache/servicecomb-mesher/proxy/protocol/dubbo/schema"
	"github.com/apache/servicecomb-mesher/proxy/protocol/dubbo/utils"
	"github.com/apache/servicecomb-mesher/proxy/resolver"
	"github.com/go-chassis/go-chassis/core/common"
	chassisCommon "github.com/go-chassis/go-chassis/core/common"
	"github.com/go-chassis/go-chassis/core/handler"
	"github.com/go-chassis/go-chassis/core/invocation"
	"github.com/go-chassis/go-chassis/core/lager"
	"github.com/go-chassis/go-chassis/core/loadbalancer"
	"github.com/go-chassis/go-chassis/pkg/runtime"
	"github.com/go-chassis/go-chassis/pkg/util/tags"
	"github.com/go-chassis/go-chassis/third_party/forked/afex/hystrix-go/hystrix"
	"github.com/go-mesh/openlogging"
)

var dr = resolver.GetDestinationResolver("http")
var sr = resolver.GetSourceResolver()

const (
	ProxyTag = "mesherproxy"
)

//IsProvider is variable of type boolean used for tag proxyed dubbo service as provider(true) or consumer(false)
var IsProvider bool

// DubboListenAddr is a variable of type string used for storing listen address
var DubboListenAddr string

//ProxyError is a struct
type ProxyError struct {
	Message string
}

//Error is a method which returns error
func (e ProxyError) Error() string {
	return e.Message
}

//SetLocalServiceAddress assign invocation endpoint a local service address
// it uses config in cmd or env fi
// if it is empty, then try to use original port from client as local port
func SetLocalServiceAddress(inv *invocation.Invocation) error {
	inv.Endpoint = cmd.Configs.PortsMap[inv.Protocol]
	if inv.Endpoint == "" {
		if inv.Port != "" {
			inv.Endpoint = "127.0.0.1:" + inv.Port
			cmd.Configs.PortsMap[inv.Protocol] = inv.Endpoint
			return nil
		} else {
			return fmt.Errorf("[%s] is not supported, [%s] didn't set env [%s] or cmd parameter --service-ports before mesher start",
				inv.Protocol, inv.MicroServiceName, mesherCommon.EnvServicePorts)
		}
	}
	return nil
}

//Handle is a function
func Handle(ctx *dubbo.InvokeContext) error {
	interfaceName := ctx.Req.GetAttachment(dubbo.PathKey, "")
	svc := schema.GetSvcByInterface(interfaceName)
	if svc == nil {
		return &util.BaseError{ErrMsg: "can't find the svc by " + interfaceName}
	}

	inv := new(invocation.Invocation)
	inv.SourceServiceID = runtime.ServiceID
	inv.SourceMicroService = ctx.Req.GetAttachment(common.HeaderSourceName, "")
	inv.Args = ctx.Req
	inv.Ctx = context.WithValue(context.Background(), chassisCommon.ContextHeaderKey{}, ctx.Req.GetAttachments())
	inv.MicroServiceName = svc.ServiceName
	inv.RouteTags = utiltags.NewDefaultTag(svc.Version, svc.AppID)
	inv.Protocol = "dubbo"
	inv.URLPathFormat = ""
	inv.Reply = &dubboclient.WrapResponse{nil} //&rest.Response{Resp: &ctx.Response}
	var err error
	err = SetLocalServiceAddress(inv) //select local service
	if err != nil {
		openlogging.GetLogger().Warn(err.Error())
		IsProvider = false
	} else {
		IsProvider = true
	}

	var c *handler.Chain
	//发送请求
	//value := ctx.Req.GetAttachment(ProxyTag, "")
	if !IsProvider || inv.MicroServiceName != runtime.ServiceName { //come from proxyedDubboSvc
		ctx.Req.SetAttachment(common.HeaderSourceName, runtime.ServiceName)
		ctx.Req.SetAttachment(ProxyTag, "true")

		if mesherRuntime.Role == mesherCommon.RoleSidecar {
			c, err = handler.GetChain(common.Consumer, mesherCommon.ChainConsumerOutgoing)
			if err != nil {
				openlogging.Error("Get Consumer chain failed: " + err.Error())
				return err
			}
		}
		c.Next(inv, func(ir *invocation.Response) error {
			return handleDubboRequest(inv, ctx, ir)
		})
	} else { //come from other mesher
		ctx.Req.SetAttachment(ProxyTag, "")
		c, err = handler.GetChain(common.Provider, mesherCommon.ChainProviderIncoming)
		if err != nil {
			openlogging.Error("Get Provider Chain failed: " + err.Error())
			return err
		}
		c.Next(inv, func(ir *invocation.Response) error {
			return handleDubboRequest(inv, ctx, ir)
		})
	}

	return nil
}

func handleDubboRequest(inv *invocation.Invocation, ctx *dubbo.InvokeContext, ir *invocation.Response) error {
	if ir != nil {
		if ir.Err != nil {
			switch ir.Err.(type) {
			case hystrix.FallbackNullError:
				ctx.Rsp.SetStatus(dubbo.Ok)
			case hystrix.CircuitError:
				ctx.Rsp.SetStatus(dubbo.ServiceError)
			case loadbalancer.LBError:
				ctx.Rsp.SetStatus(dubbo.ServiceNotFound)
			default:
				ctx.Rsp.SetStatus(dubbo.ServerError)
			}
			ctx.Rsp.SetErrorMsg(ir.Err.Error())
			return ir.Err
		}
		if inv.Endpoint == "" {
			ctx.Rsp.SetStatus(dubbo.ServerError)
			ctx.Rsp.SetErrorMsg(protocol.ErrUnknown.Error())
			return protocol.ErrUnknown
		}
	} else {
		ctx.Rsp.SetStatus(dubbo.ServerError)
		ctx.Rsp.SetErrorMsg(protocol.ErrUnExpectedHandlerChainResponse.Error())
		return protocol.ErrUnExpectedHandlerChainResponse
	}
	if ir.Result != nil {
		ctx.Rsp = ir.Result.(*dubboclient.WrapResponse).Resp
	} else {
		err := protocol.ErrNilResult
		lager.Logger.Error("CAll Chain  failed: " + err.Error())
		return err
	}

	return nil
}

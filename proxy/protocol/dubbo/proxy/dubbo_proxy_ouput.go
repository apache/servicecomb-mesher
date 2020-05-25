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
	"encoding/json"
	"fmt"
	"github.com/apache/servicecomb-mesher/proxy/cmd"
	"net/http"
	"net/url"

	mesherCommon "github.com/apache/servicecomb-mesher/proxy/common"
	mesherRuntime "github.com/apache/servicecomb-mesher/proxy/pkg/runtime"
	"github.com/apache/servicecomb-mesher/proxy/protocol"
	"github.com/apache/servicecomb-mesher/proxy/protocol/dubbo/client"
	"github.com/apache/servicecomb-mesher/proxy/protocol/dubbo/dubbo"
	"github.com/apache/servicecomb-mesher/proxy/protocol/dubbo/schema"
	"github.com/apache/servicecomb-mesher/proxy/protocol/dubbo/utils"
	"github.com/apache/servicecomb-mesher/proxy/resolver"
	"github.com/go-chassis/go-chassis/client/rest"
	"github.com/go-chassis/go-chassis/core/common"
	chassisCommon "github.com/go-chassis/go-chassis/core/common"
	"github.com/go-chassis/go-chassis/core/handler"
	"github.com/go-chassis/go-chassis/core/invocation"
	"github.com/go-chassis/go-chassis/core/lager"
	"github.com/go-chassis/go-chassis/core/loadbalancer"
	"github.com/go-chassis/go-chassis/pkg/runtime"
	stringutil "github.com/go-chassis/go-chassis/pkg/string"
	"github.com/go-chassis/go-chassis/pkg/util/httputil"
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

//ConvertDubboReqToHTTPReq is a method which converts dubbo requesto to http request
func ConvertDubboReqToHTTPReq(ctx *dubbo.InvokeContext, dubboReq *dubbo.Request) *http.Request {
	restReq := &http.Request{
		URL:    &url.URL{},
		Header: make(http.Header),
	}
	args := dubboReq.GetArguments()
	operateID := dubboReq.GetMethodName()
	iName := dubboReq.GetAttachment(dubbo.PathKey, "")

	methd := schema.GetMethodByInterface(iName, operateID)
	if methd == nil {
		lager.Logger.Error("GetMethodByInterface failed: Cannot find the method")
		return nil
	}
	ctx.Method = methd
	restReq.Method = methd.Verb

	var (
		i         = 0
		qureyNum  = 0
		paramsStr = "?"
		body      = []byte{}
	)

	for i = 0; i < len(args); i++ {
		_, in := methd.GetParamNameAndWhere(i)
		paraSchema := methd.GetParamSchema(i)
		v := args[i]
		if in == schema.InBody {
			b, _ := json.Marshal(v.GetValue())
			body = append(body, b...)
		} else {
			var fmtStr string
			var value string
			if paraSchema.Dtype == util.SchemaArray {
				value = util.ArrayToQueryString(paraSchema.Name, v.GetValue())
				fmtStr += value
			} else {
				value, _ = util.ObjectToString(paraSchema.Dtype, v.GetValue()) // (v.GetValue()).(string)
				if qureyNum == 0 {
					fmtStr = fmt.Sprintf("%s=%s", paraSchema.Name, url.QueryEscape(value))
					qureyNum++
				} else {
					fmtStr = fmt.Sprintf("&%s=%s", paraSchema.Name, url.QueryEscape(value))
				}
			}
			paramsStr += fmtStr
		}
	}
	httputil.SetBody(restReq, body)

	uri := methd.Path
	if paramsStr != "?" {
		uri += paramsStr
	}
	httputil.SetURI(restReq, uri)
	tmpName := schema.GetSvcNameByInterface(iName)
	if tmpName == "" {
		lager.Logger.Error("GetSvcNameByInterface failed: Cannot find the svc")
		return nil
	}
	restReq.URL.Host = tmpName // must after setURI
	return restReq
}

//ConvertRestRspToDubboRsp is a function which converts rest response to dubbo response
func ConvertRestRspToDubboRsp(ctx *dubbo.InvokeContext, resp *http.Response, dubboRsp *dubbo.DubboRsp) {
	var v interface{}
	var err error
	status := resp.StatusCode
	body := httputil.ReadBody(resp)
	if status >= http.StatusBadRequest {
		dubboRsp.SetStatus(dubbo.ServerError)
		if dubboRsp.GetErrorMsg() == "" && body != nil {
			dubboRsp.SetErrorMsg(string(body))
		}
		return
	}
	dubboRsp.SetStatus(dubbo.Ok)
	if body != nil {
		rspSchema := (*(ctx.Method)).GetRspSchema(status)
		if rspSchema != nil {
			v, err = util.RestByteToValue(rspSchema.DType, body)
			if err != nil {
				dubboRsp.SetStatus(dubbo.BadResponse)
				dubboRsp.SetErrorMsg(err.Error())
			} else {
				dubboRsp.SetValue(v)
			}
		} else {
			dubboRsp.SetErrorMsg(string(body))
			dubboRsp.SetStatus(dubbo.ServerError)
		}
	}

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
	} else {
		IsProvider = true
	}

	var c *handler.Chain
	if inv.Protocol == "dubbo" {
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
	} else {
		return ProxyRestHandler(ctx)
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

func preHandleToRest(ctx *dubbo.InvokeContext) (*http.Request, *invocation.Invocation, string) {
	restReq := ConvertDubboReqToHTTPReq(ctx, ctx.Req)
	if restReq == nil {
		return nil, nil, ""
	}
	inv := new(invocation.Invocation)
	inv.SourceServiceID = runtime.ServiceID
	inv.Args = restReq
	inv.Protocol = "rest"
	inv.Reply = rest.NewResponse()
	inv.URLPathFormat = restReq.URL.String()
	inv.SchemaID = ""
	inv.OperationID = ""
	inv.Ctx = context.Background()
	err := SetLocalServiceAddress(inv) //select local service
	if err != nil {
		openlogging.Error(err.Error())
	}
	source := stringutil.SplitFirstSep(ctx.RemoteAddr, ":")
	return restReq, inv, source
}

//ProxyRestHandler is a function
func ProxyRestHandler(ctx *dubbo.InvokeContext) error {
	var err error
	var c *handler.Chain

	req, inv, source := preHandleToRest(ctx)
	if req == nil {
		return &util.BaseError{ErrMsg: "request is invalid "}
	}

	source = "127.0.0.1" //"10.57.75.87"
	//Resolve Source
	si := sr.Resolve(source)
	h := make(map[string]string)
	for k := range req.Header {
		h[k] = req.Header.Get(k)
	}
	//Resolve Destination
	destination, _, err := dr.Resolve(source, "", inv.URLPathFormat, h)
	if err != nil {
		return err
	}
	inv.MicroServiceName = destination
	if mesherRuntime.Role == mesherCommon.RoleSidecar {
		c, err = handler.GetChain(common.Consumer, mesherCommon.ChainConsumerOutgoing)
		if err != nil {
			lager.Logger.Error("Get chain failed: " + err.Error())
			return err
		}
		if si == nil {
			lager.Logger.Info("Can not resolve " + source + " to Source info")
		}
	}

	c.Next(inv, func(ir *invocation.Response) error {
		//Send the request to the destination
		return handleRequest(ctx, req, inv.Reply.(*http.Response), ctx.Rsp, inv, ir)
	})
	ConvertRestRspToDubboRsp(ctx, inv.Reply.(*http.Response), ctx.Rsp)
	return nil
}

func handleRequest(ctx *dubbo.InvokeContext, req *http.Request, resp *http.Response,
	dubboRsp *dubbo.DubboRsp, inv *invocation.Invocation, ir *invocation.Response) error {
	if ir != nil {
		if ir.Err != nil {
			switch ir.Err.(type) {
			case hystrix.FallbackNullError:
				resp.StatusCode = http.StatusOK
				dubboRsp.SetErrorMsg(ir.Err.Error())
			case hystrix.CircuitError:
				ir.Status = http.StatusServiceUnavailable
				resp.StatusCode = http.StatusServiceUnavailable
				dubboRsp.SetErrorMsg(ir.Err.Error())
			case loadbalancer.LBError:
				ir.Status = http.StatusBadGateway
				resp.StatusCode = http.StatusBadGateway
				dubboRsp.SetErrorMsg(ir.Err.Error())
			default:
				ir.Status = http.StatusInternalServerError
				resp.StatusCode = http.StatusInternalServerError
				dubboRsp.SetErrorMsg(ir.Err.Error())
			}
			return ir.Err
		}
		if inv.Endpoint == "" {
			ir.Status = http.StatusInternalServerError
			resp.StatusCode = http.StatusInternalServerError
			dubboRsp.SetErrorMsg(ir.Err.Error())
			return protocol.ErrUnknown
		}
	} else {
		dubboRsp.SetErrorMsg(protocol.ErrUnExpectedHandlerChainResponse.Error())
		return protocol.ErrUnExpectedHandlerChainResponse
	}

	ir.Status = resp.StatusCode
	return nil
}

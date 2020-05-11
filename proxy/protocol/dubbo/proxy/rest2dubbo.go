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
	"github.com/go-chassis/go-chassis/pkg/runtime"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	mesherCommon "github.com/apache/servicecomb-mesher/proxy/common"
	"github.com/apache/servicecomb-mesher/proxy/protocol/dubbo/client"
	"github.com/apache/servicecomb-mesher/proxy/protocol/dubbo/dubbo"
	"github.com/apache/servicecomb-mesher/proxy/protocol/dubbo/schema"
	"github.com/apache/servicecomb-mesher/proxy/protocol/dubbo/utils"

	"github.com/apache/servicecomb-mesher/proxy/protocol"
	"github.com/go-chassis/go-chassis/core/common"
	"github.com/go-chassis/go-chassis/core/handler"
	"github.com/go-chassis/go-chassis/core/invocation"
	"github.com/go-chassis/go-chassis/core/loadbalancer"
	"github.com/go-chassis/go-chassis/pkg/string"
	"github.com/go-chassis/go-chassis/pkg/util/tags"
	"github.com/go-chassis/go-chassis/third_party/forked/afex/hystrix-go/hystrix"
	"github.com/go-mesh/openlogging"
)

//ConvertDubboRspToRestRsp is a function which converts dubbo response to rest response
func ConvertDubboRspToRestRsp(dubboRsp *dubbo.DubboRsp, w http.ResponseWriter, ctx *dubbo.InvokeContext) error {
	status := dubboRsp.GetStatus()
	if status == dubbo.Ok {
		w.WriteHeader(http.StatusOK)
		rspSchema := (*(ctx.Method)).GetRspSchema(http.StatusOK)
		if rspSchema != nil {
			v, err := util.ObjectToString(rspSchema.DType, dubboRsp.GetValue())
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			} else {
				if _, err := w.Write([]byte(v)); err != nil {
					return err
				}
			}
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
	return nil
}

//ConvertHTTPReqToDubboReq is a function which converts http request in to dubbo request
func ConvertHTTPReqToDubboReq(restReq *http.Request, ctx *dubbo.InvokeContext, inv *invocation.Invocation) error {
	req := ctx.Req
	uri := restReq.URL
	i := 0
	var dubboArgs []util.Argument
	queryAgrs := uri.Query()
	arg := &util.Argument{}

	svcSchema, methd := schema.GetSchemaMethodBySvcURL(inv.MicroServiceName, "", inv.RouteTags.Version(), inv.RouteTags.AppID(),
		strings.ToLower(restReq.Method), string(restReq.URL.String()))
	if methd == nil {
		return &util.BaseError{"Method not been found"}
	}
	req.SetMethodName(methd.OperaID)
	req.SetAttachment(dubbo.DubboVersionKey, dubbo.DubboVersion)
	req.SetAttachment(dubbo.PathKey, svcSchema.Info["x-java-interface"]) //interfaceSchema.JavaClsName
	req.SetAttachment(dubbo.VersionKey, "0.0.0")
	ctx.Method = methd
	var err error

	//处理参数
	dubboArgs = make([]util.Argument, len(methd.Paras))

	for _, v := range methd.Paras {
		var byteTmp []byte
		var bytesTmp [][]byte
		itemType := "string" //默认为string
		if strings.EqualFold(v.Where, "query") {
			byteTmp = []byte(queryAgrs.Get(v.Name))
		} else if restReq.Body != nil {
			byteTmp, _ = ioutil.ReadAll(restReq.Body)
		}
		if byteTmp == nil && v.Required {
			return &util.BaseError{"Param is null"}
		}
		var realJvmType string
		bytesTmp, realJvmType = getJVMType(v, arg, bytesTmp, restReq.URL)
		if bytesTmp == nil {
			arg.Value, err = util.RestByteToValue(arg.JavaType, byteTmp)
			if err != nil {
				return err
			}
		} else {
			arg.Value, err = util.RestBytesToLstValue(itemType, bytesTmp)
			if err != nil {
				return err
			}
		}

		if realJvmType != "" {
			arg.JavaType = realJvmType
		}
		dubboArgs[i] = *arg
		i++
	}

	req.SetArguments(dubboArgs)

	return nil
}

func getJVMType(v schema.MethParam, arg *util.Argument, bytesTmp [][]byte, queryAgrs *url.URL) ([][]byte, string) {
	var realJvmType string
	queryAgrsTmp := queryAgrs.Query()
	if _, ok := util.SchemeTypeMAP[v.Dtype]; ok {
		arg.JavaType = util.SchemeTypeMAP[v.Dtype]
		if v.Dtype == util.SchemaArray {
			realJvmType = util.JavaList
			if v.Items != nil {
				if val, ok := v.Items["x-java-class"]; ok {
					realJvmType = fmt.Sprintf("L%s;", val)
				}
				if valType, ok := v.Items["type"]; ok {
					realJvmType = fmt.Sprintf("L%s;", valType)
				}
			}
			bytesTmp = util.S2ByteSlice(queryAgrsTmp[v.Name])
		} else if arg.JavaType == util.JavaObject {
			realJvmType = fmt.Sprintf("L%s;", v.ObjRef.JvmClsName)
			if v.AdditionalProps != nil { //处理map
				if val, ok := v.AdditionalProps["x-java-class"]; ok {
					realJvmType = fmt.Sprintf("L%s;", val)
				} else {
					realJvmType = util.JavaMap
				}
			}
		}
		//Lcom.alibaba.dubbo.demo.user; need convert to  Lcom/alibaba/dubbo/demo/User;
		realJvmType = strings.Replace(realJvmType, ".", "/", -1)
	}
	return bytesTmp, realJvmType
}

func preHandleToDubbo(req *http.Request) (*invocation.Invocation, string) {
	inv := new(invocation.Invocation)
	inv.MicroServiceName = runtime.ServiceName
	inv.RouteTags = utiltags.NewDefaultTag(runtime.Version, runtime.App)

	inv.Protocol = "dubbo"
	inv.URLPathFormat = req.URL.Path
	inv.Reply = &dubboclient.WrapResponse{nil}
	source := stringutil.SplitFirstSep(req.RemoteAddr, ":")
	return inv, source
}

//TransparentForwardHandler is a function
func TransparentForwardHandler(w http.ResponseWriter, r *http.Request) {
	inv, _ := preHandleToDubbo(r)
	dubboCtx := &dubbo.InvokeContext{dubbo.NewDubboRequest(), &dubbo.DubboRsp{}, nil, "", ""}
	err := ConvertHTTPReqToDubboReq(r, dubboCtx, inv)
	if err != nil {
		openlogging.Error("Invalid Request: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	inv.Args = dubboCtx.Req

	c, err := handler.GetChain(common.Provider, mesherCommon.ChainProviderIncoming)
	if err != nil {
		openlogging.Error("Get Chain failed: " + err.Error())
		return
	}
	c.Next(inv, func(ir *invocation.Response) error {
		return handleRequestForDubbo(w, inv, ir)
	})
	dubboRsp := inv.Reply.(*dubboclient.WrapResponse).Resp
	if dubboRsp != nil {
		err := ConvertDubboRspToRestRsp(dubboRsp, w, dubboCtx)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(stringutil.Str2bytes(err.Error()))
		}
	}
}

func handleRequestForDubbo(w http.ResponseWriter, inv *invocation.Invocation, ir *invocation.Response) error {
	if ir != nil {
		var err error
		if ir.Err != nil {
			switch ir.Err.(type) {
			case hystrix.FallbackNullError:
				w.WriteHeader(http.StatusOK)
				ir.Status = http.StatusOK
			case hystrix.CircuitError:
				w.WriteHeader(http.StatusServiceUnavailable)
				ir.Status = http.StatusServiceUnavailable
				_, err = w.Write([]byte(ir.Err.Error()))
			case loadbalancer.LBError:
				w.WriteHeader(http.StatusBadGateway)
				ir.Status = http.StatusBadGateway
				_, err = w.Write([]byte(ir.Err.Error()))
			default:
				w.WriteHeader(http.StatusInternalServerError)
				ir.Status = http.StatusInternalServerError
				_, err = w.Write([]byte(ir.Err.Error()))
			}
			if err != nil {
				return err
			}

			return ir.Err
		}
		if inv.Endpoint == "" {
			w.WriteHeader(http.StatusInternalServerError)
			ir.Status = http.StatusInternalServerError
			_, err = w.Write([]byte(protocol.ErrUnknown.Error()))
			if err != nil {
				return err
			}
			return protocol.ErrUnknown
		}
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(protocol.ErrUnExpectedHandlerChainResponse.Error()))
		if err != nil {
			return err
		}
		return protocol.ErrUnExpectedHandlerChainResponse
	}

	return nil
}

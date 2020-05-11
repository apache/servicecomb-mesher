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

package http

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/apache/servicecomb-mesher/proxy/common"
	"github.com/apache/servicecomb-mesher/proxy/pkg/egress"
	"github.com/apache/servicecomb-mesher/proxy/pkg/metrics"
	"github.com/apache/servicecomb-mesher/proxy/protocol"
	"github.com/apache/servicecomb-mesher/proxy/resolver"
	"github.com/apache/servicecomb-mesher/proxy/util"
	"github.com/go-chassis/go-chassis/client/rest"
	chassisCommon "github.com/go-chassis/go-chassis/core/common"
	"github.com/go-chassis/go-chassis/core/fault"
	"github.com/go-chassis/go-chassis/core/handler"
	"github.com/go-chassis/go-chassis/core/invocation"
	"github.com/go-chassis/go-chassis/core/lager"
	"github.com/go-chassis/go-chassis/core/loadbalancer"
	"github.com/go-chassis/go-chassis/pkg/runtime"
	"github.com/go-chassis/go-chassis/pkg/string"
	"github.com/go-chassis/go-chassis/pkg/util/tags"
	"github.com/go-chassis/go-chassis/third_party/forked/afex/hystrix-go/hystrix"
	"github.com/go-mesh/openlogging"
)

var dr = resolver.GetDestinationResolver("http")
var sr = resolver.GetSourceResolver()

//constants for headers
const (
	XForwardedPort = "X-Forwarded-Port"
	XForwardedHost = "X-Forwarded-Host"
)

var (
	//ErrRestFaultAbort is a variable of type error
	ErrRestFaultAbort = errors.New("injecting abort")
	//ErrRestFault is a variable of type error
	ErrRestFault = errors.New("injecting abort and delay")
	//ErrNilResponse is a variable of type error
	ErrNilResponse = errors.New("http response is nil")
)

func preHandler(req *http.Request) *invocation.Invocation {
	inv := &invocation.Invocation{}
	inv.Args = req
	inv.Reply = rest.NewResponse()
	inv.Protocol = "rest"
	inv.URLPathFormat = req.URL.Path
	return inv
}

func consumerPreHandler(req *http.Request) *invocation.Invocation {
	inv := preHandler(req)
	inv.SourceServiceID = runtime.ServiceID
	req.Header.Set(chassisCommon.HeaderSourceName, runtime.ServiceName)
	inv.Ctx = context.TODO()
	return inv
}

func providerPreHandler(req *http.Request) *invocation.Invocation {
	inv := preHandler(req)
	inv.MicroServiceName = runtime.ServiceName
	inv.RouteTags = utiltags.NewDefaultTag(runtime.Version, runtime.App)
	inv.SourceMicroService = req.Header.Get(chassisCommon.HeaderSourceName)
	inv.Ctx = context.TODO()
	return inv
}

//LocalRequestHandler is for request from local
func LocalRequestHandler(w http.ResponseWriter, r *http.Request) {
	prepareRequest(r)
	inv := consumerPreHandler(r)
	remoteIP := stringutil.SplitFirstSep(r.RemoteAddr, ":")

	var err error
	h := make(map[string]string)
	for k := range r.Header {
		h[k] = r.Header.Get(k)
	}
	//Resolve Destination
	destination, port, err := dr.Resolve(remoteIP, r.Host, r.URL.String(), h)
	if err != nil {
		handleErrorResponse(inv, w, http.StatusBadRequest, err)
		return
	}
	inv.MicroServiceName = destination
	if port != "" {
		h[XForwardedPort] = port
	}

	//transfer header into ctx
	inv.Ctx = context.WithValue(inv.Ctx, chassisCommon.ContextHeaderKey{}, h)

	var c *handler.Chain
	ok, egressRule := egress.Match(inv.MicroServiceName)
	if ok {
		var targetPort int32 = 80
		for _, port := range egressRule.Ports {
			if strings.EqualFold(port.Protocol, common.HTTPProtocol) {
				targetPort = port.Port
				break
			}
		}
		inv.Endpoint = inv.MicroServiceName + ":" + strconv.Itoa(int(targetPort))
		c, err = handler.GetChain(common.ConsumerEgress, common.ChainConsumerEgress)
		if err != nil {
			handleErrorResponse(inv, w, http.StatusBadGateway, err)
			openlogging.Error("Get chain failed" + err.Error())
			return
		}

	} else {
		c, err = handler.GetChain(chassisCommon.Consumer, common.ChainConsumerOutgoing)
		if err != nil {
			handleErrorResponse(inv, w, http.StatusBadGateway, err)
			openlogging.Error("Get chain failed: " + err.Error())
			return
		}
	}
	defer func(begin time.Time) {
		timeTaken := time.Since(begin).Seconds()
		serviceLabelValues := map[string]string{metrics.LServiceName: inv.MicroServiceName, metrics.LApp: inv.RouteTags.AppID(), metrics.LVersion: inv.RouteTags.Version()}
		metrics.RecordLatency(serviceLabelValues, timeTaken)
	}(time.Now())
	var invRsp *invocation.Response
	c.Next(inv, func(ir *invocation.Response) error {
		//Send the request to the destination
		invRsp = ir
		if invRsp != nil {
			return invRsp.Err
		}
		return nil
	})
	resp, err := handleRequest(w, inv, invRsp)
	if err != nil {
		openlogging.Error("handle request failed: " + err.Error())
		return
	}
	RecordStatus(inv, resp.StatusCode)
}

//RemoteRequestHandler is for request from remote
func RemoteRequestHandler(w http.ResponseWriter, r *http.Request) {
	prepareRequest(r)
	inv := providerPreHandler(r)

	if inv.SourceMicroService == "" {
		source := stringutil.SplitFirstSep(r.RemoteAddr, ":")
		//Resolve Source
		si := sr.Resolve(source)
		if si != nil {
			inv.SourceMicroService = si.Name
		}
	}
	h := make(map[string]string)
	for k := range r.Header {
		h[k] = r.Header.Get(k)
	}
	//transfer header into ctx
	inv.Ctx = context.WithValue(inv.Ctx, chassisCommon.ContextHeaderKey{}, h)
	c, err := handler.GetChain(chassisCommon.Provider, common.ChainProviderIncoming)
	if err != nil {
		handleErrorResponse(inv, w, http.StatusBadGateway, err)
		lager.Logger.Error("Get chain failed: " + err.Error())
		return
	}
	if err = util.SetLocalServiceAddress(inv, r.Header.Get("X-Forwarded-Port")); err != nil {
		handleErrorResponse(inv, w, http.StatusBadGateway,
			err)
	}
	if r.Header.Get(XForwardedHost) == "" {
		r.Header.Set(XForwardedHost, r.Host)
	}
	var invRsp *invocation.Response
	c.Next(inv, func(ir *invocation.Response) error {
		//Send the request to the destination
		invRsp = ir
		if invRsp != nil {
			return invRsp.Err
		}
		return nil
	})
	if _, err = handleRequest(w, inv, invRsp); err != nil {
		lager.Logger.Error("Handle request failed: " + err.Error())
	}
}

func copyChassisResp2HttpResp(w http.ResponseWriter, resp *http.Response) {
	postProcessResponse(resp)
	copyHeader(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)
	if resp == nil {
		openlogging.GetLogger().Warn("response is nil because of unknown reason")
		return
	}
	_, err := io.Copy(w, resp.Body)
	if err != nil {
		openlogging.Error("can not copy: " + err.Error())
	}
	err = resp.Body.Close()
	if err != nil {
		openlogging.Error("Http response close error: " + err.Error())
	}
}
func handleRequest(w http.ResponseWriter, inv *invocation.Invocation, ir *invocation.Response) (*http.Response, error) {
	if ir != nil {
		if ir.Err != nil {
			//handler only mesher errors, ignore http response err
			switch ir.Err.(type) {
			case hystrix.FallbackNullError:
				handleErrorResponse(inv, w, http.StatusOK, nil)
			case loadbalancer.LBError:
				handleErrorResponse(inv, w, http.StatusBadGateway, ir.Err)
			case hystrix.CircuitError:
				handleErrorResponse(inv, w, http.StatusServiceUnavailable, ir.Err)
			case fault.Fault:
				handleErrorResponse(inv, w, ir.Status, ir.Err)
			default: //for other error, check response and response body, if there is body, just transparent response
				resp, ok := inv.Reply.(*http.Response)

				if ok { // return raw transport error
					if resp != nil {
						if resp.Body == nil {
							//resp.Resp can be nil, for example network error, must handle it
							handleErrorResponse(inv, w, http.StatusBadGateway, ir.Err)
							return nil, ir.Err
						}
						copyChassisResp2HttpResp(w, resp)
						RecordStatus(inv, resp.StatusCode)
					} else {
						// unknown error, resp is nil, e.g. connection refused
						handleErrorResponse(inv, w, http.StatusBadGateway, ir.Err)
					}
				} else { // unknown err in handler chain
					handleErrorResponse(inv, w, http.StatusInternalServerError, ir.Err)
				}
			}
			return nil, ir.Err
		}
		if inv.Endpoint == "" {
			handleErrorResponse(inv, w, http.StatusBadGateway, protocol.ErrUnknown)
			return nil, protocol.ErrUnknown
		}
		if ir.Result == nil {
			if ir.Err != nil {
				handleErrorResponse(inv, w, http.StatusBadGateway, ir.Err)
				return nil, ir.Err
			}
			handleErrorResponse(inv, w, http.StatusBadGateway, ErrNilResponse)
			return nil, protocol.ErrUnknown
		}
		resp, ok := ir.Result.(*http.Response)
		if !ok {
			err := errors.New("invocationResponse result is not type *rest.Response")
			handleErrorResponse(inv, w, http.StatusBadGateway, err)
			return nil, err
		}
		//transparent proxy
		copyChassisResp2HttpResp(w, resp)

		return resp, nil
	} else {
		handleErrorResponse(inv, w, http.StatusBadGateway, protocol.ErrUnExpectedHandlerChainResponse)
		return nil, protocol.ErrUnExpectedHandlerChainResponse
	}

}

//handleErrorResponse return proxy errors, not err from real service
func handleErrorResponse(inv *invocation.Invocation, w http.ResponseWriter, statusCode int, err error) {
	w.WriteHeader(statusCode)
	if err != nil {
		_, err := w.Write([]byte(err.Error()))
		if err != nil {
			openlogging.Error("can not write err to client: " + err.Error())
		}
	}
	RecordStatus(inv, statusCode)
}

//RecordStatus record an operation status
func RecordStatus(inv *invocation.Invocation, statusCode int) {
	LabelValues := map[string]string{metrics.LServiceName: inv.MicroServiceName, metrics.LApp: inv.RouteTags.AppID(), metrics.LVersion: inv.RouteTags.Version()}
	metrics.RecordStatus(LabelValues, statusCode)
}
func copyHeader(dst, src http.Header) {
	for k, vs := range src {
		for _, v := range vs {
			dst.Add(k, v)
		}
	}
}

func prepareRequest(req *http.Request) {
	if req.ContentLength == 0 {
		req.Body = nil
	}
	req.RequestURI = "" // client is forbidden to set RequestURI
	req.Close = false

	req.Header.Del("Connection")

}

func postProcessResponse(rsp *http.Response) {
	rsp.Header.Del("Connection")
}

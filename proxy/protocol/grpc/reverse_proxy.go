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

package grpc

import (
	"context"
	"errors"
	"io"
	"net/http"
	"time"

	"fmt"
	"github.com/apache/servicecomb-mesher/proxy/common"
	"github.com/apache/servicecomb-mesher/proxy/pkg/metrics"
	"github.com/apache/servicecomb-mesher/proxy/protocol"
	"github.com/apache/servicecomb-mesher/proxy/resolver"
	"github.com/apache/servicecomb-mesher/proxy/util"
	"github.com/go-chassis/go-chassis/client/rest"
	chassisCommon "github.com/go-chassis/go-chassis/core/common"
	"github.com/go-chassis/go-chassis/core/handler"
	"github.com/go-chassis/go-chassis/core/invocation"
	"github.com/go-chassis/go-chassis/core/lager"
	"github.com/go-chassis/go-chassis/pkg/runtime"
	"github.com/go-chassis/go-chassis/pkg/string"
	"github.com/go-chassis/go-chassis/pkg/util/tags"
	"github.com/go-mesh/openlogging"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	inv.Protocol = "grpc"
	inv.Reply = rest.NewResponse()
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
	source := stringutil.SplitFirstSep(r.RemoteAddr, ":")

	var err error
	h := make(map[string]string)
	for k := range r.Header {
		h[k] = r.Header.Get(k)
	}
	//Resolve Destination
	if r.URL.Scheme == "" {
		r.URL.Scheme = "http"
	}
	if r.URL.Host == "" {
		r.URL.Host = r.Host
	}
	serviceName, port, err := dr.Resolve(source, "", r.URL.String(), h)
	if err != nil {
		WriteErrorResponse(inv, w, r, http.StatusBadRequest, err)
		return
	}
	inv.MicroServiceName = serviceName
	if port != "" {
		h[XForwardedPort] = port
	}

	//transfer header into ctx
	inv.Ctx = context.WithValue(inv.Ctx, chassisCommon.ContextHeaderKey{}, h)
	c, err := handler.GetChain(chassisCommon.Consumer, common.ChainConsumerOutgoing)
	if err != nil {
		WriteErrorResponse(inv, w, r, http.StatusBadGateway, err)
		lager.Logger.Error("Get chain failed: " + err.Error())
		return
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
	resp, err := handleRequest(w, r, inv, invRsp)
	if err != nil {
		lager.Logger.Error("Handle request failed: " + err.Error())
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
		WriteErrorResponse(inv, w, r, http.StatusBadGateway, err)
		lager.Logger.Error("Get chain failed: " + err.Error())
		return
	}
	if err = util.SetLocalServiceAddress(inv, r.Header.Get("X-Forwarded-Port")); err != nil {
		WriteErrorResponse(inv, w, r, http.StatusBadGateway,
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
	if _, err = handleRequest(w, r, inv, invRsp); err != nil {
		lager.Logger.Error("Handle request failed: " + err.Error())
	}
}

func copyChassisResp2HttpResp(w http.ResponseWriter, resp *http.Response) {
	if resp == nil || resp.StatusCode == 0 {
		lager.Logger.Warn("response is nil or empty because of unknown reason, plz report issue")
		return
	}
	copyHeader(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)
	_, err := io.Copy(w, resp.Body)
	if err != nil {
		openlogging.Error("can not copy resp: " + err.Error())
	}
	resp.Body.Close()
}
func handleRequest(w http.ResponseWriter, r *http.Request, inv *invocation.Invocation, ir *invocation.Response) (*http.Response, error) {
	if ir != nil {
		if ir.Err != nil {
			WriteErrorResponse(inv, w, r, ir.Status, ir.Err)
			return nil, ir.Err
		}
		if inv.Endpoint == "" {
			WriteErrorResponse(inv, w, r, http.StatusBadGateway, protocol.ErrUnknown)
			return nil, protocol.ErrUnknown
		}
		if ir.Result == nil {
			if ir.Err != nil {
				WriteErrorResponse(inv, w, r, http.StatusBadGateway, ir.Err)
				return nil, ir.Err
			}
			WriteErrorResponse(inv, w, r, http.StatusBadGateway, ErrNilResponse)
			return nil, protocol.ErrUnknown
		}
		resp, ok := ir.Result.(*http.Response)
		if !ok {
			err := errors.New("invocationResponse result is not type *http.Response")
			WriteErrorResponse(inv, w, r, http.StatusBadGateway, err)
			return nil, err
		}
		//transparent proxy
		copyChassisResp2HttpResp(w, resp)

		return resp, nil
	} else {
		WriteErrorResponse(inv, w, r, http.StatusBadGateway, protocol.ErrUnExpectedHandlerChainResponse)
		return nil, protocol.ErrUnExpectedHandlerChainResponse
	}

}

//WriteErrorResponse return proxy errors, not err from real service
func WriteErrorResponse(inv *invocation.Invocation, w http.ResponseWriter, r *http.Request, statusCode int, err error) {
	stat, ok := status.FromError(err)
	if !ok {
		stat = status.New(codes.Unknown, err.Error())
	}
	openlogging.GetLogger().Errorf("grpc error: [%s]: [%s]", stat.Code().String(), stat.Message())
	w.Header().Set("Content-Type", r.Header.Get("content-type"))
	w.Header().Set("User-Agent", r.Header.Get("User-Agent"))
	w.Header().Set("Grpc-Status", fmt.Sprintf("%d", stat.Code()))
	if m := stat.Message(); m != "" {
		w.Header().Set("Grpc-Message", m)
	}
	RecordStatus(inv, int(stat.Code()))
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
	dst.Set("Trailer", "Grpc-Status, Grpc-Message")
	dst.Set("Grpc-Status", "0")
	dst.Set("Grpc-Message", "")
}

func prepareRequest(req *http.Request) {
	if req.ContentLength == 0 {
		req.Body = nil
	}
	req.RequestURI = "" // client is forbidden to set RequestURI
	req.Close = false

	req.Header.Del("Connection")

}

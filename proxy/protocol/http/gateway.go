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
	"github.com/apache/servicecomb-mesher/proxy/common"
	"github.com/apache/servicecomb-mesher/proxy/ingress"
	"github.com/go-chassis/go-chassis/client/rest"
	chassiscommon "github.com/go-chassis/go-chassis/core/common"
	"github.com/go-chassis/go-chassis/core/handler"
	"github.com/go-chassis/go-chassis/core/invocation"
	"github.com/go-mesh/openlogging"
	"net/http"
)

func handleIncomingTraffic(inv *invocation.Invocation) (*invocation.Response, error) {
	c, err := handler.GetChain(chassiscommon.Provider, common.ChainProviderIncoming)
	if err != nil {
		openlogging.Error("Get chain failed: " + err.Error())
		return nil, err
	}
	var invRsp *invocation.Response
	c.Next(inv, func(ir *invocation.Response) error {
		invRsp = ir
		if invRsp != nil {
			return invRsp.Err
		}
		return nil
	})
	return invRsp, nil
}

//HandleIngressTraffic is api gateway http handler
func HandleIngressTraffic(w http.ResponseWriter, r *http.Request) {
	inv := &invocation.Invocation{}
	inv.Reply = rest.NewResponse()
	inv.Protocol = "rest"
	inv.Args = r
	h := make(map[string]string)
	for k := range r.Header {
		h[k] = r.Header.Get(k)
	}
	inv.Ctx = chassiscommon.NewContext(h)
	invResp, err := handleIncomingTraffic(inv)
	if err != nil {
		handleErrorResponse(inv, w, http.StatusInternalServerError, err)
		return
	}
	if invResp != nil {
		if invResp.Status != 0 || invResp.Err != nil {
			handleErrorResponse(inv, w, invResp.Status, invResp.Err)
			return
		}
	}
	rule, err := ingress.DefaultFetcher.Fetch("http", r.Host, r.URL.Path, r.Header)
	if err != nil {
		handleErrorResponse(inv, w, http.StatusInternalServerError, err)
		return
	}
	inv.MicroServiceName = rule.Service.Name
	targetAPI := r.URL.Path
	if rule.Service.RedirectPath != "" {
		targetAPI = rule.Service.RedirectPath
	}
	newReq, err := http.NewRequest(r.Method, "http://"+inv.MicroServiceName+targetAPI, r.Body)
	if err != nil {
		handleErrorResponse(inv, w, http.StatusInternalServerError, err)
		return
	}
	inv.Args = newReq
	h[XForwardedPort] = rule.Service.Port.Value
	c, err := handler.GetChain(chassiscommon.Consumer, common.ChainConsumerOutgoing)
	if err != nil {
		handleErrorResponse(inv, w, http.StatusBadGateway, err)
		openlogging.Error("Get chain failed: " + err.Error())
		return
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
	resp, err := handleRequest(w, inv, invRsp)
	if err != nil {
		openlogging.Error("Handle request failed: " + err.Error())
		return
	}
	RecordStatus(inv, resp.StatusCode)
}

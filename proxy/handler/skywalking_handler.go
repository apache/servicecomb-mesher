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

package handler

import (
	"fmt"
	"github.com/apache/servicecomb-mesher/proxy/pkg/skywalking"
	"github.com/go-chassis/go-chassis/v2/core/handler"
	"github.com/go-chassis/go-chassis/v2/core/invocation"
	"github.com/go-chassis/openlog"
	"github.com/tetratelabs/go2sky"
	skycom "github.com/tetratelabs/go2sky/reporter/grpc/common"
	"strconv"
)

const (
	HTTPPrefix = "http://"
)

const (
	HTTPClientComponentID  = 2
	ServiceCombComponentID = 28
	HTTPServerComponentID  = 49
)

//SkyWalkingProviderHandler struct
type SkyWalkingProviderHandler struct {
}

//Handle is for provider
func (sp *SkyWalkingProviderHandler) Handle(chain *handler.Chain, i *invocation.Invocation, cb invocation.ResponseCallBack) {
	openlog.Debug(fmt.Sprintf("SkyWalkingProviderHandler begin. inv:%#v", *i))
	span, _, err := skywalking.CreateEntrySpan(i)
	if err != nil {
		openlog.Error(fmt.Sprintf("CreateEntrySpan error:%s", err.Error()))
	}
	chain.Next(i, func(r *invocation.Response) {
		cb(r)
		span.Tag(go2sky.TagHTTPMethod, i.Protocol)
		span.Tag(go2sky.TagURL, HTTPPrefix+i.MicroServiceName+i.URLPath)
		span.Tag(go2sky.TagStatusCode, strconv.Itoa(r.Status))
		span.SetSpanLayer(skycom.SpanLayer_Http)
		span.SetComponent(HTTPServerComponentID)
		span.End()
	})
}

//Name return provider name
func (sp *SkyWalkingProviderHandler) Name() string {
	return skywalking.SkyWalkingProvider
}

//NewSkyWalkingProvier return provider handler for SkyWalking
func NewSkyWalkingProvier() handler.Handler {
	return &SkyWalkingProviderHandler{}
}

//SkyWalkingConsumerHandler struct
type SkyWalkingConsumerHandler struct {
}

//Handle is for consumer
func (sc *SkyWalkingConsumerHandler) Handle(chain *handler.Chain, i *invocation.Invocation, cb invocation.ResponseCallBack) {
	openlog.Debug(fmt.Sprintf("SkyWalkingConsumerHandler begin:%#v", *i))
	span, ctx, err := skywalking.CreateEntrySpan(i)
	if err != nil {
		openlog.Error(fmt.Sprintf("CreateEntrySpan error:%s", err.Error()))
	}
	spanExit, err := skywalking.CreateExitSpan(ctx, i)
	if err != nil {
		openlog.Error(fmt.Sprintf("CreateExitSpan error:%s", err.Error()))
	}
	chain.Next(i, func(r *invocation.Response) {
		cb(r)
		span.Tag(go2sky.TagHTTPMethod, i.Protocol)
		span.Tag(go2sky.TagURL, HTTPPrefix+i.MicroServiceName+i.URLPath)
		span.Tag(go2sky.TagStatusCode, strconv.Itoa(r.Status))
		span.SetSpanLayer(skycom.SpanLayer_Http)
		span.SetComponent(HTTPServerComponentID)

		spanExit.Tag(go2sky.TagHTTPMethod, i.Protocol)
		spanExit.Tag(go2sky.TagURL, HTTPPrefix+i.MicroServiceName+i.URLPath)
		spanExit.Tag(go2sky.TagStatusCode, strconv.Itoa(r.Status))
		spanExit.SetSpanLayer(skycom.SpanLayer_Http)
		spanExit.SetComponent(HTTPClientComponentID)

		spanExit.End()
		span.End()
		openlog.Debug("SkyWalkingConsumerHandler end.")
	})
}

//Name return consumer name
func (sc *SkyWalkingConsumerHandler) Name() string {
	return skywalking.SkyWalkingConsumer
}

//NewSkyWalkingConsumer return consumer handler for SkyWalking
func NewSkyWalkingConsumer() handler.Handler {
	return &SkyWalkingConsumerHandler{}
}

func init() {
	err := handler.RegisterHandler(skywalking.SkyWalkingProvider, NewSkyWalkingProvier)
	if err != nil {
		openlog.Error(fmt.Sprintf("Handler [%s] register error: %s", skywalking.SkyWalkingProvider, err.Error()))
	}
	err = handler.RegisterHandler(skywalking.SkyWalkingConsumer, NewSkyWalkingConsumer)
	if err != nil {
		openlog.Error(fmt.Sprintf("Handler [%s] register error: %s", skywalking.SkyWalkingConsumer, err.Error()))
	}
}

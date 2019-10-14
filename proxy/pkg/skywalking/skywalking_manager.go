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

package skywalking

import (
	"context"
	"github.com/apache/servicecomb-mesher/proxy/config"
	gcconfig "github.com/go-chassis/go-chassis/core/config"
	"github.com/go-chassis/go-chassis/core/invocation"
	"github.com/go-mesh/openlogging"
	"github.com/tetratelabs/go2sky"
	"github.com/tetratelabs/go2sky/reporter"
)

const (
	CrossProcessProtocolV2 = "Sw6"
	SkyWalkingConsumer     = "skywalking-consumer"
	SkyWalkingProvider     = "skywalking-provider"
	SkyWalkingName         = "skywalking"
	DeafaultSWServerURI    = "127.0.0.1:11800"
)

var r go2sky.Reporter
var tracer *go2sky.Tracer

//CreateEntrySpan use tracer to create and start an entry span for incoming request
func CreateEntrySpan(i *invocation.Invocation) (go2sky.Span, context.Context, error) {
	return tracer.CreateEntrySpan(i.Ctx, i.MicroServiceName+i.URLPathFormat, func() (string, error) {
		return i.Headers()[CrossProcessProtocolV2], nil
	})
}

//CreateExitSpan use tracer to create and start an exit span for client
func CreateExitSpan(ctx context.Context, i *invocation.Invocation) (go2sky.Span, error) {
	return tracer.CreateExitSpan(ctx, i.MicroServiceName+i.URLPathFormat, i.Endpoint+i.URLPathFormat, func(header string) error {
		i.SetHeader(CrossProcessProtocolV2, header)
		return nil
	})
}

//CreateLocalSpan use tracer to create and start a span for local usage
func CreateLocalSpan(ctx context.Context, opts ...go2sky.SpanOption) (go2sky.Span, context.Context, error) {
	return tracer.CreateLocalSpan(ctx, opts...)
}

//Init skywalking manager
func Init() {
	openlogging.GetLogger().Debugf("SkyWalking manager Init begin config:%#v", config.GetConfig().ServiceComb.APM)
	var err error
	serverURI := DeafaultSWServerURI
	if config.GetConfig().ServiceComb.APM.Tracing.ServerURI != "" && config.GetConfig().ServiceComb.APM.Tracing.Enable {
		serverURI = config.GetConfig().ServiceComb.APM.Tracing.ServerURI
	}
	r, err = reporter.NewGRPCReporter(serverURI)
	if err != nil {
		openlogging.GetLogger().Errorf("NewGRPCReporter error:%s ", err.Error())
	}
	tracer, err = go2sky.NewTracer(gcconfig.MicroserviceDefinition.ServiceDescription.Name, go2sky.WithReporter(r))
	if err != nil {
		openlogging.GetLogger().Errorf("NewTracer error " + err.Error())
	}
	//tracer.WaitUntilRegister()
	openlogging.GetLogger().Debugf("SkyWalking manager Init end")
}

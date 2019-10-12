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
	"container/list"
	"github.com/apache/servicecomb-mesher/proxy/pkg/apm"
	"github.com/go-chassis/go-chassis/core/config"
	"github.com/go-chassis/go-chassis/core/invocation"
	"github.com/go-mesh/openlogging"
	"github.com/tetratelabs/go2sky"
	"github.com/tetratelabs/go2sky/reporter"
	"github.com/tetratelabs/go2sky/reporter/grpc/common"
	"strconv"
)

const (
	HttpPrefix             = "http://"
	CrossProcessProtocolV2 = "Sw6"
	Name                   = "skywalking"
)

const (
	HttpClientComponentID = 2
	HttpServerComponentID = 49
)

type SkyWalkingClient struct {
	reporter go2sky.Reporter
	tracer   *go2sky.Tracer
}

//CreateEntrySpan create entry span
func (s *SkyWalkingClient) CreateEntrySpan(i *invocation.Invocation) (interface{}, error) {
	openlogging.GetLogger().Debugf("CreateEntrySpan begin. inv:%#v", i)
	span, ctx, err := s.tracer.CreateEntrySpan(i.Ctx, i.MicroServiceName, func() (string, error) {
		return i.Headers()[CrossProcessProtocolV2], nil
	})
	if err != nil {
		openlogging.GetLogger().Errorf("CreateExitSpan error:%s", err.Error())
		return &span, err
	}
	span.Tag(go2sky.TagHTTPMethod, i.Protocol)
	span.Tag(go2sky.TagURL, i.Endpoint+i.URLPathFormat)
	span.SetSpanLayer(common.SpanLayer_Http)
	span.SetComponent(HttpServerComponentID)
	i.Ctx = ctx
	return &span, err
}

//CreateExitSpan create end span
func (s *SkyWalkingClient) CreateExitSpan(i *invocation.Invocation) (interface{}, error) {
	openlogging.GetLogger().Debugf("CreateExitSpan begin. inv:%v", i)
	span, err := s.tracer.CreateExitSpan(i.Ctx, i.MicroServiceName, i.Endpoint+i.URLPathFormat, func(header string) error {
		i.SetHeader(CrossProcessProtocolV2, header)
		return nil
	})
	if err != nil {
		openlogging.GetLogger().Errorf("CreateExitSpan error:%s", err.Error())
		return &span, err
	}
	span.Tag(go2sky.TagHTTPMethod, i.Protocol)
	span.Tag(go2sky.TagURL, HttpPrefix+i.MicroServiceName+i.URLPathFormat)
	span.SetSpanLayer(common.SpanLayer_Http)
	span.SetComponent(HttpClientComponentID)
	return &span, err
}

//EndSpan make span end and report to skywalking
func (s *SkyWalkingClient) EndSpan(sp interface{}, statusCode int) error {
	span, ok := (sp).(*go2sky.Span)
	if !ok || span == nil {
		return nil
	}
	(*span).Tag(go2sky.TagStatusCode, strconv.Itoa(statusCode))
	(*span).End()
	return nil
}

//CreateSpans create entry and exit spans for report
func (s *SkyWalkingClient) CreateSpans(i *invocation.Invocation) ([]interface{}, error) {
	openlogging.GetLogger().Debugf("CreateSpans begin. inv:%#v", i)
	var spans []interface{}
	span, ctx, err := s.tracer.CreateEntrySpan(i.Ctx, config.MicroserviceDefinition.ServiceDescription.Name, func() (string, error) {
		return i.Headers()[CrossProcessProtocolV2], nil
	})
	if err != nil {
		openlogging.GetLogger().Errorf("CreateSpans error:%s", err.Error())
		return spans, err
	}
	l := list.New()
	l.PushBack(1)
	span.Tag(go2sky.TagHTTPMethod, i.Protocol)
	span.Tag(go2sky.TagURL, HttpPrefix+i.MicroServiceName+i.URLPathFormat)
	span.SetSpanLayer(common.SpanLayer_Http)
	span.SetComponent(HttpServerComponentID)
	spans = append(spans, &span)
	spanExit, err := s.tracer.CreateExitSpan(ctx, i.MicroServiceName, i.Endpoint+i.URLPathFormat, func(header string) error {
		i.SetHeader(CrossProcessProtocolV2, header)
		return nil
	})
	if err != nil {
		openlogging.GetLogger().Errorf("CreateSpans error:%s", err.Error())
		return spans, err
	}
	spanExit.Tag(go2sky.TagHTTPMethod, i.Protocol)
	spanExit.Tag(go2sky.TagURL, i.Endpoint+i.URLPathFormat)
	spanExit.SetSpanLayer(common.SpanLayer_Http)
	spanExit.SetComponent(HttpClientComponentID)
	spans = append(spans, &spanExit)
	return spans, nil

}

//EndSpans make spans end and report to skywalking
func (s *SkyWalkingClient) EndSpans(spans []interface{}, status int) error {
	openlogging.GetLogger().Debugf("EndSpans spans:%u status:%#v", len(spans), status)
	for i := len(spans) - 1; i >= 0; i-- {
		span, ok := (spans[i]).(*go2sky.Span)
		if !ok || spans[i] == nil {
			continue
		}
		(*span).Tag(go2sky.TagStatusCode, strconv.Itoa(status))
		(*span).End()
	}
	return nil
}

//NewApmClient init report and tracer for connecting and sending messages to skywalking server
func NewApmClient(opts apm.Options) (apm.ApmClient, error) {
	var (
		err    error
		client SkyWalkingClient
	)
	client.reporter, err = reporter.NewGRPCReporter(opts.ServerUri)
	if err != nil {
		openlogging.GetLogger().Errorf("NewGRPCReporter error:%s", err.Error())
		return &client, err
	}
	client.tracer, err = go2sky.NewTracer(config.MicroserviceDefinition.ServiceDescription.Name, go2sky.WithReporter(client.reporter))
	//t.WaitUntilRegister()
	if err != nil {
		openlogging.GetLogger().Errorf("NewTracer error:%s", err.Error())
		return &client, err

	}
	openlogging.GetLogger().Debugf("NewApmClient succ. name:%s uri:%s", config.MicroserviceDefinition.ServiceDescription.Name, opts.ServerUri)
	return &client, err
}

func init() {
	apm.InstallClientPlugins(Name, NewApmClient)
}

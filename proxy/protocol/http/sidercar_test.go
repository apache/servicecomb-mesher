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
	"bytes"
	"github.com/apache/servicecomb-mesher/proxy/cmd"
	"github.com/go-chassis/go-chassis/client/rest"
	"github.com/go-chassis/go-chassis/core/common"
	"github.com/go-chassis/go-chassis/core/handler"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/apache/servicecomb-mesher/proxy/pkg/metrics"
	"github.com/go-chassis/go-chassis/core/lager"
)

func init() {
	lager.Init(&lager.Options{
		LoggerLevel:   "INFO",
		RollingPolicy: "size",
	})

	cmd.Init()

	metrics.Init()

}

func TestLocalRequestHandler(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(LocalRequestHandler))
	api := svr.URL
	rsp, err := http.Get(api)
	if err != nil {
		return
	}
	defer rsp.Body.Close()
}

func TestRemoteRequestHandler(t *testing.T) {
	handler.CreateChains(
		common.Provider, map[string]string{
			"incoming": strings.Join([]string{}, ","),
		},
	)

	handler.CreateChains(
		common.Consumer, map[string]string{
			"outgoing": strings.Join([]string{}, ","),
		},
	)

	svr := httptest.NewServer(http.HandlerFunc(RemoteRequestHandler))
	api := svr.URL
	rsp, err := http.Get(api)
	if err != nil {
		return
	}
	defer rsp.Body.Close()
}

func TestCopyChassisResp2HttpResp(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := rest.NewResponse()
		resp.StatusCode = 200
		b, _ := ioutil.ReadAll(r.Body)
		resp.Body = ioutil.NopCloser(bytes.NewBuffer(b))
		copyChassisResp2HttpResp(w, resp)

	}))
	api := svr.URL
	rsp, err := http.Get(api)
	if err != nil {
		return
	}
	defer rsp.Body.Close()
}

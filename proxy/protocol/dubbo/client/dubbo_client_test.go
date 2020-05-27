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

package dubboclient

import (
	"github.com/apache/servicecomb-mesher/proxy/protocol/dubbo/dubbo"
	"github.com/go-chassis/go-chassis/core/lager"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func init() {
	lager.Init(&lager.Options{
		LoggerLevel:   "INFO",
		RollingPolicy: "size",
	})
}

func TestClientMgr_GetClient(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	addr := ts.URL
	u, _ := url.Parse(ts.URL)
	addr = u.Host
	clientMgr := NewClientMgr()
	// case timeout=0
	c, err := clientMgr.GetClient(addr, 0)
	assert.NoError(t, err)
	assert.NotNil(t, c)
	c.close()

	c, err = clientMgr.GetClient(addr, 0)
	assert.NoError(t, err)
	assert.NotNil(t, c)

	req := dubbo.NewDubboRequest()
	c.Send(req)

	// case RspCallBack
	resp := &dubbo.DubboRsp{}
	resp.Init()
	resp.SetStatus(dubbo.ServerError)
	c.RspCallBack(resp)

	// case get addr
	c.GetAddr()

	// case net error
	ts.Close()
	c, err = clientMgr.GetClient(addr, 0)
	assert.Error(t, err)
	assert.Nil(t, c)

}

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

package chassisclient

import (
	"context"
	"fmt"
	mesherCommon "github.com/apache/servicecomb-mesher/proxy/common"
	dubboClient "github.com/apache/servicecomb-mesher/proxy/protocol/dubbo/client"
	"github.com/apache/servicecomb-mesher/proxy/protocol/dubbo/dubbo"
	dubboproxy "github.com/apache/servicecomb-mesher/proxy/protocol/dubbo/proxy"
	util "github.com/apache/servicecomb-mesher/proxy/protocol/dubbo/utils"
	"github.com/go-chassis/go-chassis/core/client"
	"github.com/go-chassis/go-chassis/core/invocation"
	"github.com/go-chassis/go-chassis/core/lager"
	"github.com/stretchr/testify/assert"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
	"time"
)

func init() {
	lager.Init(&lager.Options{LoggerLevel: "INFO", RollingPolicy: "size"})
}

func TestDubboChassisClient(t *testing.T) {
	addr := "127.0.0.1:31011"
	tcpAddr, _ := net.ResolveTCPAddr("tcp", addr)
	l, _ := net.ListenTCP("tcp", tcpAddr)
	writeError := false
	go func(l *net.TCPListener, writeError bool) {
		conn, _ := l.AcceptTCP()
		for {
			buf := make([]byte, dubbo.HeaderLength)
			_, err := conn.Read(buf)
			if err != nil {
				continue
			}
			req := new(dubbo.Request)
			bodyLen := 0
			coder := dubbo.DubboCodec{}
			ret := coder.DecodeDubboReqHead(req, buf, &bodyLen)
			fmt.Println("ret: ", ret)
			var buffer util.WriteBuffer
			buffer.Init(0)
			rsp := &dubbo.DubboRsp{}
			if !writeError {
				rsp.SetStatus(dubbo.Ok)
			}
			coder.EncodeDubboRsp(rsp, &buffer)

			hf := buffer.GetValidData()
			conn.Write(hf)

			// case header[0] != MagicHigh
			if writeError {
				hf[0] = 0
				conn.Write(hf)
			}
		}
	}(l, writeError)

	c, err := NewDubboChassisClient(client.Options{
		Service:  "dubbotest",
		PoolSize: 10,
		Timeout:  time.Second * 10,
		Endpoint: "127.0.0.1:23101",
	})
	assert.NoError(t, err)

	err = c.Close()
	assert.NoError(t, err)
	assert.Equal(t, "highway_client", c.String())

	c.GetOptions()

	c.ReloadConfigs(client.Options{
		Service:  "dubbotest",
		PoolSize: 10,
		Timeout:  time.Second * 10,
		Endpoint: "127.0.0.1:23101"})

	inv := &invocation.Invocation{}
	inv.Args = &dubbo.Request{}
	rsp := &dubboClient.WrapResponse{}

	dubboproxy.DubboListenAddr = addr
	endPoint := addr
	os.Setenv(mesherCommon.EnvSpecificAddr, addr)
	// case endPoint==""
	err = c.Call(context.Background(), "", inv, rsp)
	assert.Error(t, err)

	// case endPoint error
	err = c.Call(context.Background(), "127.0.0.1:23101", inv, rsp)
	assert.Error(t, err)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// case service error
	u, _ := url.Parse(ts.URL)
	err = c.Call(context.Background(), u.Host, inv, rsp)
	assert.Error(t, err)

	// case endPoint == dubboproxy.DubboListenAddr
	err = c.Call(context.Background(), endPoint, inv, rsp)
	assert.NoError(t, err)

	// writeError == true
	writeError = true
	c.Call(context.Background(), endPoint, inv, rsp)

	// writeError == false
	writeError = false
	c.Call(context.Background(), endPoint, inv, rsp)

}

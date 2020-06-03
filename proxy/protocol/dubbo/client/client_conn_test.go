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
	"fmt"
	"github.com/apache/servicecomb-mesher/proxy/protocol/dubbo/dubbo"
	util "github.com/apache/servicecomb-mesher/proxy/protocol/dubbo/utils"
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
	"time"
)

func TestClientConn(t *testing.T) {
	addr := "127.0.0.1:32100"
	tcpAddr, _ := net.ResolveTCPAddr("tcp", addr)
	l, _ := net.ListenTCP("tcp", tcpAddr)
	go func(l *net.TCPListener) {
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
			coder.EncodeDubboRsp(rsp, &buffer)

			hf := buffer.GetValidData()
			conn.Write(hf)

			// case header[0] != MagicHigh
			hf[0] = 0
			conn.Write(hf)
		}
	}(l)

	c, errDial := net.DialTimeout("tcp", addr, time.Second*5)
	assert.NoError(t, errDial)
	conn, ok := c.(*net.TCPConn)
	assert.Equal(t, true, ok)

	connClinet := NewDubboClientConnetction(conn, NewDubboClient(addr, nil, time.Second*5), nil)
	go func(c *DubboClientConnection) {
		t := time.NewTimer(time.Second)
		for range t.C {
			c.SendMsg(dubbo.NewDubboRequest())
		}
	}(connClinet)

	// Case conn open
	connClinet.Open()

	// Case conn closed
	select {
	case <-time.After(time.Second * 3):
		conn.Close()
		connClinet.SendMsg(dubbo.NewDubboRequest())
	}
	// case close
	connClinet.Close()
	assert.Equal(t, true, connClinet.Closed())
	connClinet.SendMsg(dubbo.NewDubboRequest())

	select {
	case <-time.After(time.Second * 5):
		connClinet.Close()
	}

	// case conn closed
	NewDubboClientConnetction(conn, NewDubboClient(addr, nil, time.Second*5), nil)

}

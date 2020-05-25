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
	"github.com/apache/servicecomb-mesher/proxy/protocol/dubbo/utils"
	"github.com/go-chassis/go-chassis/core/lager"
	"net"
	"sync"
)

//SndTask is a struct
type SndTask struct{}

//Svc is a method
func (this SndTask) Svc(arg interface{}) interface{} {
	dubboConn := arg.(*DubboClientConnection)
	dubboConn.MsgSndLoop()
	return nil
}

//RecvTask is a struct
type RecvTask struct {
}

//Svc is a method
func (this RecvTask) Svc(arg interface{}) interface{} {
	dubboConn := arg.(*DubboClientConnection)
	dubboConn.MsgRecvLoop()
	return nil
}

//ProcessTask is a struct
type ProcessTask struct {
	conn    *DubboClientConnection
	rsp     *dubbo.DubboRsp
	bufBody []byte
}

//Svc is a method
func (this ProcessTask) Svc(arg interface{}) interface{} {
	if this.conn != nil {
		this.conn.ProcessBody(this.rsp, this.bufBody)
	}
	return nil
}

//DubboClientConnection is a struct which has attributes for dubbo protocol connection
type DubboClientConnection struct {
	msgque     *util.MsgQueue
	remoteAddr string
	conn       *net.TCPConn
	codec      dubbo.DubboCodec
	client     *DubboClient
	mtx        sync.Mutex
	routineMgr *util.RoutineManager
	closed     bool
}

//NewDubboClientConnetction is a function which create new dubbo client connection
func NewDubboClientConnetction(conn *net.TCPConn, client *DubboClient, routineMgr *util.RoutineManager) *DubboClientConnection {
	tmp := new(DubboClientConnection)
	err := conn.SetKeepAlive(true)
	if err != nil {
		lager.Logger.Error("TCPConn SetKeepAlive error:" + err.Error())
	}
	tmp.conn = conn
	tmp.codec = dubbo.DubboCodec{}
	tmp.client = client
	tmp.msgque = util.NewMsgQueue()
	tmp.closed = false
	if routineMgr == nil {
		tmp.routineMgr = util.NewRoutineManager()
	}
	return tmp
}

//Open is a method which open connection
func (this *DubboClientConnection) Open() {
	this.routineMgr.Spawn(SndTask{}, this, fmt.Sprintf("client Snd-%s->%s", this.conn.LocalAddr().String(), this.conn.RemoteAddr().String()))
	this.routineMgr.Spawn(RecvTask{}, this, fmt.Sprintf("client Recv-%s->%s", this.conn.LocalAddr().String(), this.conn.RemoteAddr().String()))
}

//Close is a method which closes connection
func (this *DubboClientConnection) Close() {
	this.mtx.Lock()
	defer this.mtx.Unlock()
	if this.closed {
		return
	}
	this.closed = true
	this.msgque.Deavtive()
	err := this.conn.Close()
	if err != nil {
		lager.Logger.Error("Dubbo client connection close error:" + err.Error())
	}
}

//MsgRecvLoop is a method which receives message
func (this *DubboClientConnection) MsgRecvLoop() {
	//通知处理应答消息
	for {
		//先处理消息头
		buf := make([]byte, dubbo.HeaderLength)
		size, err := this.conn.Read(buf)
		if err != nil {
			if nerr, ok := err.(net.Error); ok && nerr.Timeout() {
				lager.Logger.Error("client Recv head time err:" + err.Error())
				//time.Sleep(time.Second * 3)
				continue
			}
			lager.Logger.Error("client Recv head err:" + err.Error())
			break
		}

		if size < dubbo.HeaderLength {
			continue
		}
		rsp := new(dubbo.DubboRsp)
		bodyLen := 0
		ret := this.codec.DecodeDubboRsqHead(rsp, buf, &bodyLen)
		if ret != dubbo.Success {
			lager.Logger.Info("Recv DecodeDubboRsqHead failed")
			continue
		}
		body := make([]byte, bodyLen)
		count := 0
		for {
			redBuff := body[count:]
			size, err = this.conn.Read(redBuff)
			if err != nil {
				if nerr, ok := err.(net.Error); ok && nerr.Timeout() {
					continue
				}
				//通知关闭连接
				lager.Logger.Error("Recv client body err:" + err.Error())
				goto exitloop
			}
			count += size
			if count == bodyLen {
				break
			}
		}
		this.routineMgr.Spawn(ProcessTask{this, rsp, body}, nil, fmt.Sprintf("Client ProcessTask-%d", rsp.GetID()))
	}
exitloop:
	this.Close()
}

//ProcessBody is a method which process body data
func (this *DubboClientConnection) ProcessBody(rsp *dubbo.DubboRsp, bufBody []byte) {
	var buffer util.ReadBuffer
	buffer.SetBuffer(bufBody)
	this.codec.DecodeDubboRspBody(&buffer, rsp)
	this.HandleMsg(rsp)
}

//HandleMsg is a method which returns message from dubbo response
func (this *DubboClientConnection) HandleMsg(rsp *dubbo.DubboRsp) {
	this.client.RspCallBack(rsp)
}

//SendMsg is a method which send a request
func (this *DubboClientConnection) SendMsg(req *dubbo.Request) {
	//这里发送Rest请求以及收发送应答
	err := this.msgque.Enqueue(req)
	if err != nil {
		lager.Logger.Error("Msg Enqueue:" + err.Error())
	}
}

//MsgSndLoop is a method which send data
func (this *DubboClientConnection) MsgSndLoop() {
	for {
		msg, err := this.msgque.Dequeue()
		if err != nil {
			lager.Logger.Error("MsgSndLoop Dequeue:" + err.Error())
			break
		}
		var buffer util.WriteBuffer
		buffer.Init(0)
		this.codec.EncodeDubboReq(msg.(*dubbo.Request), &buffer)
		_, err = this.conn.Write(buffer.GetValidData())
		if err != nil {
			lager.Logger.Error("Send exception:" + err.Error())
			break
		}
	}
	this.Close()
}

//Closed is a method which checks connnection is closed or not
func (this *DubboClientConnection) Closed() bool {
	return this.closed
}

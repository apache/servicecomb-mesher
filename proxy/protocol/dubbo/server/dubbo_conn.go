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

package server

import (
	"fmt"
	"github.com/apache/servicecomb-mesher/proxy/protocol/dubbo/dubbo"
	"github.com/apache/servicecomb-mesher/proxy/protocol/dubbo/proxy"
	"github.com/apache/servicecomb-mesher/proxy/protocol/dubbo/utils"
	"github.com/go-chassis/go-chassis/core/lager"
	"net"
	"sync"
)

//SndTask is a struct
type SndTask struct{}

//Svc is a method
func (this SndTask) Svc(arg interface{}) interface{} {
	dubboConn := arg.(*DubboConnection)
	dubboConn.MsgSndLoop()
	return nil
}

//RecvTask is a struct
type RecvTask struct {
}

//Svc is a method
func (this RecvTask) Svc(arg interface{}) interface{} {
	dubboConn := arg.(*DubboConnection)
	dubboConn.MsgRecvLoop()
	return nil
}

//ProcessTask is a struct
type ProcessTask struct {
	conn    *DubboConnection
	req     *dubbo.Request
	bufBody []byte
}

//Svc is a method
func (this ProcessTask) Svc(arg interface{}) interface{} {
	if this.conn != nil {
		this.conn.ProcessBody(this.req, this.bufBody)
	}
	return nil
}

//DubboConnection is a struct which has attributes for dubbo connection
type DubboConnection struct {
	msgque     *util.MsgQueue
	remoteAddr string
	conn       *net.TCPConn
	codec      dubbo.DubboCodec
	mtx        sync.Mutex
	routineMgr *util.RoutineManager
	closed     bool
}

//NewDubboConnetction is a function to create new dubbo connection
func NewDubboConnetction(conn *net.TCPConn, routineMgr *util.RoutineManager) *DubboConnection {
	tmp := new(DubboConnection)
	tmp.conn = conn
	tmp.codec = dubbo.DubboCodec{}
	tmp.msgque = util.NewMsgQueue()
	tmp.remoteAddr = conn.RemoteAddr().String()
	tmp.closed = false
	if routineMgr == nil {
		tmp.routineMgr = util.NewRoutineManager()
	}
	return tmp
}

//Open is a function to open a connection
func (this *DubboConnection) Open() {
	this.routineMgr.Spawn(SndTask{}, this, fmt.Sprintf("Snd-%s->%s", this.conn.LocalAddr().String(), this.conn.RemoteAddr().String()))
	this.routineMgr.Spawn(RecvTask{}, this, fmt.Sprintf("Recv-%s->%s", this.conn.LocalAddr().String(), this.conn.RemoteAddr().String()))
}

//Close is a function to close a connection
func (this *DubboConnection) Close() {
	this.mtx.Lock()
	defer this.mtx.Unlock()
	if this.closed {
		return
	}
	this.closed = true
	this.msgque.Deavtive()
	this.conn.Close()
}

//MsgRecvLoop is a method receive data
func (this *DubboConnection) MsgRecvLoop() {
	//通知处理应答消息
	for {
		//先处理消息头
		buf := make([]byte, dubbo.HeaderLength)
		size, err := this.conn.Read(buf)
		if err != nil {
			if nerr, ok := err.(net.Error); ok && nerr.Timeout() {
				lager.Logger.Error("Dubbo server Recv head: " + err.Error())
				continue
			}
			lager.Logger.Error("Dubbo server Recv head: " + err.Error())
			break
		}

		if size < dubbo.HeaderLength {
			lager.Logger.Info("Invalid msg head")
			continue
		}
		req := new(dubbo.Request)
		bodyLen := 0
		ret := this.codec.DecodeDubboReqHead(req, buf, &bodyLen)
		if ret != dubbo.Success {
			lager.Logger.Info("Invalid msg head")
			continue
		}
		body := make([]byte, bodyLen)
		count := 0
		for {
			redBuff := body[count:]
			size, err = this.conn.Read(redBuff)

			if err != nil {
				//通知关闭连接
				lager.Logger.Error("Recv: " + err.Error())
				goto exitloop
			}
			count += size
			if count == bodyLen {
				break
			}
		}
		this.routineMgr.Spawn(ProcessTask{this, req, body}, nil, fmt.Sprintf("ProcessTask-%d", req.GetMsgID()))
	}
exitloop:
	this.Close()
}

//ProcessBody is a method to process the body of response
func (this *DubboConnection) ProcessBody(req *dubbo.Request, bufBody []byte) {
	var buffer util.ReadBuffer
	buffer.SetBuffer(bufBody)
	this.codec.DecodeDubboReqBody(req, &buffer)
	this.HandleMsg(req)
}

//HandleMsg is a method
func (this *DubboConnection) HandleMsg(req *dubbo.Request) {
	//这里发送Rest请求以及收发送应答
	ctx := &dubbo.InvokeContext{req, &dubbo.DubboRsp{}, nil, "", this.remoteAddr}
	ctx.Rsp.Init()
	ctx.Rsp.SetID(req.GetMsgID())
	if req.IsHeartbeat() {
		ctx.Rsp.SetValue(nil)
		ctx.Rsp.SetEvent(true)
		ctx.Rsp.SetStatus(dubbo.Ok)
	} else {
		//这里重新分配MSGID
		srcMsgID := ctx.Req.GetMsgID()
		dstMsgID := dubbo.GenerateMsgID()
		//lager.Logger.Info(fmt.Sprintf("dubbo2dubbo srcMsgID=%d, newMsgID=%d", srcMsgID, dstMsgID))
		ctx.Req.SetMsgID(dstMsgID)
		err := dubboproxy.Handle(ctx)
		if err != nil {
			ctx.Rsp.SetErrorMsg(err.Error())
			lager.Logger.Error("request: " + err.Error())
			ctx.Rsp.SetStatus(dubbo.ServerError)
		}
		ctx.Req.SetMsgID(srcMsgID)
		ctx.Rsp.SetID(srcMsgID)
	}
	if req.IsTwoWay() {
		this.msgque.Enqueue(ctx.Rsp)

	}
}

//MsgSndLoop is a method to send data
func (this *DubboConnection) MsgSndLoop() {
	for {
		msg, err := this.msgque.Dequeue()
		if err != nil {
			lager.Logger.Error("MsgSndLoop Dequeue: " + err.Error())
			break
		}
		var buffer util.WriteBuffer
		buffer.Init(0)
		this.codec.EncodeDubboRsp(msg.(*dubbo.DubboRsp), &buffer)
		bs := buffer.GetValidData()
		_, err = this.conn.Write(bs /*buffer.GetValidData()*/)
		if err != nil {
			lager.Logger.Error("Send exception: " + err.Error())
			break
		}
	}
	this.Close()
}

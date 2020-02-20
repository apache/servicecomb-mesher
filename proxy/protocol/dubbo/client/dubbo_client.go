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
	"time"
)

//DubboClient is a struct which has attributes for dubboClient
type DubboClient struct {
	addr          string
	mtx           sync.Mutex
	mapMutex      sync.Mutex
	msgWaitRspMap map[int64]*RespondResult
	conn          *DubboClientConnection
	closed        bool
	routeMgr      *util.RoutineManager
	Timeout       time.Duration
}

//WrapResponse is a struct
type WrapResponse struct {
	Resp *dubbo.DubboRsp
}

//CachedClients is a variable which stores
var CachedClients *ClientMgr

func init() {
	CachedClients = NewClientMgr()
}

//RespondResult is a struct which has attribute for dubbo response
type RespondResult struct {
	Rsp  *dubbo.DubboRsp
	Wait *chan int
}

//ClientMgr is a struct which has attributes for managing client
type ClientMgr struct {
	mapMutex sync.Mutex
	clients  map[string]*DubboClient
}

//NewClientMgr is a function which creates new clientmanager and returns it
func NewClientMgr() *ClientMgr {
	tmp := new(ClientMgr)
	tmp.clients = make(map[string]*DubboClient)
	return tmp
}

//GetClient is a function which returns the particular client for that address
func (this *ClientMgr) GetClient(addr string, timeout time.Duration) (*DubboClient, error) {
	this.mapMutex.Lock()
	defer this.mapMutex.Unlock()
	if tmp, ok := this.clients[addr]; ok {
		if timeout <= 0 {
			timeout = 30 * time.Second
		}
		if tmp.Timeout != timeout {
			tmp.Timeout = timeout
			this.clients[addr] = tmp
		}
		if !tmp.Closed() {
			lager.Logger.Info("GetClient from cached addr:" + addr)
			return tmp, nil
		} else {
			err := tmp.ReOpen()
			lager.Logger.Info("GetClient repopen addr:" + addr)
			if err != nil {
				delete(this.clients, addr)
				return nil, err
			} else {
				return tmp, nil
			}
		}
	}
	lager.Logger.Info("GetClient from new open addr:" + addr)
	tmp := NewDubboClient(addr, nil, timeout)
	err := tmp.Open()
	if err != nil {
		return nil, err
	} else {
		this.clients[addr] = tmp
		return tmp, nil
	}
}

//NewDubboClient is a function which creates new dubbo client for given value
func NewDubboClient(addr string, routeMgr *util.RoutineManager, timeout time.Duration) *DubboClient {
	tmp := &DubboClient{}
	tmp.addr = addr
	tmp.Timeout = timeout
	tmp.conn = nil
	tmp.closed = true
	tmp.msgWaitRspMap = make(map[int64]*RespondResult)
	if routeMgr == nil {
		tmp.routeMgr = util.NewRoutineManager()
	}
	return tmp
}

//GetAddr is a method which returns address of particular client
func (this *DubboClient) GetAddr() string {
	return this.addr
}

//ReOpen is a method which reopens connection
func (this *DubboClient) ReOpen() error {
	this.mtx.Lock()
	defer this.mtx.Unlock()
	this.close()
	return this.open()
}

//Open is a method which opens a connection
func (this *DubboClient) Open() error {
	this.mtx.Lock()
	defer this.mtx.Unlock()
	return this.open()
}

func (this *DubboClient) open() error {
	c, errDial := net.DialTimeout("tcp", this.addr, this.Timeout)
	if errDial != nil {
		lager.Logger.Errorf("the addr: %s %s ", this.addr, errDial)
		return errDial
	}
	conn, ok := c.(*net.TCPConn)
	if !ok {
		return fmt.Errorf("not TCPConn type")
	}
	this.conn = NewDubboClientConnetction(conn, this, nil)
	this.conn.Open()
	this.closed = false
	return nil
}

//Close is a method which closes a connection
func (this *DubboClient) Close() {
	this.mtx.Lock()
	defer this.mtx.Unlock()
	this.close()
	this.routeMgr.Done()
	this.routeMgr.Wait()
}

func (this *DubboClient) close() {
	if this.closed {
		return
	}
	this.closed = true
	this.mapMutex.Lock()
	for _, v := range this.msgWaitRspMap {
		*v.Wait <- 1
	}
	this.msgWaitRspMap = make(map[int64]*RespondResult) //清空map
	this.mapMutex.Unlock()
	this.conn.Close()
}

//AddWaitMsg is a method which adds wait message in the response
func (this *DubboClient) AddWaitMsg(msgID int64, result *RespondResult) {
	this.mapMutex.Lock()
	if this.msgWaitRspMap != nil {
		this.msgWaitRspMap[msgID] = result
	}
	this.mapMutex.Unlock()
}

//RemoveWaitMsg is a method which delete waiting message
func (this *DubboClient) RemoveWaitMsg(msgID int64) {
	this.mapMutex.Lock()
	if this.msgWaitRspMap != nil {
		delete(this.msgWaitRspMap, msgID)
	}
	this.mapMutex.Unlock()
}

//Svc is a method
func (this *DubboClient) Svc(agr interface{}) interface{} {
	this.conn.SendMsg(agr.(*dubbo.Request))
	return nil
}

//Send is a method which send request from dubbo client
func (this *DubboClient) Send(dubboReq *dubbo.Request) (*dubbo.DubboRsp, error) {
	this.mapMutex.Lock()
	if this.closed {
		if err := this.open(); err != nil {
			return nil, err
		}
	}
	this.mapMutex.Unlock()
	wait := make(chan int)
	result := &RespondResult{nil, &wait}
	msgID := dubboReq.GetMsgID()
	this.AddWaitMsg(msgID, result)

	this.routeMgr.Spawn(this, dubboReq, fmt.Sprintf("SndMsgID-%d", dubboReq.GetMsgID()))
	var timeout = false
	select {
	case <-wait:
		timeout = false
	case <-time.After(this.Timeout):
		timeout = true
	}
	if this.closed {
		lager.Logger.Info("Client been closed.")
		return nil, &util.BaseError{"Client been closed."}
	}
	this.RemoveWaitMsg(msgID)
	if timeout {
		dubboReq.SetBroken(true)
		lager.Logger.Info("Client send timeout.")
		return nil, &util.BaseError{"timeout"}
	} else {
		return result.Rsp, nil
	}
}

//RspCallBack is a method
func (this *DubboClient) RspCallBack(rsp *dubbo.DubboRsp) {
	msgID := rsp.GetID()
	var result *RespondResult
	this.mapMutex.Lock()
	defer this.mapMutex.Unlock()
	if this.msgWaitRspMap == nil {
		return
	}
	if _, ok := this.msgWaitRspMap[msgID]; ok {
		result = this.msgWaitRspMap[msgID]
		result.Rsp = rsp
		*result.Wait <- 1
	}
}

//Closed is a method which checks whether connection has been closed or not
func (this *DubboClient) Closed() bool {
	this.mtx.Lock()
	defer this.mtx.Unlock()
	if this.conn.Closed() {
		this.close()
	}
	return this.closed
}

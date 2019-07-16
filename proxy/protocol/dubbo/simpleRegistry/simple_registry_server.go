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

package simpleregistry

import (
	"github.com/apache/servicecomb-mesher/proxy/protocol/dubbo/dubbo"
	"github.com/apache/servicecomb-mesher/proxy/protocol/dubbo/proxy"
	"github.com/apache/servicecomb-mesher/proxy/protocol/dubbo/utils"
	"github.com/go-chassis/go-chassis/core/lager"
	"github.com/go-chassis/go-chassis/core/server"
	"net"
	"sync"
)

const (
	NAME = "dubboSimpleRegistry"
)

func init() {
	server.InstallPlugin(NAME, newServer)
}

func newServer(opts server.Options) server.ProtocolServer {

	return &SimDubboRegistryServer{
		opts: opts,
	}
}

//SimDubboRegistryServer is a struct
type SimDubboRegistryServer struct {
	opts       server.Options
	mux        sync.RWMutex
	exit       chan chan error
	routineMgr *util.RoutineManager
}

func (d *SimDubboRegistryServer) String() string {
	return NAME
}

//Init is a method which initialized server config
func (d *SimDubboRegistryServer) Init(opts ...server.Options) error {
	lager.Logger.Info("Dubbo Simple Registry server init.")
	return nil
}

//Register is a method to register schema to that server
func (d *SimDubboRegistryServer) Register(schema interface{}, options ...server.RegisterOption) (string, error) {
	return "", nil
}

//Stop is a method to stop the server
func (d *SimDubboRegistryServer) Stop() error {
	return nil
}

//Start is a method to start the server
func (d *SimDubboRegistryServer) Start() error {
	err := d.Init()
	if err != nil {
		return err
	}
	host, _, err := net.SplitHostPort(d.opts.Address)
	if err != nil {
		return err
	}
	ip := net.ParseIP(host)
	if ip == nil {
		return &util.BaseError{ErrMsg: "invalid host"}
	}
	tcpAddr, err := net.ResolveTCPAddr("tcp", d.opts.Address)
	if err != nil {
		lager.Logger.Error("ResolveTCPAddr err: " + err.Error())
		return err
	}
	l, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		lager.Logger.Error("listening falied, reason: " + err.Error())
		return err
	}
	go d.AcceptLoop(l)
	return nil
}

//AcceptLoop is a method to receive data in loop
func (d *SimDubboRegistryServer) AcceptLoop(l *net.TCPListener) {
	for {
		for {
			conn, err := l.Accept()
			if err != nil {
				lager.Logger.Error("tcp conn error: " + err.Error())
				continue
			}

			lager.Logger.Debug("Received message")

			go handleConn(conn)
		}

	}
	defer l.Close()
}

func handleConn(conn net.Conn) {
	defer conn.Close()
	var buf []byte
	buf = make([]byte, 1024)
	for {
		_, err := conn.Read(buf)
		if err != nil {
			return
		}
		req := &dubbo.Request{}
		codec := dubbo.DubboCodec{}

		var bodyLen int
		codec.DecodeDubboReqHead(req, buf[0:dubbo.HeaderLength], &bodyLen)

		SendVoidRespond(conn, req)
	}
}

//SendVoidRespond is a method to send void respose
func SendVoidRespond(conn net.Conn, req *dubbo.Request) {
	dubboproxy.IsProvider = true
	var rsp dubbo.DubboRsp
	var wBuf util.WriteBuffer
	wBuf = util.WriteBuffer{}
	wBuf.Init(1024)
	rsp = dubbo.DubboRsp{}
	rsp.Init()
	rsp.SetEvent(req.IsEvent())
	rsp.SetID(req.GetMsgID())
	rsp.SetValue(nil)
	codec := dubbo.DubboCodec{}
	codec.EncodeDubboRsp(&rsp, &wBuf)
	conn.Write(wBuf.GetValidData())
}

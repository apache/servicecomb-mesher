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
	"github.com/go-chassis/go-chassis/core/config"
	"github.com/go-chassis/go-chassis/core/config/schema"
	"github.com/go-mesh/openlogging"
	"gopkg.in/yaml.v2"
	"net"
	"sync"
	"time"

	"github.com/apache/servicecomb-mesher/proxy/protocol/dubbo/proxy"
	"github.com/apache/servicecomb-mesher/proxy/protocol/dubbo/utils"
	"github.com/go-chassis/go-chassis/core/lager"
	"github.com/go-chassis/go-chassis/core/server"
)

const (
	NAME = "dubbo"
)

//ConnectionMgr -------连接管理
type ConnectionMgr struct {
	conns map[int]*DubboConnection
	count int
}

//NewConnectMgr is a function which new connection manager and returns it
func NewConnectMgr() *ConnectionMgr {
	tmp := new(ConnectionMgr)
	tmp.count = 0
	tmp.conns = make(map[int]*DubboConnection)
	return tmp
}

//GetConnection is a method to get connection
func (this *ConnectionMgr) GetConnection(conn *net.TCPConn) *DubboConnection {
	dubbConn := NewDubboConnetction(conn, nil)
	key := this.count
	this.conns[key] = dubbConn
	this.count++
	return dubbConn
}

//DeactiveAllConn is a function to close all connection
func (this *ConnectionMgr) DeactiveAllConn() {
	for _, v := range this.conns {
		v.Close()
	}
}

func init() {
	server.InstallPlugin(NAME, newServer)
}

func newServer(opts server.Options) server.ProtocolServer {

	return &DubboServer{
		opts:       opts,
		routineMgr: util.NewRoutineManager(),
	}
}

//-----------------------dubbo server---------------------------

//DubboServer is a struct
type DubboServer struct {
	connMgr    *ConnectionMgr
	opts       server.Options
	mux        sync.RWMutex
	exit       chan chan error
	routineMgr *util.RoutineManager
}

func (d *DubboServer) String() string {
	return NAME
}

//Init is a method to initialize the server
func (d *DubboServer) Init() error {
	initSchema()
	d.connMgr = NewConnectMgr()
	lager.Logger.Info("Dubbo server init success.")
	return nil
}

//Register is a method to register the schema to the server
func (d *DubboServer) Register(schema interface{}, options ...server.RegisterOption) (string, error) {
	return "", nil
}

//Stop is a method to disconnect all connection
func (d *DubboServer) Stop() error {
	d.connMgr.DeactiveAllConn()
	d.routineMgr.Done()
	return nil
}

//Start is a method to start server
func (d *DubboServer) Start() error {
	err := d.Init()
	if err != nil {
		return err
	}
	dubboproxy.DubboListenAddr = d.opts.Address
	host, _, err := net.SplitHostPort(d.opts.Address)
	if err != nil {
		return err
	}
	ip := net.ParseIP(host)
	if ip == nil {
		return &util.BaseError{"invalid host"}
	}
	tcpAddr, err := net.ResolveTCPAddr("tcp", d.opts.Address)
	if err != nil {
		lager.Logger.Error("ResolveTCPAddr err: " + err.Error())
		return err
	}
	l, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		lager.Logger.Error("listening failed, reason: " + err.Error())
		return err
	}
	d.routineMgr.Spawn(d, l, "Acceptloop")
	return nil
}

//Svc is a method
func (d *DubboServer) Svc(arg interface{}) interface{} {
	d.AcceptLoop(arg.(*net.TCPListener))
	return nil
}

//AcceptLoop is a method
func (d *DubboServer) AcceptLoop(l *net.TCPListener) {
	for {
		conn, err := l.AcceptTCP()
		if err != nil {
			select {
			case <-time.After(time.Second * 3):
				lager.Logger.Info("Sleep three second")
			}
		}
		dubbConn := d.connMgr.GetConnection(conn)
		dubbConn.Open()
	}

	defer l.Close()
}

// initSchema is a method to ini the schema ids
func initSchema() {
	m := make(map[string]string, 0)
	service := config.MicroserviceDefinition
	if len(service.ServiceDescription.Schemas) == 0 {
		return
	}

	for _, inter := range service.ServiceDescription.Schemas {
		if len(inter) == 0 {
			openlogging.GetLogger().Warnf("interfaces is empty")
			break
		}
		schemaContent := struct {
			Swagger string            `yaml:"swagger"`
			Info    map[string]string `yaml:"info"`
		}{
			Swagger: "2.0",
			Info: map[string]string{
				"version":          "1.0.0",
				"title":            fmt.Sprintf("swagger definition for %s", inter),
				"x-java-interface": inter,
			},
		}

		b, err := yaml.Marshal(&schemaContent)
		if err != nil {
			break
		}

		m[inter] = string(b)
	}

	err := schema.SetSchemaInfoByMap(m)
	if err != nil {
		openlogging.Error("Set schemaInfo failed: " + err.Error())
	}
}

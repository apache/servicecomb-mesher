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
	"errors"
	"github.com/apache/servicecomb-mesher/proxy/common"
	"github.com/apache/servicecomb-mesher/proxy/pkg/runtime"
	"github.com/go-chassis/go-chassis/core/config"
	"github.com/go-chassis/go-chassis/core/config/model"
	"github.com/go-chassis/go-chassis/core/lager"
	"github.com/go-chassis/go-chassis/core/server"
	"github.com/stretchr/testify/assert"
	"net"
	"net/http"
	"sync"
	"testing"
	"time"
)

func init() {
	lager.Init(&lager.Options{LoggerLevel: "INFO", RollingPolicy: "size"})
}

func externalIP() (net.IP, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return nil, err
		}
		for _, addr := range addrs {
			ip := getIpFromAddr(addr)
			if ip == nil {
				continue
			}
			return ip, nil
		}
	}
	return nil, errors.New("connected to the network?")
}

func getIpFromAddr(addr net.Addr) net.IP {
	var ip net.IP
	switch v := addr.(type) {
	case *net.IPNet:
		ip = v.IP
	case *net.IPAddr:
		ip = v.IP
	}
	if ip == nil || ip.IsLoopback() {
		return nil
	}
	ip = ip.To4()
	if ip == nil {
		return nil // not an ipv4 address
	}

	return ip
}

func TestHttpServer(t *testing.T) {
	config.Init()

	protoMap := make(map[string]model.Protocol)
	config.GlobalDefinition = &model.GlobalCfg{
		Cse: model.CseStruct{
			Protocols: protoMap,
		},
	}

	defaultChain := make(map[string]string)
	defaultChain["default"] = ""

	config.GlobalDefinition.Cse.Handler.Chain.Provider = defaultChain
	config.GlobalDefinition.Cse.Handler.Chain.Consumer = defaultChain

	f, err := server.GetServerFunc("http")
	assert.NoError(t, err)

	// case split port error
	s := f(server.Options{
		Address:   "0.0.0.130201",
		ChainName: "default",
	})
	err = s.Start()
	assert.Error(t, err)

	// case invalid host
	s = f(server.Options{
		Address:   "2.2.2.1990:30201",
		ChainName: "default",
	})
	err = s.Start()
	assert.Error(t, err)

	// case listening error
	s = f(server.Options{
		Address:   "99.0.0.1:30201",
		ChainName: "default",
	})
	err = s.Start()
	assert.Error(t, err)

	// case listen on 127.0.0.1
	s = f(server.Options{
		Address:   "127.0.0.1:30201",
		ChainName: "default",
	})
	err = s.Start()
	assert.NoError(t, err)

	s.Stop()
	time.Sleep(time.Second * 5)

	// case forbidden to listen on 0.0.0.0
	runtime.Role = common.RoleSidecar
	s = f(server.Options{
		Address:   "0.0.0.0:50201",
		ChainName: "default",
	})
	err = s.Start()
	assert.Error(t, err)

	// start sider
	eIP, err := externalIP()
	runtime.Role = common.RoleSidecar
	if err != nil {
		return
	}
	s = f(server.Options{
		Address:   eIP.String() + ":30303",
		ChainName: "default",
	})

	runtime.Role = common.RoleSidecar
	err = s.Start()
	assert.NoError(t, err)

	s.Stop()
	time.Sleep(time.Second * 5)
}

func TestHttpServer_Start(t *testing.T) {
	config.Init()

	protoMap := make(map[string]model.Protocol)
	config.GlobalDefinition = &model.GlobalCfg{
		Cse: model.CseStruct{
			Protocols: protoMap,
		},
	}

	defaultChain := make(map[string]string)
	defaultChain["default"] = ""

	config.GlobalDefinition.Cse.Handler.Chain.Provider = defaultChain
	config.GlobalDefinition.Cse.Handler.Chain.Consumer = defaultChain

	f, err := server.GetServerFunc("http")
	assert.NoError(t, err)
	s := f(server.Options{
		Address:   "127.0.0.1:40201",
		ChainName: "default",
	})

	s.Register(map[string]string{})

	err = s.Start()
	assert.NoError(t, err)

	name := s.String()
	assert.Equal(t, "http", name)

	var wg sync.WaitGroup
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		resp, err := http.Get("http://127.0.0.1:40201")
		assert.NoError(t, err)
		if err != nil {
			return
		}
		defer resp.Body.Close()

	}(&wg)

	wg.Wait()

	err = s.Stop()
	assert.NoError(t, err)
}

func TestGenTag(t *testing.T) {
	str := genTag("s1", "s2", "s3")
	assert.Equal(t, "s1.s2.s3", str)
}

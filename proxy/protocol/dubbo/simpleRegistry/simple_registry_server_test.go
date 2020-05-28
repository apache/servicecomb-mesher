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
	dubboclient "github.com/apache/servicecomb-mesher/proxy/protocol/dubbo/client"
	"github.com/apache/servicecomb-mesher/proxy/protocol/dubbo/dubbo"
	"github.com/apache/servicecomb-mesher/proxy/protocol/dubbo/proxy"
	"github.com/go-chassis/go-chassis/core/config"
	"github.com/go-chassis/go-chassis/core/config/model"
	"github.com/go-chassis/go-chassis/core/lager"
	"github.com/go-chassis/go-chassis/core/server"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

func init() {
	lager.Init(&lager.Options{LoggerLevel: "DEBUG"})
}

func TestSimpleDubboRegistryServer_Start(t *testing.T) {
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

	f, err := server.GetServerFunc("dubboSimpleRegistry")
	assert.NoError(t, err)

	// case split port error
	s := f(server.Options{
		Address:   "0.0.0.10201",
		ChainName: "default",
	})
	err = s.Start()
	assert.Error(t, err)
	// case invalid host
	s = f(server.Options{
		Address:   "2.2.2.1990:50201",
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

	// case ok
	s = f(server.Options{
		Address:   "127.0.0.1:50201",
		ChainName: "default",
	})
	err = s.Start()
	assert.NoError(t, err)

	s.Stop()
	time.Sleep(time.Second * 5)
}

func TestDubboServer(t *testing.T) {
	t.Log("Test dubboSimpleRegistry server function")

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

	f, err := server.GetServerFunc("dubboSimpleRegistry")
	assert.NoError(t, err)
	addr := "127.0.0.1:30401"
	s := f(server.Options{
		Address:   addr,
		ChainName: "default",
	})

	s.Register(map[string]string{})

	err = s.Start()
	assert.NoError(t, err)

	name := s.String()
	assert.Equal(t, "dubboSimpleRegistry", name)

	var wg sync.WaitGroup
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		clientMgr := dubboclient.NewClientMgr()
		var dubboClient *dubboclient.DubboClient
		dubboClient, err := clientMgr.GetClient(addr, time.Second*5)
		assert.NoError(t, err)

		req := new(dubbo.Request)
		req.SetMsgID(int64(11111111))
		req.SetVersion("1.0.0")
		req.SetEvent("ok")

		_, err = dubboClient.Send(req)

	}(&wg)

	wg.Wait()

	assert.True(t, dubboproxy.IsProvider)
	err = s.Stop()
	assert.NoError(t, err)
}

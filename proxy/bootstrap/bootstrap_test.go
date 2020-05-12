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

package bootstrap

import (
	"github.com/apache/servicecomb-mesher/proxy/cmd"
	"github.com/apache/servicecomb-mesher/proxy/common"
	_ "github.com/apache/servicecomb-mesher/proxy/pkg/egress/archaius"
	"github.com/go-chassis/go-archaius"
	"github.com/go-chassis/go-chassis/core/config"
	"github.com/go-chassis/go-chassis/core/config/model"
	"github.com/go-chassis/go-chassis/core/lager"
	"github.com/go-chassis/go-chassis/pkg/util/fileutil"
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"path/filepath"
	"sync"
	"testing"
	// rate limiter handler
	_ "github.com/go-chassis/go-chassis/middleware/ratelimiter"
)

var o sync.Once = sync.Once{}

var yamlContent = `---
egress:
  infra: cse # pilot or cse
  address: http://istio-pilot.istio-system:15010
egressRule:
  google-ext:
    - hosts:
        - "www.google.com"
        - "*.yahoo.com"
      ports:
        - port: 80
          protocol: HTTP
  facebook-ext:
    - hosts:
        - "www.facebook.com"
      ports:
        - port: 80
          protocol: HTTP`

func TestBootstrap(t *testing.T) {
	lager.Init(&lager.Options{LoggerLevel: "DEBUG"})

	// init work dir
	os.Setenv(fileutil.ChassisHome, filepath.Join("...", "..."))
	os.Setenv(fileutil.ChassisConfDir, filepath.Join("...", "...", "conf"))
	t.Log(os.Getenv("CHASSIS_HOME"))

	// init archaius
	archaius.Init(archaius.WithENVSource())

	//ini config
	config.Init()

	protoMap := make(map[string]model.Protocol)
	protoMap["http"] = model.Protocol{
		Listen: "127.0.0.1:90909",
	}
	config.GlobalDefinition = &model.GlobalCfg{
		Cse: model.CseStruct{
			Protocols: protoMap,
		},
	}

	configMesher := "../../conf/mesher.yaml"
	os.Args = []string{"test", "--config", configMesher}
	if err := cmd.Init(); err != nil {
		panic(err)
	}
	if err := cmd.Configs.GeneratePortsMap(); err != nil {
		panic(err)
	}

	// init egress.yaml file
	d, _ := os.Getwd()
	os.Mkdir(filepath.Join(d, "conf"), os.ModePerm)
	filename := filepath.Join(d, "conf", "egress.yaml")
	os.Remove(filename)
	f1, err := os.Create(filename)
	assert.NoError(t, err)
	defer f1.Close()
	_, err = io.WriteString(f1, yamlContent)
	assert.NoError(t, err)

	t.Run("Test RegisterFramework", func(t *testing.T) {
		// case cmd.Configs.Role is empty
		cmd.Configs.Role = ""
		RegisterFramework()
		// case cmd.Configs.Role == common.RoleSidecar
		cmd.Configs.Role = common.RoleSidecar
		RegisterFramework()
	})

	t.Run("Test Start", func(t *testing.T) {
		// case Protocols is empty
		config.GlobalDefinition.Cse.Protocols = map[string]model.Protocol{}
		err := Start()
		assert.Error(t, err)

		// cmd.Configs.LocalServicePorts = "http:9090"
		cmd.Configs.LocalServicePorts = "http:9090"
		err = Start()

		cmd.Configs.LocalServicePorts = ""
		RegisterFramework()
		SetHandlers()
		err = InitEgressChain()
		assert.NoError(t, err)

		err = Start()
		assert.NoError(t, err)

	})
}

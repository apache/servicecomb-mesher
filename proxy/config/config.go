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

package config

import (
	"github.com/apache/servicecomb-mesher/proxy/cmd"
	"github.com/apache/servicecomb-mesher/proxy/common"
	"github.com/go-chassis/go-archaius"
	"github.com/go-chassis/go-archaius/core"
	"github.com/go-chassis/go-chassis/core/config"
	"github.com/go-chassis/go-chassis/core/config/model"
	"github.com/go-chassis/go-chassis/core/lager"
	"github.com/go-chassis/go-chassis/core/server"
	"github.com/go-chassis/go-chassis/pkg/util/fileutil"
	"github.com/go-mesh/openlogging"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
)

//Constant for mesher conf file
const (
	ConfFile        = "mesher.yaml"
	EgressConfFile  = "egress.yaml"
	ConfigCenterKie = "servicecomb-kie"
)

var mesherConfig *MesherConfig
var egressConfig *EgressConfig

//ConfigCenterEvent event implemented.
type ConfigCenterEvent struct {
}

//Event is used for Configcenter changes.
func (configEvent *ConfigCenterEvent) Event(event *core.Event) {
	openlogging.GetLogger().Debugf("Configcenter Event. event:%s %s %s", event.EventType, event.Key, event.EventSource)
	if event.EventType == "UPDATE" {
		globalDef := &model.GlobalCfg{}
		err := yaml.Unmarshal([]byte(event.Value.(string)), globalDef)
		if err != nil {
			openlogging.GetLogger().Warn("Config chassis.yaml change unmarshal err:" + err.Error())
			return
		}
		openlogging.GetLogger().Debugf("Configcener Event succ. event:%s %s %s", event.EventType, event.Key, event.EventSource)
	}
}

//GetConfig returns mesher config
func GetConfig() *MesherConfig {
	return mesherConfig
}

//SetConfig sets new mesher config from input config
func SetConfig(nc *MesherConfig) {
	if mesherConfig == nil {
		mesherConfig = &MesherConfig{}
	}
	*mesherConfig = *nc
}

//GetEgressConfig returns Egress config
func GetEgressConfig() *EgressConfig {
	return egressConfig
}

//GetConfigFilePath returns config file path
func GetConfigFilePath(key string) (string, error) {
	if cmd.Configs.ConfigFile == "" || key == EgressConfFile {
		wd, err := fileutil.GetWorkDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(wd, "conf", key), nil
	}
	return cmd.Configs.ConfigFile, nil
}

//InitProtocols initiates protocols
func InitProtocols() error {
	// todo if sdk init failed, do not call the data
	if len(config.GlobalDefinition.Cse.Protocols) == 0 {
		config.GlobalDefinition.Cse.Protocols = map[string]model.Protocol{
			common.HTTPProtocol: {Listen: "127.0.0.1:30101"},
		}

		return server.Init()
	}
	return nil
}

//Init reads config and initiates
func Init() error {
	mesherConfig = &MesherConfig{}
	contents, err := GetConfigContents(ConfFile)
	if err != nil {
		return err
	}
	if err := yaml.Unmarshal([]byte(contents), mesherConfig); err != nil {
		return err
	}

	egressConfig = &EgressConfig{}
	egressContents, err := GetConfigContents(EgressConfFile)
	if err != nil {
		return err
	}
	if err := yaml.Unmarshal([]byte(egressContents), egressConfig); err != nil {
		return err
	}
	//Register event in archaius for kie.
	if config.GetConfigCenterConf().Type == ConfigCenterKie {
		var configEvent = &ConfigCenterEvent{}
		err := archaius.RegisterListener(configEvent, fileutil.Global)
		if err != nil {
			openlogging.GetLogger().Errorf("Archaius RegisterListener failed. key:%s, err:%s", fileutil.Global, err.Error())
		}
		openlogging.GetLogger().Debugf("Archaius RegisterListener succ. key:%s", fileutil.Global)
	}

	return nil
}

//GetConfigContents returns config contents
func GetConfigContents(key string) (string, error) {
	f, err := GetConfigFilePath(key)
	if err != nil {
		return "", err
	}
	var contents string
	//route rule yaml file's content is value of a key
	//So read from config center first,if it is empty, Try to set file content into memory key value
	contents = archaius.GetString(key, "")
	if contents == "" {
		contents = SetKeyValueByFile(key, f)
	}
	return contents, nil
}

//SetKeyValueByFile reads mesher.yaml and gets key and value
func SetKeyValueByFile(key, f string) string {
	var contents string
	if _, err := os.Stat(f); err != nil {
		openlogging.GetLogger().Warn(err.Error())
		return ""
	}
	b, err := ioutil.ReadFile(f)
	if err != nil {
		lager.Logger.Error("Can not read yaml file" + err.Error())
		return ""
	}
	contents = string(b)
	archaius.AddKeyValue(key, contents)
	return contents
}

// GetEgressEndpoints returns the pilot address
func GetEgressEndpoints() string {
	return egressConfig.Egress.Address
}

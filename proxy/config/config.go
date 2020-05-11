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
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/apache/servicecomb-mesher/proxy/cmd"
	"github.com/apache/servicecomb-mesher/proxy/common"
	"github.com/go-chassis/go-archaius"
	"github.com/go-chassis/go-chassis/core/config"
	"github.com/go-chassis/go-chassis/core/config/model"
	"github.com/go-chassis/go-chassis/core/lager"
	"github.com/go-chassis/go-chassis/core/server"
	"github.com/go-chassis/go-chassis/pkg/util/fileutil"
	"github.com/go-mesh/openlogging"
	"gopkg.in/yaml.v2"
)

//Constant for mesher conf file
const (
	ConfFile       = "mesher.yaml"
	EgressConfFile = "egress.yaml"
)

var mesherConfig *MesherConfig
var egressConfig *EgressConfig

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
	egressConfig = &EgressConfig{}
	p, err := GetConfigFilePath(ConfFile)
	if err != nil {
		return err
	}
	err = archaius.AddFile(p)
	if err != nil {
		return err
	}

	contents, err := GetConfigContents(ConfFile)
	if err != nil {
		return err
	}
	if err := yaml.Unmarshal([]byte(contents), mesherConfig); err != nil {
		return err
	}

	egressContents, err := GetConfigContents(EgressConfFile)
	if err != nil {
		return err
	}
	if err := yaml.Unmarshal([]byte(egressContents), egressConfig); err != nil {
		return err
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
	err = archaius.Set(key, contents)
	if err != nil {
		lager.Logger.Error("Archaius set error: " + err.Error())
	}
	return contents
}

// GetEgressEndpoints returns the pilot address
func GetEgressEndpoints() string {
	return egressConfig.Egress.Address
}

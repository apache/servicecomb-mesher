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

package chassis

import (
	"errors"
	"github.com/go-chassis/go-archaius/core"
	"github.com/go-chassis/go-archaius/core/config-manager"
	"github.com/go-chassis/go-archaius/core/event-system"
	"github.com/go-chassis/go-archaius/sources/file-source"
	"github.com/go-chassis/go-chassis/core/lager"
	"github.com/go-chassis/go-chassis/pkg/util/fileutil"
	"github.com/go-mesh/mesher/config/model"
	controlarchaius "github.com/go-mesh/mesher/control/archiaus"
	"github.com/go-mesh/mesher/pkg/egress"
	"gopkg.in/yaml.v2"
	"path/filepath"
)

//EgressYaml egress yaml file name
const EgressYaml = "egress.yaml"

var egressRuleMgr core.ConfigMgr

type egressRuleEventListener struct{}

// update egress rule of a service
func (r *egressRuleEventListener) Event(e *core.Event) {
	if e == nil {
		lager.Logger.Warn("Event pointer is nil", nil)
		return
	}

	v := egressRuleMgr.GetConfigurationsByKey(e.Key)
	if v == nil {
		controlarchaius.SaveToEgressCache(nil)
		lager.Logger.Infof("[%s] Egress rule is removed", e.Key)
		return
	}
	var egressconfig model.EgressConfig

	if err := yaml.Unmarshal([]byte(v.([]byte)), &egressconfig); err != nil {
		lager.Logger.Error("yaml unmarshal failed", nil)
		return
	}
	var egressRules []*model.EgressRule

	for key, value := range egressconfig.Destinations {
		ok, _ := egress.ValidateEgressRule(map[string][]*model.EgressRule{key: value})
		if !ok {
			lager.Logger.Warn("Validating Egress Rule Failed")
			return

		}
		egressRules = append(egressRules, value...)
	}

	controlarchaius.SaveToEgressCache(&egressconfig)
	lager.Logger.Infof("Update [%s] egress rule SUCCESS", e.Key)
}

// initialize the config mgr and add several sources
func initEgressManager() error {
	d := eventsystem.NewDispatcher()
	l := &egressRuleEventListener{}
	d.RegisterListener(l, ".*")
	egressRuleMgr = configmanager.NewConfigurationManager(d)

	if err := AddEgressRuleSource(); err != nil {
		return err
	}
	return nil
}

// AddEgressRuleSource adds a config source to egress rule manager
// Do not call this method until egress init success
func AddEgressRuleSource() error {
	if egressRuleMgr == nil {
		return errors.New("egressRuleMgr is nil, please init it firstly")
	}
	fsource := filesource.NewFileSource()
	fsource.AddFile(filepath.Join(fileutil.GetConfDir(), EgressYaml), filesource.DefaultFilePriority, filesource.Convert2configMap)
	if err := egressRuleMgr.AddSource(fsource, filesource.DefaultFilePriority); err != nil {
		return err
	}
	lager.Logger.Infof("Add [%s] source success", fsource.GetSourceName())
	return nil
}

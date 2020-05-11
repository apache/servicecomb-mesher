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

package archaius

import (
	"github.com/go-chassis/go-archaius/event"
	"github.com/go-chassis/go-archaius/source/util"
	"gopkg.in/yaml.v2"
	"path/filepath"
	"strings"

	"github.com/apache/servicecomb-mesher/proxy/config"
	"github.com/apache/servicecomb-mesher/proxy/pkg/egress"
	"github.com/go-chassis/go-archaius"
	"github.com/go-chassis/go-chassis/core/lager"
	"github.com/go-chassis/go-chassis/pkg/util/fileutil"
)

//EgressYaml egress yaml file name
const EgressYaml = "egress.yaml"

type egressRuleEventListener struct{}

// update egress rule of a service
func (r *egressRuleEventListener) Event(e *event.Event) {
	if e == nil {
		lager.Logger.Warn("Event pointer is nil", nil)
		return
	}
	if !strings.Contains(e.Key, EgressYaml) {
		return
	}
	v := archaius.Get(e.Key)
	if v == nil {
		lager.Logger.Infof("[%s] Error getting egress key", e.Key)
		return
	}

	var egressconfig config.EgressConfig

	if err := yaml.Unmarshal([]byte(v.([]byte)), &egressconfig); err != nil {
		lager.Logger.Error("yaml unmarshal failed", nil)
		return
	}
	var egressRules []*config.EgressRule

	for key, value := range egressconfig.Destinations {
		ok, _ := egress.ValidateEgressRule(map[string][]*config.EgressRule{key: value})
		if !ok {
			lager.Logger.Warn("Validating Egress Rule Failed")
			return

		}
		egressRules = append(egressRules, value...)
	}

	SetEgressRule(map[string][]*config.EgressRule{e.Key: egressRules})
	lager.Logger.Infof("Update [%s] egress rule SUCCESS", e.Key)
}

// initialize the config mgr and add several sources
func initEgressManager() error {
	egressListener := &egressRuleEventListener{}
	err := archaius.AddFile(filepath.Join(fileutil.GetConfDir(), EgressYaml), archaius.WithFileHandler(util.UseFileNameAsKeyContentAsValue))
	if err != nil {
		lager.Logger.Infof("Archaius add file failed: ", err)
	}
	err = archaius.RegisterListener(egressListener, ".*")
	if err != nil {
		lager.Logger.Infof("Archaius add file failed: ", err)
	}
	return nil
}

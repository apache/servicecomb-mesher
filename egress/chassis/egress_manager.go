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
	"sync"

	"github.com/go-chassis/go-archaius/core"
	"github.com/go-chassis/go-archaius/core/config-manager"
	"github.com/go-chassis/go-archaius/core/event-system"
	"github.com/go-chassis/go-chassis/core/lager"
	"github.com/go-mesh/mesher/config/model"
	"github.com/go-mesh/mesher/egress"
)

const egressFileSourceName = "EgressFileSource"
const egressFileSourcePriority = 16

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
		DeleteEgressRuleByKey(e.Key)
		lager.Logger.Infof("[%s] Egress rule is removed", e.Key)
		return
	}
	egressRules, ok := v.([]*model.EgressRule)
	if !ok {
		lager.Logger.Error("value is not type []*RouteRule", nil)
		return
	}

	ok, _ = egress.ValidateEgressRule(map[string][]*model.EgressRule{e.Key: egressRules})
	if ok {
		SetEgressRuleByKey(e.Key, egressRules)
		lager.Logger.Infof("Update [%s] egress rule success", e.Key)
	}
}

// egressFileSource keeps the egress rule in egress file,
// after init, it's data does not change
type egressFileSource struct {
	once sync.Once
	d    map[string]interface{}
}

func newEgressFileSource() *egressFileSource {
	r := &egressFileSource{}
	r.once.Do(func() {
		egressRules := GetEgressRule()

		d := make(map[string]interface{}, 0)
		if egressRules == nil {
			r.d = d
			lager.Logger.Error("Can not get any egress config", nil)
			return
		}
		for k, v := range egressRules {
			d[k] = v
		}
		r.d = d
	})
	return r
}

func (r *egressFileSource) GetSourceName() string {
	return egressFileSourceName
}
func (r *egressFileSource) GetConfigurations() (map[string]interface{}, error) {
	configMap := make(map[string]interface{})
	for k, v := range r.d {
		configMap[k] = v
	}
	return configMap, nil
}
func (r *egressFileSource) GetConfigurationsByDI(dimensionInfo string) (map[string]interface{}, error) {
	return nil, nil
}
func (r *egressFileSource) GetConfigurationByKey(k string) (interface{}, error) {
	v, ok := r.d[k]
	if !ok {
		return nil, errors.New("key " + k + " not exist")
	}
	return v, nil
}
func (r *egressFileSource) GetConfigurationByKeyAndDimensionInfo(key, dimensionInfo string) (interface{}, error) {
	return nil, nil
}
func (r *egressFileSource) AddDimensionInfo(dimensionInfo string) (map[string]string, error) {
	return nil, nil
}
func (r *egressFileSource) DynamicConfigHandler(core.DynamicConfigCallback) error {
	return nil
}
func (r *egressFileSource) GetPriority() int {
	return egressFileSourcePriority
}
func (r *egressFileSource) Cleanup() error { return nil }

// initialize the config mgr and add several sources
func initEgressManager() error {
	d := eventsystem.NewDispatcher()
	l := &egressRuleEventListener{}
	d.RegisterListener(l, ".*")
	egressRuleMgr = configmanager.NewConfigurationManager(d)
	if err := AddEgressRuleSource(newEgressFileSource()); err != nil {
		return err
	}
	return nil
}

// AddEgressRuleSource adds a config source to egress rule manager
// Do not call this method until egress init success
func AddEgressRuleSource(s core.ConfigSource) error {
	if s == nil {
		return errors.New("source nil")
	}
	if egressRuleMgr == nil {
		return errors.New("egressRuleMgr is nil, please init it firstly")
	}
	if err := egressRuleMgr.AddSource(s, s.GetPriority()); err != nil {
		return err
	}
	lager.Logger.Infof("Add [%s] source success", s.GetSourceName())
	return nil
}

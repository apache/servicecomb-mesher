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
	"fmt"
	"github.com/go-mesh/mesher/config/model"
	"github.com/go-mesh/mesher/egress"
	"sync"
)

//Egress is cse Egress service
type Egress struct {
}

//FetchEgressRule return all rules
func (r *Egress) FetchEgressRule() map[string][]*model.EgressRule {
	return GetEgressRule()
}

//SetEgressRule set rules
func (r *Egress) SetEgressRule(rr map[string][]*model.EgressRule) {
	SetEgressRule(rr)
}

//FetchEgressRuleByName get rules by name
func (r *Egress) FetchEgressRuleByName(name string) []*model.EgressRule {
	return GetEgressRuleByKey(name)
}

//Init init router config
func (r *Egress) Init(op egress.Options) error {
	// the manager use dests to init, so must init after dests
	if err := initEgressManager(); err != nil {
		return err
	}
	return refresh()
}

func newEgress() (egress.Egress, error) {
	return &Egress{}, nil
}

// refresh all the egress config
func refresh() error {
	configs := egressRuleMgr.GetConfigurations()
	d := make(map[string][]*model.EgressRule)
	for k, v := range configs {
		rules, ok := v.([]*model.EgressRule)
		if !ok {
			err := fmt.Errorf("Egress rule type assertion fail, key: %s", k)
			return err
		}
		d[k] = rules
	}

	ok, _ := egress.ValidateEgressRule(d)
	if ok {
		dests = d
	}
	return nil
}

var dests = make(map[string][]*model.EgressRule)
var lock sync.RWMutex

// SetEgressRuleByKey set egress rule by key
func SetEgressRuleByKey(k string, r []*model.EgressRule) {
	lock.Lock()
	dests[k] = r
	lock.Unlock()
}

// DeleteEgressRuleByKey set egress rule by key
func DeleteEgressRuleByKey(k string) {
	lock.Lock()
	delete(dests, k)
	lock.Unlock()
}

// GetEgressRuleByKey get egress rule by key
func GetEgressRuleByKey(k string) []*model.EgressRule {
	lock.RLock()
	defer lock.RUnlock()
	return dests[k]
}

// GetEgressRule get egress rule
func GetEgressRule() map[string][]*model.EgressRule {
	lock.RLock()
	defer lock.RUnlock()
	return dests
}

// SetEgressRule set egress rule
func SetEgressRule(rule map[string][]*model.EgressRule) {
	lock.RLock()
	defer lock.RUnlock()
	dests = rule
}
func init() {
	egress.InstallEgressService("cse", newEgress)
}

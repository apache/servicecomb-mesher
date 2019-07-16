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
	"fmt"
	"github.com/apache/servicecomb-mesher/proxy/config"
	"github.com/apache/servicecomb-mesher/proxy/pkg/egress"
	"sync"
)

//Egress is cse Egress service
type Egress struct {
}

//FetchEgressRule return all rules
func (r *Egress) FetchEgressRule() map[string][]*config.EgressRule {
	return GetEgressRule()
}

//SetEgressRule set rules
func (r *Egress) SetEgressRule(rr map[string][]*config.EgressRule) {
	SetEgressRule(rr)
}

//Init init egress config
func (r *Egress) Init(op egress.Options) error {
	// the manager use dests to init, so must init after dests
	if err := initEgressManager(); err != nil {
		return err
	}
	return refresh()
}

// refresh all the egress config
func refresh() error {
	configs := config.GetEgressConfig()
	ok, _ := egress.ValidateEgressRule(configs.Destinations)
	if !ok {
		err := fmt.Errorf("Egress rule type assertion fail, key: %s", configs.Destinations)
		return err
	}

	dests = configs.Destinations
	return nil
}

func newEgress() (egress.Egress, error) {
	return &Egress{}, nil
}

var dests = make(map[string][]*config.EgressRule)
var lock sync.RWMutex

// GetEgressRule get egress rule
func GetEgressRule() map[string][]*config.EgressRule {
	lock.RLock()
	defer lock.RUnlock()
	return dests
}

// SetEgressRule set egress rule
func SetEgressRule(rule map[string][]*config.EgressRule) {
	lock.RLock()
	defer lock.RUnlock()
	dests = rule
}

func init() {
	egress.InstallEgressService("cse", newEgress)
}

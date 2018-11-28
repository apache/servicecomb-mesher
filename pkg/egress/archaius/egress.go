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
	"github.com/go-mesh/mesher/config/model"
	"github.com/go-mesh/mesher/pkg/egress"
	"sync"
)

//Egress is cse Egress service
type Egress struct {
}

//FetchEgressRule return all rules
func (r *Egress) FetchEgressRule() map[string][]*model.EgressRule {
	return GetEgressRule()
}

//Init init egress config
func (r *Egress) Init(op egress.Options) error {
	// the manager use dests to init, so must init after dests
	if err := initEgressManager(); err != nil {
		return err
	}
	return nil
}

func newEgress() (egress.Egress, error) {
	return &Egress{}, nil
}

var dests = make(map[string][]*model.EgressRule)
var lock sync.RWMutex

// GetEgressRule get egress rule
func GetEgressRule() map[string][]*model.EgressRule {
	lock.RLock()
	defer lock.RUnlock()
	return dests
}

func init() {
	egress.InstallEgressService("cse", newEgress)
}

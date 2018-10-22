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

package egress

import (
	"errors"
	"github.com/go-chassis/go-chassis/control"
	"github.com/go-mesh/mesher/config/model"
	meshercontrol "github.com/go-mesh/mesher/control"
	"regexp"
	"sync"
)

var lock sync.RWMutex

var plainHosts = make(map[string]*model.EgressRule)
var regexHosts = make(map[string]*model.EgressRule)

//Egress return egress rule, you can also set custom egress rule
type Egress interface {
	Init(Options) error
	SetEgressRule(map[string][]*model.EgressRule)
	FetchEgressRule() map[string][]*model.EgressRule
	FetchEgressRuleByName(string) []*model.EgressRule
}

// ErrNoExist means if there is no egress implementation
var ErrNoExist = errors.New("Egress not exists")
var egressServices = make(map[string]func() (Egress, error))

// DefaultEgress is current egress implementation
var DefaultEgress Egress

// InstallEgressService install router service for developer
func InstallEgressService(name string, f func() (Egress, error)) {
	egressServices[name] = f
}

//BuildEgress create a Egress
func BuildEgress(name string) error {
	f, ok := egressServices[name]
	if !ok {
		return ErrNoExist
	}
	r, err := f()
	if err != nil {
		return err
	}
	DefaultEgress = r
	return nil
}

//Match Check Egress rule matches
func Match(hostname string) (bool, *control.EgressConfig) {
	EgressRules := meshercontrol.DefaultPanelEgress.GetEgressRule()
	for _, egress := range EgressRules {

		for _, host := range egress.Hosts {
			// Check host length greater than 0 and does not
			// start with *
			if len(host) > 0 && string(host[0]) != "*" {
				if host == hostname {
					return true, &egress
				}
			} else if string(host[0]) == "*" {
				substring := host[1:]
				match, _ := regexp.MatchString(substring+"$", hostname)
				if match == true {
					return true, &egress
				}
			}
		}
	}
	return false, nil
}

//SplitEgressRules Check Egress rule matches
func SplitEgressRules() (map[string]*model.EgressRule, map[string]*model.EgressRule) {
	EgressRules := DefaultEgress.FetchEgressRule()
	for _, egressRules := range EgressRules {
		for _, egress := range egressRules {

			for _, host := range egress.Hosts {
				if len(host) > 1 && string(host[0]) != "*" {
					plainHosts[host] = egress
				} else if string(host[0]) == "*" {
					substring := host[1:]
					regexHosts[substring] = egress
				}
			}
		}
	}

	return plainHosts, regexHosts
}

//MatchHost Check Egress rule matches
func MatchHost(hostname string) (bool, *model.EgressRule) {
	if val, ok := plainHosts[hostname]; ok {
		return true, val
	}

	for key, value := range regexHosts {
		match, _ := regexp.MatchString(key+"$", hostname)
		if match == true {
			return true, value

		}
	}
	return false, nil
}

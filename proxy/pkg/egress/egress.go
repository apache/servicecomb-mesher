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
	"github.com/apache/servicecomb-mesher/proxy/config"
	meshercontrol "github.com/apache/servicecomb-mesher/proxy/control"
	"github.com/go-chassis/go-chassis/control"
	"regexp"
	"sync"
)

var lock sync.RWMutex

var plainHosts = make(map[string]*config.EgressRule)
var regexHosts = make(map[string]*config.EgressRule)

//Egress return egress rule, you can also set custom egress rule
type Egress interface {
	Init(Options) error
	SetEgressRule(map[string][]*config.EgressRule)
	FetchEgressRule() map[string][]*config.EgressRule
}

// ErrNoExist means if there is no egress implementation
var ErrNoExist = errors.New("Egress not exists")
var egressServices = make(map[string]func() (Egress, error))

// DefaultEgress is current egress implementation
var DefaultEgress Egress

// InstallEgressService install egress service for developer
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
	var EgressRules []control.EgressConfig
	if meshercontrol.DefaultPanelEgress != nil {
		EgressRules = meshercontrol.DefaultPanelEgress.GetEgressRule()
	} else {
		mapEgressRules := DefaultEgress.FetchEgressRule()
		for _, value := range mapEgressRules {
			for _, rule := range value {
				var Ports []*control.EgressPort
				for _, port := range rule.Ports {
					p := control.EgressPort{
						Port:     (*port).Port,
						Protocol: (*port).Protocol,
					}
					Ports = append(Ports, &p)
				}
				c := control.EgressConfig{
					Hosts: rule.Hosts,
					Ports: Ports,
				}
				EgressRules = append(EgressRules, c)
			}
		}
	}

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

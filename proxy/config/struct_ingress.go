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
	"github.com/go-chassis/foundation/string"
	"gopkg.in/yaml.v2"
)

//Mesher is prefix
type Mesher struct {
	Ingress Ingress `yaml:"ingress"`
}

//Ingress hold rules and other settings
type Ingress struct {
	Rule map[string]string `yaml:"rule"`
	Type string            `yaml:"type"`
}

//IngressRules is ingress rules slice
type IngressRules []*IngressRule

//Len return the length of rule
func (r IngressRules) Len() int {
	return len(r)
}

//Value return the rule
func (r IngressRules) Value() []*IngressRule {
	return r
}

//NewRules create a rule by raw data
func NewRules(raw string) (*IngressRules, error) {
	b := stringutil.Str2bytes(raw)
	r := &IngressRules{}
	err := yaml.Unmarshal(b, r)
	return r, err
}

//IngressRule is a ingress rule
type IngressRule struct {
	Host    string  `yaml:"host"`
	Limit   int     `yaml:"limit"`
	APIPath string  `yaml:"apiPath"`
	Service Service `yaml:"service"`
}

//Service is upstream info
type Service struct {
	Name         string            `yaml:"name"`
	Tags         map[string]string `yaml:"tags"`
	RedirectPath string            `yaml:"redirectPath"`
	Port         Port              `yaml:"port"`
}

//Port is service port information
type Port struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

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

//EgressConfig is the struct having info about egress rule destinations
type EgressConfig struct {
	Egress       Egress                   `yaml:"egress"`
	Destinations map[string][]*EgressRule `yaml:"egressRule"`
}

// Egress define where rule comes from
type Egress struct {
	Infra   string `yaml:"infra"`
	Address string `yaml:"address"`
}

//EgressRule has hosts and ports information
type EgressRule struct {
	Hosts []string      `yaml:"hosts"`
	Ports []*EgressPort `yaml:"ports"`
}

//EgressPort protocol and the corresponding port
type EgressPort struct {
	Port     int32  `yaml:"port"`
	Protocol string `yaml:"protocol"`
}

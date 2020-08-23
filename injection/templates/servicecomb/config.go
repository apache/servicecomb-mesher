/*
 *  Licensed to the Apache Software Foundation (ASF) under one or more
 *  contributor license agreements.  See the NOTICE file distributed with
 *  this work for additional information regarding copyright ownership.
 *  The ASF licenses this file to You under the Apache License, Version 2.0
 *  (the "License"); you may not use this file except in compliance with
 *  the License.  You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 */

package servicecomb

// Config of ServiceComb sidecar template
type Config struct {
	App           App           `yaml:"-"`
	Mesher        Mesher        `yaml:"mesher"`
	ServiceCenter ServiceCenter `yaml:"serviceCenter"`
}

// App metadata, need to read from the k8s pod
type App struct {
	Name    string
	Host    string
	Port    int32
	Version string
}

// Mesher startup configuration
type Mesher struct {
	Name      string `yaml:"name"`
	Image     string `yaml:"image"`
	Tag       string `yaml:"tag"`
	GRPCPort  int    `yaml:"grpcPort"`
	HTTPPort  int    `yaml:"httpPort"`
	AdminPort int    `yaml:"adminPort"`
}

// ServiceCenter registry configuration
type ServiceCenter struct {
	Name      string     `yaml:"name"`
	Namespace string     `yaml:"namespace"`
	Address   string     `yaml:"address"`
	TlsConfig *TLSConfig `yaml:"tlsConfig"`
}

// TLSConfig tls configuration
type TLSConfig struct {
	CaFile   string `yaml:"caFile"`
	CertFile string `yaml:"certFile"`
	KeyFile  string `yaml:"keyFile"`
}

// DefaultConfig return default configuration
func DefaultConfig() *Config {
	return &Config{
		App: App{},
		Mesher: Mesher{
			Name:      "svccomb-mesher",
			Image:     "servicecomb/mesher-sidecar",
			Tag:       "latest",
			GRPCPort:  40101,
			HTTPPort:  30101,
			AdminPort: 30102,
		},
		ServiceCenter: ServiceCenter{
			Name:      "servicecenter",
			Namespace: "svccomb-system",
		},
	}
}

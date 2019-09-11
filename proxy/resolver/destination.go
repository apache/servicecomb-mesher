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

package resolver

import (
	"errors"
	"log"
	"net/url"

	"fmt"
	"github.com/apache/servicecomb-mesher/proxy/config"
	"github.com/go-mesh/openlogging"
)

var drMap = make(map[string]DestinationResolver)

//DestinationResolverPlugins is a map
var DestinationResolverPlugins map[string]func() DestinationResolver

//SelfEndpoint is a string
var SelfEndpoint = "#To be init#"

//DefaultPlugin is a constant which stores default plugin name
const DefaultPlugin = "host"

//ErrUnknownResolver is of type error
var ErrUnknownResolver = errors.New("unknown Destination Resolver")

//DestinationResolver is a interface with Resolve method
type DestinationResolver interface {
	Resolve(remoteIP, host, rawURI string, header map[string]string) (string, string, error)
}

//DefaultDestinationResolver is a struct
//mesher as sidecar must use DefaultDestinationResolver
type DefaultDestinationResolver struct {
}

//Resolve resolves service's endpoint
//service may have multiple port for same protocol
func (dr *DefaultDestinationResolver) Resolve(remoteIP, host, rawURI string, header map[string]string) (string, string, error) {
	u, err := url.Parse(rawURI)
	if err != nil {
		openlogging.Error("Can not parse url: " + err.Error())
		return "", "", err
	}

	if u.Host == "" {
		return "", "", errors.New(`Invalid uri, please check:
1, For provider, mesher listens on external ip
2, Set http_proxy as mesher address, before sending request`)
	}

	if u.Host == SelfEndpoint {
		return "", "", errors.New(`uri format must be: http://serviceName/api`)
	}

	return u.Hostname(), u.Port(), nil
}

//New function returns new DefaultDestinationResolver struct object
func New() DestinationResolver {
	return &DefaultDestinationResolver{}
}

//GetDestinationResolver returns destinationResolver pointer
func GetDestinationResolver(name string) DestinationResolver {
	return drMap[name]
}

//InstallDestinationResolverPlugin function installs new plugin
func InstallDestinationResolverPlugin(name string, newFunc func() DestinationResolver) {
	DestinationResolverPlugins[name] = newFunc
	log.Printf("Installed DestinationResolver Plugin, name=%s", name)
}

//SetDefaultDestinationResolver set the a default implementation for a protocol, so that you don't need to set config file
func SetDefaultDestinationResolver(name string, dr DestinationResolver) {
	drMap[name] = dr
	log.Printf("Installed default DestinationResolver for [%s]", name)
}
func init() {
	DestinationResolverPlugins = make(map[string]func() DestinationResolver)
	InstallDestinationResolverPlugin(DefaultPlugin, New)
	SetDefaultDestinationResolver("http", &DefaultDestinationResolver{})
}

//Init function reads config and initiates it
func Init() error {
	if config.GetConfig().Plugin != nil {
		for name, v := range config.GetConfig().Plugin.DestinationResolver {
			if v == "" {
				v = DefaultPlugin
			}
			f, ok := DestinationResolverPlugins[v]
			if !ok {
				return fmt.Errorf("unknown destination resolver [%s]", v)
			}
			drMap[name] = f()
		}
	}

	return nil
}

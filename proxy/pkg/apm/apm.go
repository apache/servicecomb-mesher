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

package apm

import (
	"github.com/apache/servicecomb-mesher/proxy/config"
	"github.com/go-chassis/go-chassis/core/invocation"
	"github.com/go-mesh/openlogging"
	"strconv"
)

type ApmClient interface {
	CreateSpans(i *invocation.Invocation) ([]interface{}, error)
	EndSpans(spans []interface{}, status int) error
	CreateEntrySpan(i *invocation.Invocation) (interface{}, error)
	CreateExitSpan(i *invocation.Invocation) (interface{}, error)
	EndSpan(sp interface{}, statusCode int) error
}

type ApmManager struct {
	apmClientPlugins map[string]func(Options) (ApmClient, error)
	apmClients       map[string]ApmClient
}

var apmClientPlugins = make(map[string]func(Options) (ApmClient, error))
var apmClients = make(map[string]ApmClient)

//InstallClientPlugins register apmclient create func
func InstallClientPlugins(name string, f func(Options) (ApmClient, error)) {
	apmClientPlugins[name] = f
	openlogging.Info("Install apm client: " + name)
}

//CreateSpans use invocation to make spans for apm
func CreateSpans(i *invocation.Invocation) ([]interface{}, error) {
	openlogging.Info("CreateSpans")
	if client, ok := apmClients[config.GetConfig().APM.ApmName]; ok {
		openlogging.Info("client.CreateSpans")
		return client.CreateSpans(i)
	}
	var spans []interface{}
	return spans, nil
}

//EndSpans use invocation to make spans of apm end
func EndSpans(spans []interface{}, status int) error {
	openlogging.Info("EndSpans" + strconv.Itoa(status))
	if client, ok := apmClients[config.GetConfig().APM.ApmName]; ok {
		return client.EndSpans(spans, status)
	}
	return nil
}

//Init apm client
func Init() {
	openlogging.Info("Apm Init " + config.GetConfig().APM.ApmName + " " + config.GetConfig().APM.ServerURI + " " + strconv.FormatBool(config.GetConfig().APM.Enable))
	if config.GetConfig().APM.Enable == true {
		f, ok := apmClientPlugins[config.GetConfig().APM.ApmName]
		if ok {
			client, err := f(Options{Name: config.GetConfig().APM.ApmName, ServerUri: config.GetConfig().APM.ServerURI})
			if err == nil {
				apmClients[config.GetConfig().APM.ApmName] = client
			} else {
				openlogging.Error("apmClients init failed. " + err.Error())
			}
		}
	}
}

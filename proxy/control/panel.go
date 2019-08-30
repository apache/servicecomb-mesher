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

package control

import (
	"fmt"
	"github.com/go-chassis/go-chassis/control"
	"github.com/go-chassis/go-chassis/core/config"
)

var panelPlugin = make(map[string]func(options Options) control.Panel)

//DefaultPanelEgress get fetch config
var DefaultPanelEgress control.Panel

//InstallPlugin install implementation
func InstallPlugin(name string, f func(options Options) control.Panel) {
	panelPlugin[name] = f
}

//Options is options
type Options struct {
	Address string
}

//Init initialize DefaultPanel
func Init() error {
	infra := config.GlobalDefinition.Panel.Infra
	if infra == "" || infra == "archaius" {
		return nil
	}

	f, ok := panelPlugin[infra]
	if !ok {
		return fmt.Errorf("do not support [%s] panel", infra)
	}

	DefaultPanelEgress = f(Options{
		Address: config.GlobalDefinition.Panel.Settings["address"],
	})
	return nil
}

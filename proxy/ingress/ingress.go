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

package ingress

import (
	"errors"
	"fmt"
	"github.com/apache/servicecomb-mesher/proxy/config"
	"github.com/go-chassis/go-archaius"
)

//error in ingress package
var (
	ErrNotMatch = errors.New("no matching rule")
)
var plugin = make(map[string]func() (RuleFetcher, error))

//RuleFetcher query ingress rule
type RuleFetcher interface {
	Fetch(protocol, host, apiPath string, headers map[string][]string) (*config.IngressRule, error)
}

//DefaultFetcher fetch config
var DefaultFetcher RuleFetcher

//InstallPlugin install implementation
func InstallPlugin(name string, f func() (RuleFetcher, error)) {
	plugin[name] = f
}

//Init initialize
func Init() error {
	t := archaius.GetString("mesher.ingress.type", "servicecomb")
	f, ok := plugin[t]
	if !ok {
		return fmt.Errorf("do not support [%s]", t)
	}
	var err error
	DefaultFetcher, err = f()
	return err
}

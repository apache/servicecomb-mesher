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

package templates

import (
	"fmt"

	"github.com/go-mesh/openlogging"
	corev1 "k8s.io/api/core/v1"
)

var templateMap = map[string]func(configPath, templatePath string) (Templater, error){}

// Templater sidecar's configured template parser
type Templater interface {
	// PodSpecFromTemplate get PodSpec from template
	PodSpecFromTemplate(*corev1.Pod) (*corev1.PodSpec, error)

	// UpdateConfig update sidecar config
	UpdateConfig(configPath string) error

	// UpdateTemplate update sidecar template
	UpdateTemplate(tmplConfig string) error
}

// Register templater to manage
func Register(name string, fn func(configPath, templatePath string) (Templater, error)) {
	_, ok := templateMap[name]
	if ok {
		openlogging.Warn(fmt.Sprintf("sidecar already exists, name = %s", name))
		return
	}
	templateMap[name] = fn
}

// NewTemplater templater by name
func NewTemplater(name, configPath, templatePath string) (Templater, error) {
	fn, ok := templateMap[name]
	if !ok {
		err := fmt.Errorf("sidecar not found, , name = %s", name)
		openlogging.Error(err.Error())
		return nil, err
	}
	return fn(configPath, templatePath)
}

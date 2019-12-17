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

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"sync"
	"text/template"

	"github.com/apache/servicecomb-mesher/injection/templates"
	"github.com/go-mesh/openlogging"
	"gopkg.in/yaml.v2"
	corev1 "k8s.io/api/core/v1"
)

const (
	TemplateName          = "servicecomb-mesher"
	envVersion            = "VERSION"
	envListenPortName     = "LISTEN_PORT_NAME"
	defaultAppVersion     = "1.0.0"
	defaultListenPortName = "http"
)

// ServiceComb template
type ServiceComb struct {
	mu   sync.RWMutex
	conf *Config
	tmpl *template.Template
}

func init() {
	templates.Register(TemplateName, NewServiceComb)
}

// NewServiceComb returns ServiceComb Templater
func NewServiceComb(configPath, templatePath string) (templates.Templater, error) {
	sc := &ServiceComb{
		conf: DefaultConfig(),
		tmpl: DefaultTmpl(),
	}

	if configPath != "" {
		err := sc.UpdateConfig(configPath)
		if err != nil {
			return nil, err
		}
	}

	if templatePath != "" {
		err := sc.UpdateTemplate(templatePath)
		if err != nil {
			return nil, err
		}
	}
	return sc, nil
}

// UpdateConfig update ServiceComb configuration
func (s *ServiceComb) UpdateConfig(configPath string) error {
	conf, err := loadConfig(configPath)
	if err != nil {
		openlogging.Error(fmt.Sprintf("load sidecar config failed: filename = %s, error = %s", configPath, err.Error()))
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.conf = s.mergeConfig(conf, s.conf)
	return nil
}

func (s *ServiceComb) mergeConfig(dst, src *Config) *Config {
	if dst.Mesher.Name == "" {
		dst.Mesher.Name = src.Mesher.Name
	}

	if dst.Mesher.Image == "" {
		dst.Mesher.Image = src.Mesher.Image
	}

	if dst.Mesher.Tag == "" {
		dst.Mesher.Tag = src.Mesher.Tag
	}

	if dst.Mesher.GRPCPort <= 0 {
		dst.Mesher.GRPCPort = src.Mesher.GRPCPort
	}

	if dst.Mesher.HTTPPort <= 0 {
		dst.Mesher.HTTPPort = src.Mesher.HTTPPort
	}

	if dst.Mesher.AdminPort <= 0 {
		dst.Mesher.AdminPort = src.Mesher.AdminPort
	}

	if dst.ServiceCenter.Name == "" {
		dst.ServiceCenter.Name = src.ServiceCenter.Name
	}

	if dst.ServiceCenter.Namespace == "" {
		dst.ServiceCenter.Namespace = src.ServiceCenter.Namespace
	}

	protocol := "http"
	if dst.ServiceCenter.TlsConfig != nil {
		protocol = "https"
	}

	if dst.ServiceCenter.Address == "" {
		dst.ServiceCenter.Address = src.ServiceCenter.Address
		if dst.ServiceCenter.Address == "" {
			dst.ServiceCenter.Address = fmt.Sprintf("%s://%s.%s.svc.cluster.local:30100",
				protocol, dst.ServiceCenter.Name, dst.ServiceCenter.Namespace)
		}
	}
	return dst
}

// UpdateTemplate update ServiceComb template
func (s *ServiceComb) UpdateTemplate(tmplConfig string) error {
	tmpl, err := template.ParseFiles(tmplConfig)
	if err != nil {
		openlogging.Error(fmt.Sprintf("parse sidecar template failed: filename = %s, error = %s", tmplConfig, err.Error()))
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tmpl = tmpl
	return nil
}

// PodSpecFromTemplate returns PodSpec by template
func (s *ServiceComb) PodSpecFromTemplate(origin *corev1.Pod) (*corev1.PodSpec, error) {
	if len(origin.Spec.Containers) == 0 {
		err := fmt.Errorf("can not find containers in origin pod")
		openlogging.Error(err.Error())
		return nil, err
	}

	s.parseAppFromPod(origin)

	newPod, err := s.podFromTemplate()
	if err != nil {
		openlogging.Error("get pod from template failed: err = " + err.Error())
		return nil, err
	}

	return &newPod.Spec, nil
}

func (s *ServiceComb) parseAppFromPod(pod *corev1.Pod) {
	var port int32
	var version string
	name := pod.Name
	if name == "" {
		name = pod.GenerateName
	}

	least := len(name)
	for _, container := range pod.Spec.Containers {
		if strings.HasPrefix(name, container.Name) {
			if least > len(container.Name) {
				least = len(container.Name)
				name = container.Name
				version = getEnv(container.Env, envVersion, defaultAppVersion)
				listenPortMatch := getEnv(container.Env, envListenPortName, defaultListenPortName)
				port = getPort(container.Ports, listenPortMatch)
			}
		}
	}

	s.conf.App.Name = name
	s.conf.App.Host = "127.0.0.1"
	s.conf.App.Port = port
	s.conf.App.Version = version
}

func (s *ServiceComb) podFromTemplate() (*corev1.Pod, error) {
	buffer := &bytes.Buffer{}
	err := s.tmpl.Execute(buffer, s.conf)
	if err != nil {
		openlogging.Error("execute template failed: err = " + err.Error())
		return nil, err
	}

	pod := &corev1.Pod{}
	if err := yaml.Unmarshal(buffer.Bytes(), pod); err != nil {
		openlogging.Error("unmarshal template failed: err = " + err.Error())
		return nil, err
	}
	return pod, nil
}

func loadConfig(configPath string) (*Config, error) {
	content, err := ioutil.ReadFile(filepath.Clean(configPath))
	if err != nil {
		openlogging.Error(fmt.Sprintf("read config file failed: filename = %s, err = %s", configPath, err.Error()))
		return nil, err
	}

	conf := &Config{}
	err = yaml.Unmarshal(content, conf)
	if err != nil {
		openlogging.Error(fmt.Sprintf("unmarshal config file failed: filename = %s, err = %s", configPath, err.Error()))
		return nil, err
	}
	return conf, nil
}

func getPort(ports []corev1.ContainerPort, matchName string) (port int32) {
	for _, item := range ports {
		if item.Protocol == corev1.ProtocolTCP {
			port = item.ContainerPort
			if item.Name == matchName {
				return
			}
		}
	}
	return
}

func getEnv(envs []corev1.EnvVar, matchName string, noMatched string) string {
	for _, env := range envs {
		if strings.ToUpper(env.Name) == strings.ToUpper(matchName) {
			return env.Value
		}
	}
	return noMatched
}

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

package webhook

import (
	"crypto/tls"
	"sync"
)

// ConfigOption Config option
type ConfigOption func(*config)

// WithListenPort listen port to ConfigOption
func WithListenPort(port int) ConfigOption {
	return func(c *config) { c.port = port }
}

// WithCertFile tls cert file to ConfigOption
func WithCertFile(certFile string) ConfigOption {
	return func(c *config) { c.certFile = certFile }
}

// WithKeyFile tls key file to ConfigOption
func WithKeyFile(keyFile string) ConfigOption {
	return func(c *config) { c.keyFile = keyFile }
}

// WithSidecarConfig sidecar config to ConfigOption
func WithSidecarConfig(configPath string) ConfigOption {
	return func(c *config) { c.sidecarConfig = configPath }
}

// WithTemplateName template name to ConfigOption
func WithTemplateName(templateName string) ConfigOption {
	return func(c *config) { c.templateName = templateName }
}

// WithSidecarTemplate sidecar template to ConfigOption
func WithSidecarTemplate(templatePath string) ConfigOption {
	return func(c *config) { c.sidecarTemplate = templatePath }
}

func toConfig(opts ...ConfigOption) *config {
	conf := &config{}
	for _, opt := range opts {
		opt(conf)
	}
	return conf
}

type config struct {
	port            int
	certFile        string
	keyFile         string
	sidecarConfig   string
	sidecarTemplate string
	templateName    string

	mu   sync.RWMutex
	cert *tls.Certificate
}

func (c *config) loadTLSConfigFiles() error {
	pair, err := tls.LoadX509KeyPair(c.certFile, c.keyFile)
	if err != nil {
		return err
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cert = &pair
	return nil
}

func (c *config) getCertificate(*tls.ClientHelloInfo) (*tls.Certificate, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.cert, nil
}

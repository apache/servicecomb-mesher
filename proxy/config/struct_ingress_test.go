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

package config_test

import (
	"github.com/apache/servicecomb-mesher/proxy/config"
	"github.com/ghodss/yaml"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewRules(t *testing.T) {
	b := []byte(`
mesher:
  ingress:
    type: servicecomb
    rule:
      http: |
        - host: example.com
          limit: 30
          apiPath: /some/api
          service:
            name: example
            tags:
              version: 1.0.0
            redirectPath: /another/api
            port:
              name: http-legacy
              value: 8080
        - host: foo.com
          apiPath: /some/api
          service:
            name: foo
            tags:
              version: 1.0.0
            redirectPath: /another/api
            port:
              name: http
              value: 8080
`)
	c := &config.MesherConfig{}
	err := yaml.Unmarshal(b, c)
	assert.NoError(t, err)
	rules, err := config.NewRules(c.Mesher.Ingress.Rule["http"])
	assert.NoError(t, err)
	assert.Equal(t, 2, rules.Len())
	v := rules.Value()
	assert.Equal(t, "example", v[0].Service.Name)
}

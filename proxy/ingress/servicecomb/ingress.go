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

package servicecomb

import (
	"github.com/apache/servicecomb-mesher/proxy/config"
	"github.com/apache/servicecomb-mesher/proxy/ingress"
	"github.com/go-chassis/go-archaius"
	"github.com/patrickmn/go-cache"
	"regexp"
	"time"
)

const (
	cacheTTL       = 30
	ingressRuleKey = "mesher.ingress.rule.http"
)

var rulesData []*config.IngressRule

//IngressRuleFetcher query ingress rule
type IngressRuleFetcher struct {
	cache *cache.Cache
}

//Fetch get ingress rule
func (f *IngressRuleFetcher) Fetch(protocol, host, apiPath string, headers map[string][]string) (*config.IngressRule, error) {
	for _, r := range rulesData {
		if r.Host != "" && host != r.Host {
			//do not match host,then ignore path
			continue
		}
		match, err := regexp.MatchString(r.APIPath, apiPath)
		if err != nil {
			return nil, err
		}
		if match {
			return r, nil
		}
	}
	return nil, ingress.ErrNotMatch
}

func newFetcher() (ingress.RuleFetcher, error) {
	raw := archaius.GetString(ingressRuleKey, "")
	rules, err := config.NewRules(raw)
	if err != nil {
		return nil, err
	}

	err = archaius.RegisterListener(&ingressRuleEventListener{}, ingressRuleKey)
	if err != nil {
		return nil, err
	}
	rulesData = rules.Value()
	return &IngressRuleFetcher{
		cache: cache.New(cacheTTL*time.Second, 0),
	}, nil
}

func init() {
	ingress.InstallPlugin("servicecomb", newFetcher)
}

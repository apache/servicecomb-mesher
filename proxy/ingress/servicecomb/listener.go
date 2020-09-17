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
	"github.com/go-chassis/go-archaius/event"
	"github.com/go-chassis/go-chassis/v2/core/common"
	"github.com/go-chassis/openlog"
)

type ingressRuleEventListener struct{}

//Event update ingress rule
func (r *ingressRuleEventListener) Event(e *event.Event) {
	if e == nil {
		openlog.Warn("Event pointer is nil")
		return
	}
	openlog.Info("dark launch event", openlog.WithTags(openlog.Tags{
		"key":   e.Key,
		"event": e.EventType,
		"rule":  e.Value,
	}))
	raw, ok := e.Value.(string)
	if !ok {
		openlog.Error("invalid ingress rule", openlog.WithTags(openlog.Tags{
			"value": raw,
		}))
	}
	switch e.EventType {
	case common.Update:
		saveRules(raw)
	case common.Create:
		saveRules(raw)
	case common.Delete:
		rulesData = nil
		openlog.Info("ingress rule is removed", openlog.WithTags(
			openlog.Tags{
				"key": e.Key,
			}))
	}

}

func saveRules(raw string) {
	rules, err := config.NewRules(raw)
	if err != nil {
		openlog.Error("invalid ingress rule", openlog.WithTags(openlog.Tags{
			"value": raw,
		}))
	}
	rulesData = rules.Value()
	openlog.Info("update ingress rule", openlog.WithTags(openlog.Tags{
		"value": raw,
	}))
}

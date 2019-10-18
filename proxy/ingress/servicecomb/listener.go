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
	"github.com/go-chassis/go-chassis/core/common"
	"github.com/go-mesh/openlogging"
)

type ingressRuleEventListener struct{}

//Event update ingress rule
func (r *ingressRuleEventListener) Event(e *event.Event) {
	if e == nil {
		openlogging.Warn("Event pointer is nil")
		return
	}
	openlogging.Info("dark launch event", openlogging.WithTags(openlogging.Tags{
		"key":   e.Key,
		"event": e.EventType,
		"rule":  e.Value,
	}))
	raw, ok := e.Value.(string)
	if !ok {
		openlogging.Error("invalid ingress rule", openlogging.WithTags(openlogging.Tags{
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
		openlogging.Info("ingress rule is removed", openlogging.WithTags(
			openlogging.Tags{
				"key": e.Key,
			}))
	}

}

func saveRules(raw string) {
	rules, err := config.NewRules(raw)
	if err != nil {
		openlogging.Error("invalid ingress rule", openlogging.WithTags(openlogging.Tags{
			"value": raw,
		}))
	}
	rulesData = rules.Value()
	openlogging.Info("update ingress rule", openlogging.WithTags(openlogging.Tags{
		"value": raw,
	}))
}

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

package pilot

import (
	"errors"
	"fmt"
	"github.com/go-chassis/go-archaius/event"
	"github.com/go-chassis/go-archaius/source"
	"os"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/apache/servicecomb-mesher/proxy/config"
	"github.com/apache/servicecomb-mesher/proxy/control/istio"
	"github.com/apache/servicecomb-mesher/proxy/pkg/egress"
	istioinfra "github.com/apache/servicecomb-mesher/proxy/pkg/infras/istio"
	"github.com/go-chassis/go-chassis/core/lager"
)

const egressPilotSourceName = "EgressPilotSource"
const egressPilotSourcePriority = 8

//OUTBOUND outbound
const OUTBOUND = "outbound"

//clusteroriginaldst cluster type
const clusteroriginaldst = 4

// DefaultPilotRefresh is default pilot refresh time
// TODO: use stream instead
var DefaultPilotRefresh = 10 * time.Second

var pilotfetcher *source.Manager
var pilotChan = make(chan string, 10)

func setChanForPilot(k string) bool {
	select {
	case pilotChan <- k:
		return true
	default:
		return false
	}
}

// InitPilotFetcher init the config mgr and add several sources
func InitPilotFetcher(o egress.Options) error {
	// register and init pilot fetcher
	pilotfetcher = source.NewManager()
	pilotfetcher.RegisterListener(&pilotEventListener{}, ".*")
	return addEgressPilotSource(o)
}

// addEgressPilotSource adds a config source to pilotfetcher
func addEgressPilotSource(o egress.Options) error {
	if pilotfetcher == nil {
		return errors.New("pilotfetcher is nil, please init it first")
	}

	s, err := newPilotSource(o)
	if err != nil {
		return err
	}
	lager.Logger.Infof("New [%s] source success", s.GetSourceName())
	return pilotfetcher.AddSource(s)
}

// pilotSource keeps the egress rule in istio
type pilotSource struct {
	refreshInverval time.Duration
	fetcher         istioinfra.XdsClient
	priority        int
	mu              sync.RWMutex
	pmu             sync.RWMutex
	Configurations  map[string]interface{}
	PortToService   map[string]string
}

func newPilotSource(o egress.Options) (*pilotSource, error) {
	pilotAddr := o.Endpoints[0]
	nodeInfo := getNodeInfo()
	xdsClient, err := istioinfra.NewXdsClient(pilotAddr, nil, nodeInfo, o.ConfigPath)
	if err != nil {
		return nil, fmt.Errorf("connect to pilot failed: %v", err)
	}

	return &pilotSource{
		// TODO: read from config
		refreshInverval: DefaultPilotRefresh,
		Configurations:  map[string]interface{}{},
		PortToService:   map[string]string{},
		fetcher:         *xdsClient,
		priority:        egressPilotSourcePriority,
	}, nil
}

func getNodeInfo() *istioinfra.NodeInfo {
	PodName := os.Getenv("POD_NAME")
	PodNamespace := os.Getenv("POD_NAMESPACE")
	InstanceIP := os.Getenv("INSTANCE_IP")

	nodeInfo := &istioinfra.NodeInfo{
		PodName:    PodName,
		Namespace:  PodNamespace,
		InstanceIP: InstanceIP,
	}
	return nodeInfo
}

func (r *pilotSource) GetSourceName() string { return egressPilotSourceName }
func (r *pilotSource) GetPriority() int      { return r.priority }
func (r *pilotSource) Cleanup() error        { return nil }

//AddDimensionInfo no use
func (r *pilotSource) AddDimensionInfo(labels map[string]string) error {
	return nil
}

func (r *pilotSource) GetConfigurations() (map[string]interface{}, error) {
	egressConfigs, err := r.getEgressConfigFromPilot()
	if err != nil {
		lager.Logger.Error("Get egress config from pilot failed" + err.Error())
		return nil, err
	}
	d := make(map[string]interface{}, 0)
	d["pilotEgress"] = egressConfigs
	r.mu.Lock()
	r.Configurations = d
	r.mu.Unlock()
	return d, nil
}

func (r *pilotSource) GetConfigurationByKey(k string) (interface{}, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if value, ok := r.Configurations[k]; ok {
		return value, nil
	}
	return nil, fmt.Errorf("not found %s", k)
}

//Set no use
func (r *pilotSource) Set(key string, value interface{}) error {
	return nil
}

//Delete no use
func (r *pilotSource) Delete(key string) error {
	return nil
}

// get egress config from pilot
func (r *pilotSource) getEgressConfigFromPilot() ([]*config.EgressRule, error) {
	clusters, _ := r.fetcher.CDS()
	var egressRules []*config.EgressRule
	for _, cluster := range clusters {

		if cluster.Type == clusteroriginaldst {
			data := strings.Split(cluster.Name, "|")
			if len(data) > 1 && data[0] == OUTBOUND {
				intport, err := strconv.Atoi(data[1])
				if err != nil {
					return nil, nil
				}
				var rule config.EgressRule
				rule.Hosts = []string{data[len(data)-1]}
				rule.Ports = []*config.EgressPort{
					{Port: int32(intport),
						Protocol: "http"}}

				egressRules = append(egressRules, &rule)
			}

		}

	}

	return egressRules, nil
}

func (r *pilotSource) setPortForDestination(service, port string) {
	r.pmu.RLock()
	r.PortToService[port] = service
	r.pmu.RUnlock()
}

//SetPriority custom priority
func (r *pilotSource) SetPriority(priority int) {
	r.priority = priority
}
func (r *pilotSource) Watch(callback source.EventHandler) error {
	// Periodically refresh configurations
	ticker := time.NewTicker(r.refreshInverval)
	for {
		select {
		case <-pilotChan:
			data, err := r.GetConfigurations()
			if err != nil {
				lager.Logger.Error("pilot pull configuration error" + err.Error())
				continue
			}
			for k, d := range data {
				istio.SaveToEgressCache(map[string][]*config.EgressRule{k: d.([]*config.EgressRule)})
			}

		case <-ticker.C:
			data, err := r.refreshConfigurations()
			if err != nil {
				lager.Logger.Error("pilot refresh configuration error" + err.Error())
				continue
			}
			events, err := r.populateEvents(data)
			if err != nil {
				lager.Logger.Warnf("populate event error", err)
				return err
			}
			//Generate OnEvent Callback based on the events created
			lager.Logger.Debugf("event On receive %+v", events)
			for _, event := range events {
				callback.OnEvent(event)
			}
		}
	}
	return nil
}

func (r *pilotSource) refreshConfigurations() (map[string]interface{}, error) {
	data := make(map[string]interface{}, 0)

	egressConfigs, err := r.getEgressConfigFromPilot()
	if err != nil {
		lager.Logger.Error("Get egress config from pilot failed" + err.Error())
		return nil, err
	}

	data["pilotEgress"] = egressConfigs
	return data, nil
}

func (r *pilotSource) populateEvents(updates map[string]interface{}) ([]*event.Event, error) {
	events := make([]*event.Event, 0)
	new := make(map[string]interface{})

	// generate create and update event
	r.mu.RLock()
	current := r.Configurations
	r.mu.RUnlock()

	for key, value := range updates {
		new[key] = value
		currentValue, ok := current[key]
		if !ok { // if new configuration introduced
			events = append(events, constructEvent(event.Create, key, value))
		} else if !reflect.DeepEqual(currentValue, value) {
			events = append(events, constructEvent(event.Update, key, value))
		}
	}
	// generate delete event
	for key, value := range current {
		_, ok := new[key]
		if !ok { // when old config not present in new config
			events = append(events, constructEvent(event.Delete, key, value))
		}
	}

	// update with latest config
	r.mu.Lock()
	r.Configurations = new
	r.mu.Unlock()
	return events, nil
}

func constructEvent(eventType string, key string, value interface{}) *event.Event {
	return &event.Event{
		EventType:   eventType,
		EventSource: egressPilotSourceName,
		Value:       value,
		Key:         key,
	}
}

// pilotEventListener handle event dispatcher
type pilotEventListener struct{}

// update egress rule of a service
func (r *pilotEventListener) Event(e *event.Event) {
	if e == nil {
		lager.Logger.Warn("pilot event pointer is nil", nil)
		return
	}

	v := pilotfetcher.GetConfig(e.Key)
	if v == nil {
		istio.SaveToEgressCache(nil)
		return
	}
	egressRules, ok := v.([]*config.EgressRule)
	if !ok {
		lager.Logger.Error("value of pilot is not type []*EgressRule", nil)
		return
	}

	ok, _ = egress.ValidateEgressRule(map[string][]*config.EgressRule{e.Key: egressRules})
	if ok {
		istio.SaveToEgressCache(map[string][]*config.EgressRule{e.Key: egressRules})
		lager.Logger.Infof("Update [%s] egress rule of pilot success", e.Key)
	}
}

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

package health

import (
	"errors"
	"fmt"

	"github.com/apache/servicecomb-mesher/proxy/resource/v1/version"
	"github.com/go-chassis/go-chassis/core/config"
	"github.com/go-chassis/go-chassis/core/registry"
	"github.com/go-chassis/go-chassis/pkg/runtime"
	"github.com/go-mesh/openlogging"
)

//GetMesherHealth returns health
func GetMesherHealth() *Health {
	serviceName, version, err := getServiceStatus()
	if err != nil {
		openlogging.Error("health check failed: " + err.Error())
		resp := &Health{
			ServiceName: serviceName,
			Version:     version,
			Status:      Green,
			Error:       "",
		}
		resp.Status = Red
		resp.Error = err.Error()
		return resp
	}
	resp := &Health{
		ServiceName: serviceName,
		Version:     version,
		Status:      Green,
		Error:       "",
	}
	return resp
}

func getServiceStatus() (serviceName, v string, err error) {
	appID := runtime.App
	microServiceName := runtime.ServiceName
	v = runtime.Version
	if v == "" {
		v = version.DefaultVersion
	}
	environment := config.MicroserviceDefinition.ServiceDescription.Environment
	serviceID, err := registry.DefaultServiceDiscoveryService.GetMicroServiceID(appID, microServiceName, v, environment)
	if err != nil {
		return microServiceName, v, err
	}
	if len(serviceID) == 0 {
		return microServiceName, v, errors.New("serviceID is empty")
	}
	instances, err := registry.DefaultServiceDiscoveryService.GetMicroServiceInstances(serviceID, serviceID)
	if err != nil {
		return microServiceName, v, err
	}
	if len(instances) == 0 {
		return microServiceName, v, errors.New("no instance found")
	}
	for _, instance := range instances {
		ok, err := registry.DefaultRegistrator.Heartbeat(serviceID, instance.InstanceID)
		if err != nil {
			return microServiceName, v, err
		}
		if !ok {
			e := fmt.Errorf("heartbeat failed, instanceId: %s", instance.InstanceID)
			return microServiceName, v, e
		}
	}
	return microServiceName, v, nil
}

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

package v1

import (
	"github.com/apache/servicecomb-mesher/proxy/config"
	"github.com/go-chassis/go-chassis"
	"github.com/go-mesh/openlogging"
)

//RegisterWebService creates route and returns all admin APIs
func RegisterWebService() {
	chassis.RegisterSchema("rest-admin", &RouteResource{})
	chassis.RegisterSchema("rest-admin", &StatusResource{})
}

//Init function initiates admin API
func Init() (err error) {
	if !config.GetConfig().Admin.Enable {
		openlogging.Info("admin API is disabled")
		return nil
	}
	openlogging.Info("admin API is enabled")
	RegisterWebService()
	return
}

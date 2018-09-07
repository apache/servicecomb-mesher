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

package common

//Constants for default app and version
const (
	DefaultApp     = "default"
	DefaultVersion = "0.0.1"
)

//Constants for buildtag app and version
const (
	BuildInTagApp     = "app"
	BuildInTagVersion = "version"
)

//ComponentName is contant for component name
const ComponentName = "mesher"

//ModeSidecar is constant for side car mode
const ModeSidecar = "sidecar"

//ModePerHost is constant for side car mode
const ModePerHost = "per-host"

//Constants for env specific addr and service ports
const (
	//EnvSpecificAddr Deprecated
	EnvSpecificAddr = "SPECIFIC_ADDR"
	EnvServicePorts = "SERVICE_PORTS"
)

//HTTPProtocol is constant for protocol
const HTTPProtocol = "http"

//Constants for provider and consumer handlers
const (
	ChainConsumerOutgoing = "outgoing"
	ChainProviderIncoming = "incoming"
)

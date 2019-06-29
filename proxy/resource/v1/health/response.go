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

//StatusCode is type of string
type StatusCode string

const (
	//Red is a const
	Red StatusCode = "red"
	//Green is a const
	Green StatusCode = "green"
)

//Health has details about health of a service
type Health struct {
	ServiceName                 string     `json:"serviceName,omitempty"`
	Version                     string     `json:"version,omitempty"`
	Status                      StatusCode `json:"status,omitempty"`
	ConnectedConfigCenterClient bool       `json:"connectedConfigCenterClient"`
	ConnectedMonitoring         bool       `json:"connectedMonitoring"`
	Error                       string     `json:"error,omitempty"`
}

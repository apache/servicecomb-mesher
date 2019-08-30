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

package ports

import "github.com/go-chassis/go-chassis/core/common"

var defaultProtocolPort = map[string]string{
	common.ProtocolRest: "30101",
	"grpc":              "40101",
}

//SetFixedPort allows developer set a fixed port for for you protocol
func SetFixedPort(protocol, port string) {
	defaultProtocolPort[protocol] = port
}

//GetFixedPort return port pf a protocol
func GetFixedPort(protocol string) string {
	return defaultProtocolPort[protocol]

}

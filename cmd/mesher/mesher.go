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

package main

import (
	_ "net/http/pprof"

	_ "github.com/apache/servicecomb-mesher/proxy/resolver/authority"

	_ "github.com/apache/servicecomb-mesher/proxy/handler"
	//protocols
	_ "github.com/apache/servicecomb-mesher/proxy/protocol/dubbo/client/chassis"
	_ "github.com/apache/servicecomb-mesher/proxy/protocol/dubbo/server"
	_ "github.com/apache/servicecomb-mesher/proxy/protocol/dubbo/simpleRegistry"
	_ "github.com/apache/servicecomb-mesher/proxy/protocol/grpc"
	_ "github.com/apache/servicecomb-mesher/proxy/protocol/http"
	//ingress rule fetcher
	_ "github.com/apache/servicecomb-mesher/proxy/ingress/servicecomb"
	"github.com/apache/servicecomb-mesher/proxy/server"

	_ "github.com/apache/servicecomb-mesher/proxy/pkg/egress/archaius"

	_ "github.com/apache/servicecomb-mesher/proxy/handler/oauth2"

	//middle wares
	_ "github.com/go-chassis/go-chassis/v2/middleware/circuit"
	_ "github.com/go-chassis/go-chassis/v2/middleware/ratelimiter"
)

func main() {
	server.Run()
}

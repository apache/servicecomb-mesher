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

package authority

import (
	"fmt"
	"github.com/apache/servicecomb-mesher/proxy/resolver"
	"github.com/go-mesh/openlogging"
	"strings"
)

//GRPCDefaultDestinationResolver is a struct
type GRPCDefaultDestinationResolver struct {
}

//Resolve resolves service name
func (dr *GRPCDefaultDestinationResolver) Resolve(sourceAddr, host, rawURI string, header map[string]string) (string, string, error) {
	s := strings.Split(rawURI, ":")
	if len(s) != 2 {
		err := fmt.Errorf("can not parse [%s]", rawURI)
		openlogging.Error(err.Error())
		return "", "", err
	}

	return s[0], s[1], nil
}

//New return return dr
func New() resolver.DestinationResolver {
	return &GRPCDefaultDestinationResolver{}
}

func init() {
	resolver.InstallDestinationResolverPlugin("authority", New)
	resolver.SetDefaultDestinationResolver("grpc", &GRPCDefaultDestinationResolver{})
}

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
	"github.com/go-chassis/go-chassis/core/lager"
	"github.com/go-mesh/mesher/resolver"
	"strings"
)

//GPRCDefaultDestinationResolver is a struct
type GPRCDefaultDestinationResolver struct {
}

//Resolve resolves service's endpoint
//service may have multiple port for same protocol
func (dr *GPRCDefaultDestinationResolver) Resolve(sourceAddr string, header map[string]string, rawURI string, destinationName *string) (string, error) {
	s := strings.Split(rawURI, ":")
	if len(s) != 2 {
		err := fmt.Errorf("can not parse [%s]", rawURI)
		lager.Logger.Error(err.Error())
		return "", err
	}

	*destinationName = s[0]
	return s[1], nil
}

//New return return dr
func New() resolver.DestinationResolver {
	return &GPRCDefaultDestinationResolver{}
}

func init() {
	resolver.InstallDestinationResolverPlugin("authority", New)
	resolver.SetDefaultDestinationResolver("grpc", &GPRCDefaultDestinationResolver{})
}

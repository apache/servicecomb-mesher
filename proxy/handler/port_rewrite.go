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

package handler

import (
	"fmt"
	"strings"

	"github.com/apache/servicecomb-mesher/proxy/pkg/ports"
	"github.com/go-chassis/go-chassis/core/handler"
	"github.com/go-chassis/go-chassis/core/invocation"
	"github.com/go-mesh/openlogging"
)

//PortMapForPilot is a constant
const PortMapForPilot = "port-selector"

//PortSelectionHandler ..
type PortSelectionHandler struct {
}

//Handle function replace the provider port to mesher port so that traffic goes through mesher
func (ps *PortSelectionHandler) Handle(chain *handler.Chain, inv *invocation.Invocation, cb invocation.ResponseCallBack) {
	var err error
	inv.Endpoint, err = replacePort(inv.Protocol, inv.Endpoint)
	if err != nil {
		openlogging.Error("can not replace port: " + err.Error())
	}
	if inv.Endpoint == "" {
		r := &invocation.Response{
			Err: err,
		}

		if err := cb(r); err != nil {
			openlogging.Error("Response callBack error: " + err.Error())
		}
		return
	}

	chain.Next(inv, func(r *invocation.Response) error {
		return cb(r)
	})
}

//replacePort will replace the provider port with mesher port.
func replacePort(protocol, endpoint string) (string, error) {
	eps := strings.Split(endpoint, ":")
	if len(eps) != 2 {
		return "", fmt.Errorf("invalid endpoint [%s]", eps)
	}

	eps[1] = ports.GetFixedPort(protocol)

	return strings.Join(eps, ":"), nil
}

//Name returns name
func (ps *PortSelectionHandler) Name() string {
	return PortMapForPilot
}

//New create new port for pilot handler and retuns
func New() handler.Handler {
	return &PortSelectionHandler{}
}

func init() {
	handler.RegisterHandler(PortMapForPilot, New)
}

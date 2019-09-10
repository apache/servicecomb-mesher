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

package cmd

import (
	"fmt"
	"github.com/apache/servicecomb-mesher/proxy/common"
	chassiscommon "github.com/go-chassis/go-chassis/core/common"
	"github.com/urfave/cli"
	"log"
	"os"
	"strings"
)

//Local is a constant
const Local = "127.0.0.1"

//ConfigFromCmd store cmd params
type ConfigFromCmd struct {
	ConfigFile        string
	Role              string
	LocalServicePorts string
	PortsMap          map[string]string
}

//Configs is a pointer of struct ConfigFromCmd
var Configs *ConfigFromCmd

// parseConfigFromCmd
func parseConfigFromCmd(args []string) (err error) {
	app := cli.NewApp()
	app.HideVersion = true
	app.Usage = "a service mesh that governance your service traffic."
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "config",
			Usage:       "mesher config file, example: --config=mesher.yaml",
			Destination: &Configs.ConfigFile,
		},
		cli.StringFlag{
			Name:  "mode",
			Value: common.RoleSidecar,
			Usage: fmt.Sprintf("mesher role [ %s|%s|%s ]",
				common.RolePerHost, common.RoleSidecar, common.RoleEdge),
			Destination: &Configs.Role,
		},
		cli.StringFlag{
			Name:        "service-ports",
			EnvVar:      common.EnvServicePorts,
			Usage:       fmt.Sprintf("service protocol and port,examples: --service-ports=http:3000,grpc:8000"),
			Destination: &Configs.LocalServicePorts,
		},
	}
	app.Action = func(c *cli.Context) error {
		return nil
	}

	err = app.Run(args)
	return
}

//Init get config and parses those command
func Init() error {
	Configs = &ConfigFromCmd{}
	return parseConfigFromCmd(os.Args)
}

//GeneratePortsMap generates ports map
func (c *ConfigFromCmd) GeneratePortsMap() error {
	c.PortsMap = make(map[string]string)
	if c.LocalServicePorts != "" { //parse service ports
		s := strings.Split(c.LocalServicePorts, ",")
		for _, v := range s {
			p := strings.Split(v, ":")
			if len(p) != 2 {
				return fmt.Errorf("[%s] is invalid", p)
			}
			c.PortsMap[p[0]] = Local + ":" + p[1]
		}
		return nil
	}
	//support deprecated env
	addr := os.Getenv(common.EnvSpecificAddr)
	if addr != "" {
		addr = strings.TrimSpace(addr)
		log.Printf("%s is deprecated, plz use SERVICE_PORTS=http:8080,grpc:9000 instead", common.EnvSpecificAddr)
		s := strings.Split(addr, ":")
		if len(s) != 2 {
			return fmt.Errorf("[%s] is invalid", addr)
		}
		c.PortsMap[chassiscommon.ProtocolRest] = Local + ":" + s[1]
	}

	return nil
}

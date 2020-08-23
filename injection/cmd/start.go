/*
 *  Licensed to the Apache Software Foundation (ASF) under one or more
 *  contributor license agreements.  See the NOTICE file distributed with
 *  this work for additional information regarding copyright ownership.
 *  The ASF licenses this file to You under the Apache License, Version 2.0
 *  (the "License"); you may not use this file except in compliance with
 *  the License.  You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 */

package cmd

import (
	"github.com/apache/servicecomb-mesher/injection/config"
	"github.com/apache/servicecomb-mesher/injection/server"
	"github.com/go-mesh/openlogging"
	"github.com/urfave/cli"
)

var conf = config.DefaultConfig()

// GetCmdStart return service stated configure
func GetCmdStart() cli.Command {
	return cli.Command{
		Name:   "start",
		Usage:  "start sidecar injection server",
		Action: start,
		Flags: []cli.Flag{
			cli.IntFlag{
				Name:        "port",
				Usage:       "webhook server port",
				Value:       conf.Port,
				Destination: &conf.Port,
			},
			cli.StringFlag{
				Name:        "tlsCertPath",
				Usage:       "path of tls cert",
				Value:       conf.CertFile,
				Destination: &conf.CertFile,
			},
			cli.StringFlag{
				Name:        "tlsKeyPath",
				Usage:       "path of tls key",
				Value:       conf.KeyFile,
				Destination: &conf.KeyFile,
			},
			cli.StringFlag{
				Name:        "templateName",
				Usage:       "template name of sidecar",
				Value:       conf.TemplateName,
				Destination: &conf.TemplateName,
			},
			cli.StringFlag{
				Name:        "configPath",
				Usage:       "path of sidecar config",
				Destination: &conf.SidecarConfig,
			},
			cli.StringFlag{
				Name:        "templatePath",
				Usage:       "path of sidecar template",
				Destination: &conf.SidecarTemplate,
			},
		},
	}
}

func start(c *cli.Context) {
	if err := server.Run(conf); err != nil {
		openlogging.Error("start webhook server failed:" + err.Error())
		panic(err)
	}
}

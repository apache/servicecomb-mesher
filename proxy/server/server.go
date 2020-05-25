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

package server

import (
	"net/http"

	"github.com/apache/servicecomb-mesher/proxy/bootstrap"
	"github.com/apache/servicecomb-mesher/proxy/cmd"
	"github.com/apache/servicecomb-mesher/proxy/config"
	"github.com/apache/servicecomb-mesher/proxy/health"

	"github.com/apache/servicecomb-mesher/proxy/resource/v1/version"
	"github.com/go-chassis/go-chassis"
	"github.com/go-mesh/openlogging"
)

// Run run mesher proxy server
func Run() {
	// server init
	if err := cmd.Init(); err != nil {
		openlogging.Fatal(err.Error())
	}
	if err := cmd.Configs.GeneratePortsMap(); err != nil {
		openlogging.Fatal(err.Error())
	}
	bootstrap.RegisterFramework()
	bootstrap.SetHandlers()
	if err := chassis.Init(); err != nil {
		openlogging.Fatal("Go chassis init failed, Mesher is not available: " + err.Error())
	}
	if err := bootstrap.InitEgressChain(); err != nil {
		openlogging.Error("egress chain int failed: %s", openlogging.WithTags(openlogging.Tags{
			"err": err.Error(),
		}))
		openlogging.Fatal(err.Error())
	}

	if err := bootstrap.Start(); err != nil {
		openlogging.Fatal("Bootstrap failed: " + err.Error())
	}
	openlogging.Info("server start complete", openlogging.WithTags(openlogging.Tags{
		"version": version.Ver().Version,
	}))
	if err := health.Run(); err != nil {
		openlogging.Fatal("Health manager start failed: " + err.Error())
	}
	profile()
	if err := chassis.Run(); err != nil {
		openlogging.Fatal("Chassis failed: " + err.Error())
	}
}

func profile() {
	if config.GetConfig().PProf != nil {
		if config.GetConfig().PProf.Enable {
			go startProfiling()
		}
	}
}

func startProfiling() {
	if config.GetConfig().PProf.Listen == "" {
		config.GetConfig().PProf.Listen = "127.0.0.1:6060"
	}
	openlogging.Warn("Enable pprof on " + config.GetConfig().PProf.Listen)
	if err := http.ListenAndServe(config.GetConfig().PProf.Listen, nil); err != nil {
		openlogging.Error("Can not enable pprof: " + err.Error())
	}
}

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
		panic(err)
	}
	if err := cmd.Configs.GeneratePortsMap(); err != nil {
		panic(err)
	}
	bootstrap.RegisterFramework()
	bootstrap.SetHandlers()
	if err := chassis.Init(); err != nil {
		openlogging.Error("Go chassis init failed, Mesher is not available: " + err.Error())
		panic(err)
	}
	if err := bootstrap.InitEgressChain(); err != nil {
		openlogging.Error("egress chain int failed: %s", openlogging.WithTags(openlogging.Tags{
			"err": err.Error(),
		}))
		panic(err)
	}

	if err := bootstrap.Start(); err != nil {
		openlogging.Error("Bootstrap failed: " + err.Error())
		panic(err)
	}
	openlogging.Info("server start complete", openlogging.WithTags(openlogging.Tags{
		"version": version.Ver().Version,
	}))
	if err := health.Run(); err != nil {
		openlogging.Error("Health manager start failed: " + err.Error())
		panic(err)
	}
	profile()
	chassis.Run()
}

func profile() {
	if config.GetConfig().PProf != nil {
		if config.GetConfig().PProf.Enable {
			go func() {
				if config.GetConfig().PProf.Listen == "" {
					config.GetConfig().PProf.Listen = "127.0.0.1:6060"
				}
				openlogging.Warn("Enable pprof on " + config.GetConfig().PProf.Listen)
				if err := http.ListenAndServe(config.GetConfig().PProf.Listen, nil); err != nil {
					openlogging.Error("Can not enable pprof: " + err.Error())
				}
			}()
		}
	}
}

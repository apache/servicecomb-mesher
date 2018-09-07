package server

import (
	"net/http"

	"github.com/go-mesh/mesher/adminapi/version"
	"github.com/go-mesh/mesher/bootstrap"
	"github.com/go-mesh/mesher/cmd"
	"github.com/go-mesh/mesher/config"
	"github.com/go-mesh/mesher/health"

	"github.com/go-chassis/go-chassis"
	"github.com/go-chassis/go-chassis/core/lager"
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
		lager.Logger.Error("Go chassis init failed, Mesher is not available: " + err.Error())
		panic(err)
	}
	if err := bootstrap.Start(); err != nil {
		lager.Logger.Error("Bootstrap failed: " + err.Error())
		panic(err)
	}
	lager.Logger.Infof("Version is %s", version.Ver().Version)
	if err := health.Run(); err != nil {
		lager.Logger.Error("Health manager start failed: " + err.Error())
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
				lager.Logger.Warn("Enable pprof on "+config.GetConfig().PProf.Listen, nil)
				if err := http.ListenAndServe(config.GetConfig().PProf.Listen, nil); err != nil {
					lager.Logger.Error("Can not enable pprof: " + err.Error())
				}
			}()
		}
	}
}

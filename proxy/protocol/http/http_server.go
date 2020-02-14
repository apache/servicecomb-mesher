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

package http

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/apache/servicecomb-mesher/proxy/ingress"
	"github.com/apache/servicecomb-mesher/proxy/pkg/runtime"
	"net"
	"net/http"
	"strings"

	"github.com/apache/servicecomb-mesher/proxy/common"
	"github.com/apache/servicecomb-mesher/proxy/resolver"
	chassisCom "github.com/go-chassis/go-chassis/core/common"
	"github.com/go-chassis/go-chassis/core/lager"
	"github.com/go-chassis/go-chassis/core/server"
	chassisTLS "github.com/go-chassis/go-chassis/core/tls"
	chassisRuntime "github.com/go-chassis/go-chassis/pkg/runtime"
	"github.com/go-mesh/openlogging"
)

const (
	Name = "http"
)

func init() {
	server.InstallPlugin(Name, newServer)
}

func newServer(opts server.Options) server.ProtocolServer {
	return &httpServer{
		opts: opts,
	}
}

type httpServer struct {
	opts   server.Options
	server *http.Server
}

func (hs *httpServer) Register(schema interface{}, options ...server.RegisterOption) (string, error) {
	return "", nil
}

func (hs *httpServer) Start() error {
	host, port, err := net.SplitHostPort(hs.opts.Address)
	if err != nil {
		return err
	}
	ip := net.ParseIP(host)
	if ip == nil {
		return fmt.Errorf("IP format error, input is [%s]", hs.opts.Address)
	}
	if ip.To4() == nil {
		return fmt.Errorf("only support ipv4, input is [%s]", hs.opts.Address)
	}

	switch runtime.Role {
	case common.RoleSidecar:
		err = hs.startSidecar(host, port)
	default:
		err = hs.startCommonProxy()
	}
	if err != nil {
		return err
	}
	return nil
}

func (hs *httpServer) startSidecar(host, port string) error {
	mesherTLSConfig, mesherSSLConfig, mesherErr := chassisTLS.GetTLSConfigByService(
		common.ComponentName, "", chassisCom.Provider)
	if mesherErr != nil {
		if !chassisTLS.IsSSLConfigNotExist(mesherErr) {
			return mesherErr
		}
	} else {
		sslTag := genTag(common.ComponentName, chassisCom.Provider)
		lager.Logger.Warnf("%s TLS mode, verify peer: %t, cipher plugin: %s.",
			sslTag, mesherSSLConfig.VerifyPeer, mesherSSLConfig.CipherPlugin)
	}

	err := hs.listenAndServe("127.0.0.1"+":"+port, mesherTLSConfig, http.HandlerFunc(LocalRequestHandler))
	if err != nil {
		return err
	}
	resolver.SelfEndpoint = "127.0.0.1" + ":" + port

	switch host {
	case "0.0.0.0":
		return errors.New("in sidecar mode, forbidden to listen on 0.0.0.0")
	case "127.0.0.1":
		lager.Logger.Warnf("Mesher listen on 127.0.0.1, it can only proxy for consumer. " +
			"for provider, mesher must listen on external ip.")
		return nil
	default:
		serverTLSConfig, serverSSLConfig, serverErr := chassisTLS.GetTLSConfigByService(
			chassisRuntime.ServiceName, chassisCom.ProtocolRest, chassisCom.Provider)
		if serverErr != nil {
			if !chassisTLS.IsSSLConfigNotExist(serverErr) {
				return serverErr
			}
		} else {
			sslTag := genTag(chassisRuntime.ServiceName, chassisCom.ProtocolRest, chassisCom.Provider)
			lager.Logger.Warnf("%s TLS mode, verify peer: %t, cipher plugin: %s.",
				sslTag, serverSSLConfig.VerifyPeer, serverSSLConfig.CipherPlugin)
		}

		err = hs.listenAndServe(hs.opts.Address, serverTLSConfig, http.HandlerFunc(RemoteRequestHandler))
		if err != nil {
			return err
		}
	}

	return nil
}

func (hs *httpServer) startCommonProxy() error {
	if err := ingress.Init(); err != nil {
		return err
	}
	mesherTLSConfig, mesherSSLConfig, err := chassisTLS.GetTLSConfigByService(
		chassisRuntime.ServiceName, "rest", chassisCom.Provider)
	if err != nil {
		if !chassisTLS.IsSSLConfigNotExist(err) {
			return err
		}
	} else {
		lager.Logger.Warnf("TLS mode, verify peer: %t, cipher plugin: %s.",
			mesherSSLConfig.VerifyPeer, mesherSSLConfig.CipherPlugin)
	}

	err = hs.listenAndServe(hs.opts.Address, mesherTLSConfig, http.HandlerFunc(HandleIngressTraffic))
	if err != nil {
		return err
	}
	return nil
}

func (hs *httpServer) listenAndServe(addr string, t *tls.Config, h http.Handler) error {
	ln, err := net.Listen("tcp4", addr)
	if err != nil {
		return err
	}
	if t != nil {
		openlogging.Info("run as https")
		lnTLS := tls.NewListener(ln, t)
		ln = lnTLS
	}
	go func() {
		hs.server = &http.Server{
			Handler: h,
		}
		if err := hs.server.Serve(ln); err != nil {
			server.ErrRuntime <- err
			return
		}
	}()
	return nil
}

func (hs *httpServer) Stop() error {
	//go 1.8+ drain connections handleIncomingTraffic stop server
	if hs.server == nil {
		openlogging.Info("http server don't need to be stopped")
		return nil
	}
	if err := hs.server.Shutdown(context.TODO()); err != nil {
		panic(err)
	}
	openlogging.Info("Mesher gracefully stopped")
	return nil
}

func (hs *httpServer) String() string {
	return Name
}

func genTag(s ...string) string {
	return strings.Join(s, ".")
}

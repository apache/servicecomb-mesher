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

package grpc

import (
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/go-mesh/mesher/common"
	"github.com/go-mesh/mesher/config"
	"github.com/go-mesh/mesher/resolver"
	"net"
	"strings"

	chassisCom "github.com/go-chassis/go-chassis/core/common"
	chassisConfig "github.com/go-chassis/go-chassis/core/config"
	"github.com/go-chassis/go-chassis/core/lager"
	"github.com/go-chassis/go-chassis/core/server"
	chassisTLS "github.com/go-chassis/go-chassis/core/tls"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	protocolName = "grpc"
)

func init() {
	server.InstallPlugin(protocolName, newServer)
}

type grpcServer struct {
	opts   server.Options
	server *grpc.Server
}

func newServer(opts server.Options) server.ProtocolServer {
	return &grpcServer{
		opts: opts,
	}
}

func (hs *grpcServer) Register(schema interface{}, options ...server.RegisterOption) (string, error) {
	return "", nil
}

func (hs *grpcServer) Start() error {
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

	switch config.Mode {
	case common.ModeSidecar:
		err = hs.startSidecar(host, port)
	case common.ModePerHost:
		err = errors.New("do not support per host")
	}
	if err != nil {
		return err
	}
	return nil
}

func (hs *grpcServer) startSidecar(host, port string) error {
	mesherTLSConfig, mesherSSLConfig, mesherErr := chassisTLS.GetTLSConfigByService(
		common.ComponentName, protocolName, chassisCom.Provider)
	if mesherErr != nil {
		if !chassisTLS.IsSSLConfigNotExist(mesherErr) {
			return mesherErr
		}
	} else {
		sslTag := genTag(common.ComponentName, chassisCom.Provider)
		lager.Logger.Warnf("%s TLS mode, verify peer: %t, cipher plugin: %s.",
			sslTag, mesherSSLConfig.VerifyPeer, mesherSSLConfig.CipherPlugin)
	}

	if err := hs.listenAndServe("127.0.0.1"+":"+port, mesherTLSConfig, LocalRequestHandler); err != nil {
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
			chassisConfig.SelfServiceName, protocolName, chassisCom.Provider)
		if serverErr != nil {
			if !chassisTLS.IsSSLConfigNotExist(serverErr) {
				return serverErr
			}
		} else {
			sslTag := genTag(chassisConfig.SelfServiceName, chassisCom.ProtocolRest, chassisCom.Provider)
			lager.Logger.Warnf("%s TLS mode, verify peer: %t, cipher plugin: %s.",
				sslTag, serverSSLConfig.VerifyPeer, serverSSLConfig.CipherPlugin)
		}
		if err := hs.listenAndServe(hs.opts.Address, serverTLSConfig, RemoteRequestHandler); err != nil {
			return err
		}
	}

	return nil
}

func (hs *grpcServer) listenAndServe(addr string, tlsConfig *tls.Config, h grpc.StreamHandler) error {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	go func() {

		if tlsConfig != nil {
			lager.Logger.Info(protocolName + " enable TLS and listen on " + addr)
			hs.server = grpc.NewServer(grpc.Creds(credentials.NewTLS(tlsConfig)),
				grpc.UnknownServiceHandler(h),
				grpc.CustomCodec(&codec{}))
		} else {
			lager.Logger.Info(protocolName + " listen on " + addr)
			hs.server = grpc.NewServer(grpc.UnknownServiceHandler(h), grpc.CustomCodec(&codec{}))
		}

		if err := hs.server.Serve(ln); err != nil {
			server.ErrRuntime <- err
			return
		}
	}()
	return nil
}

func (hs *grpcServer) Stop() error {
	//go 1.8+ drain connections before stop server
	if hs.server == nil {
		lager.Logger.Info("grpc server doesn't need to be stopped")
		return nil
	}
	hs.server.GracefulStop()
	lager.Logger.Info("grpc server gracefully stopped")
	return nil
}

func (hs *grpcServer) String() string {
	return protocolName
}

func genTag(s ...string) string {
	return strings.Join(s, ".")
}

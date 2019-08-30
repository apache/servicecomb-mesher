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

package health

import (
	"errors"
	"github.com/apache/servicecomb-mesher/proxy/config"
	"github.com/go-chassis/go-chassis/core/lager"
	"github.com/go-chassis/go-chassis/core/registry"
	"github.com/go-chassis/go-chassis/pkg/runtime"
	"github.com/go-mesh/openlogging"
	"net"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	DefaultInterval = time.Second * 30
	DefaultTimeout  = time.Second * 10
)

//Error definitions
var (
	ErrPortEmpty  = errors.New("port is empty")
	ErrInvalidURI = errors.New("uri must start with /")
)

//Deal handle the unhealthy status
type Deal func(err error)

//L7Check is the interface for L7 checker
type L7Check func(check *config.HealthCheck, address string) error

//l7Checks save l7 check func
var l7Checks = make(map[string]L7Check)

//InstallChecker install checkers
func InstallChecker(n string, c L7Check) {
	l7Checks[n] = c
}

//UpdateInstanceStatus update status in registrator, it just works in client side discovery
func UpdateInstanceStatus(err error) {
	if registry.DefaultRegistrator == nil {
		lager.Logger.Warn("Registrator is nil, can not update instance status")
		return
	}
	if err != nil {
		if runtime.InstanceStatus == runtime.StatusRunning {
			lager.Logger.Info("service is not healthy, update status")
			ChangeStatus(runtime.StatusDown)
		}

	} else {
		if runtime.InstanceStatus == runtime.StatusDown {
			lager.Logger.Info("service is healthy, update status")
			ChangeStatus(runtime.StatusRunning)
		}
	}

}

//ChangeStatus change status in local and remote
func ChangeStatus(status string) {
	if err := registry.DefaultRegistrator.UpdateMicroServiceInstanceStatus(runtime.ServiceID, runtime.InstanceID, status); err != nil {
		lager.Logger.Error("update instance status failed:" + err.Error())
		return
	}
	runtime.InstanceStatus = status
	lager.Logger.Info("update instance status to: " + runtime.InstanceStatus)
}

//runCheckers run check routines
func runCheckers(c *config.HealthCheck, l7check L7Check, address string, deal Deal) (err error) {
	var interval = DefaultInterval
	if c.Interval != "" {
		interval, err = time.ParseDuration(c.Interval)
		if err != nil {
			return err
		}
	}
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			err := CheckService(c, l7check, address)
			if err != nil {
				lager.Logger.Errorf("health check failed for service port[%s]: %s", c.Port, err)
			}
			deal(err)
		}
	}()
	return nil
}

//CheckService check service health based on config
func CheckService(c *config.HealthCheck, l7check L7Check, address string) error {
	lager.Logger.Debugf("check port [%s]", c.Port)
	if l7check != nil {
		if err := l7check(c, address); err != nil {
			return err
		}
	} else {
		if err := L4Check(address); err != nil {
			return err
		}
	}
	lager.Logger.Debug("service is healthy: " + address)
	return nil
}

//L4Check check tcp connection
func L4Check(address string) error {
	c, err := net.DialTimeout("tcp", address, DefaultTimeout)
	if err != nil {
		return err
	}
	if err = c.Close(); err != nil {
		return err
	}
	return nil
}

//Run Launch go routines to check service health
func Run() error {
	openlogging.Info("local health manager start")
	for _, v := range config.GetConfig().HealthCheck {
		lager.Logger.Debugf("check local health [%s],protocol [%s]", v.Port, v.Protocol)
		address, check, err := ParseConfig(v)
		if err != nil {
			lager.Logger.Warn("Health keeper can not check health")
			return err
		}
		//TODO make pluggable Deal
		if err := runCheckers(v, check, address, UpdateInstanceStatus); err != nil {
			return err
		}
	}
	return nil
}

//ParseConfig validate config and return address, checker
//port name must not be empty
//port name must named as {protocol}-{name}
//protocol must has checker
func ParseConfig(c *config.HealthCheck) (string, L7Check, error) {
	if c.Port == "" {
		return "", nil, ErrPortEmpty
	}
	var check L7Check
	if c.Protocol != "" {
		var ok bool
		check, ok = l7Checks[c.Protocol]
		if !ok {
			return "", nil, errors.New("don not support L7 checker:" + c.Protocol)
		}
	} else {
		check = nil
	}

	address := "127.0.0.1:" + c.Port
	if c.URI != "" {
		if !strings.HasPrefix(c.URI, "/") {
			return "", nil, ErrInvalidURI
		}
	}
	if c.Match != nil {
		if c.Match.Status != "" {
			_, err := strconv.Atoi(c.Match.Status)
			if err != nil {
				return "", nil, err
			}
		}
		if c.Match.Body != "" {
			_, err := regexp.Compile(c.Match.Body)
			if err != nil {
				return "", nil, err
			}
		}
	}
	return address, check, nil
}

func init() {
	InstallChecker("rest", HTTPCheck)
}

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

package health_test

import (
	"context"
	"github.com/apache/servicecomb-mesher/proxy/config"
	"github.com/apache/servicecomb-mesher/proxy/health"
	"github.com/go-chassis/go-chassis/core/lager"
	"github.com/stretchr/testify/assert"
	"net"
	"net/http"
	"testing"
	"time"
)

func init() {
	lager.Init(&lager.Options{LoggerLevel: "DEBUG"})
}
func TestHttpCheck(t *testing.T) {
	server := &http.Server{
		Addr: "127.0.0.1:3000",
	}
	http.HandleFunc("/health", func(resp http.ResponseWriter, req *http.Request) {
		resp.WriteHeader(200)
		resp.Write([]byte("hello"))
	})

	t.Log("Check server stoped")
	check := &config.HealthCheck{
		Port: "rest-console",
		URI:  "/health",
	}
	addr := "127.0.0.1:3000"
	err := health.HTTPCheck(check, addr)
	assert.Error(t, err)

	t.Log("launch server")
	go server.ListenAndServe()
	time.Sleep(3 * time.Second)
	defer server.Shutdown(context.Background())

	t.Log("check real health ")
	check = &config.HealthCheck{
		Port: "rest-console",
		URI:  "/health",
	}
	addr = "127.0.0.1:3000"
	err = health.HTTPCheck(check, addr)
	assert.NoError(t, err)

	t.Log("status match 500,must fail ")
	check = &config.HealthCheck{
		Port: "rest-console",
		URI:  "/health",
		Match: &config.Match{
			Status: "201",
		},
	}
	addr = "127.0.0.1:3000"
	err = health.HTTPCheck(check, addr)
	assert.Error(t, err)

	t.Log("body match fake,must fail ")
	check = &config.HealthCheck{
		Port: "rest-console",
		URI:  "/health",
		Match: &config.Match{
			Body: "fake",
		},
	}
	addr = "127.0.0.1:3000"
	err = health.HTTPCheck(check, addr)
	assert.Error(t, err)

	t.Log("body match right,no error ")
	check = &config.HealthCheck{
		Port: "rest-console",
		URI:  "/health",
		Match: &config.Match{
			Body: "hello",
		},
	}
	addr = "127.0.0.1:3000"
	err = health.HTTPCheck(check, addr)
	assert.NoError(t, err)

	t.Log("all match,no error ")
	check = &config.HealthCheck{
		Port: "rest-console",
		URI:  "/health",
		Match: &config.Match{
			Status: "200",
			Body:   "hello",
		},
	}
	addr = "127.0.0.1:3000"
	err = health.HTTPCheck(check, addr)
	assert.NoError(t, err)
}

func TestParseConfig(t *testing.T) {
	check := &config.HealthCheck{
		Port: "8080",
		URI:  "/health",
		Match: &config.Match{
			Body: "hello",
		},
	}
	addr, c, err := health.ParseConfig(check)
	assert.NoError(t, err)
	assert.Equal(t, "127.0.0.1:8080", addr)
	assert.Nil(t, c)

	check = &config.HealthCheck{
		Port:     "8080",
		Protocol: "rest",
		URI:      "/health",
		Match: &config.Match{
			Body: "hello",
		},
	}
	_, c, err = health.ParseConfig(check)
	err = c(check, addr)
	assert.Error(t, err)
}
func TestL4Check(t *testing.T) {
	err := health.L4Check("127.0.0.1:3000")
	assert.Error(t, err)
	net.Listen("tcp", "127.0.0.1:3000")
	err = health.L4Check("127.0.0.1:3000")
	assert.NoError(t, err)
}

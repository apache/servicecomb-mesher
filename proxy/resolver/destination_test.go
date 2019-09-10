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

package resolver

import (
	"testing"

	//	"github.com/go-chassis/go-chassis/core/archaius"
	//	cConfig "github.com/go-chassis/go-chassis/core/config"
	//	"github.com/go-chassis/go-chassis/core/lager"
	//	"github.com/go-chassis/go-chassis/pkg/util/fileutil"
	//	"github.com/apache/servicecomb-mesher/cmd"
	//	"github.com/apache/servicecomb-mesher/config"
	"net/http"
	//	"os"
	//	"path/filepath"
	"github.com/go-chassis/go-chassis/core/lager"
	"github.com/stretchr/testify/assert"
)

// Testcase is trying to create files inside /tmp/build folder which is dynamic, so in travis it is not possible to create folder in prior, so can't test this case in travis
/*func TestInit(t *testing.T) {
	s, _ := fileutil.GetWorkDir()
	os.Setenv(fileutil.ChassisHome, s)
	chassisConf := filepath.Join(os.Getenv(fileutil.ChassisHome), "conf")
	os.MkdirAll(chassisConf, 0600)
	f, err := os.Create(filepath.Join(chassisConf, "chassis.yaml"))
	t.Log(f.Name())
	assert.NoError(t, err)
	f, err = os.Create(filepath.Join(chassisConf, "microservice.yaml"))
	t.Log(f.Name())
	assert.NoError(t, err)
	err = cConfig.Init()
	err = cmd.Init()
	lager.Initialize("", "INFO", "", "size", true, 1, 10, 7)
	archaius.Init()
	config.Init()
	err = Init()
	assert.NoError(t, err)
}*/

func TestResolve(t *testing.T) {
	lager.Initialize("", "DEBUG", "",
		"size", true, 1, 10, 7)
	d := &DefaultDestinationResolver{}
	header := http.Header{}
	header.Add("cookie", "user=jason")
	header.Add("X-Age", "18")
	destinationString := "Server"
	destinationString, p, err := d.Resolve("127.0.1.1", "abc", "", map[string]string{})
	assert.Error(t, err)
	assert.Equal(t, "", p)

	destinationString, p, err = d.Resolve("abc", "", "", map[string]string{})
	assert.Error(t, err)
	assert.Equal(t, "", p)

	destinationString, p, err = d.Resolve("abc", "", "http://127.0.0.1:80/test", map[string]string{})
	assert.NoError(t, err)
	assert.Equal(t, "80", p)

	destinationString, p, err = d.Resolve("abc", "", "Server:80/test", map[string]string{})
	assert.Error(t, err)

	destinationString, p, err = d.Resolve("abc", "", "127.0.0.1:80", map[string]string{})
	assert.Error(t, err)
	assert.Equal(t, "", p)
	t.Log(destinationString)
}

func TestGetDestinationResolver(t *testing.T) {
	dr := GetDestinationResolver("http")
	assert.NotNil(t, dr)
	dr = GetDestinationResolver("http_nil")
	assert.Nil(t, dr)
	SetDefaultDestinationResolver("http_nil", &DefaultDestinationResolver{})
	dr = GetDestinationResolver("http_nil")
	assert.NotNil(t, dr)
}

func BenchmarkDefaultDestinationResolver_Resolve(b *testing.B) {
	lager.Initialize("", "DEBUG", "",
		"size", true, 1, 10, 7)
	d := &DefaultDestinationResolver{}
	for i := 0; i < b.N; i++ {
		d.Resolve("abc", "", "http://127.0.0.1:80/test", map[string]string{})
	}
}

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

package authority_test

import (
	"testing"

	"github.com/apache/servicecomb-mesher/proxy/resolver"
	"github.com/apache/servicecomb-mesher/proxy/resolver/authority"
	"github.com/go-chassis/go-chassis/core/lager"
	"github.com/stretchr/testify/assert"
	"net/http"
)

func init() {
	lager.Init(&lager.Options{LoggerLevel: "DEBUG"})
}
func TestResolve(t *testing.T) {
	d := &authority.GRPCDefaultDestinationResolver{}
	header := http.Header{}
	header.Add("cookie", "user=jason")
	header.Add("X-Age", "18")
	var destinationString = "Server"
	destinationString, p, err := d.Resolve("abc", "", "127.0.1.1", map[string]string{})
	assert.Error(t, err)
	assert.Equal(t, "", p)

	destinationString, p, err = d.Resolve("abc", "", "", map[string]string{})
	assert.Error(t, err)
	assert.Equal(t, "", p)

	destinationString, p, err = d.Resolve("abc", "", "127.0.0.1:80", map[string]string{})
	assert.NoError(t, err)
	assert.Equal(t, "80", p)

	dr := resolver.GetDestinationResolver("grpc")

	destinationString, p, err = dr.Resolve("abc", "", "127.0.0.1:80", map[string]string{})
	assert.NoError(t, err)
	assert.Equal(t, "80", p)
	t.Log(destinationString)
}

func BenchmarkDefaultDestinationResolver_Resolve(b *testing.B) {
	d := &authority.GRPCDefaultDestinationResolver{}
	for i := 0; i < b.N; i++ {
		d.Resolve("abc", "", "127.0.0.1:80", map[string]string{})
	}
}

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

	"github.com/go-chassis/go-chassis/core/lager"
	"github.com/go-mesh/mesher/proxy/resolver"
	"github.com/go-mesh/mesher/proxy/resolver/authority"
	"github.com/stretchr/testify/assert"
	"net/http"
)

func TestResolve(t *testing.T) {
	lager.Initialize("", "DEBUG", "",
		"size", true, 1, 10, 7)
	d := &authority.GPRCDefaultDestinationResolver{}
	header := http.Header{}
	header.Add("cookie", "user=jason")
	header.Add("X-Age", "18")
	mystring := "Server"
	var destinationString = &mystring
	p, err := d.Resolve("abc", map[string]string{}, "127.0.1.1", destinationString)
	assert.Error(t, err)
	assert.Equal(t, "", p)

	p, err = d.Resolve("abc", map[string]string{}, "", destinationString)
	assert.Error(t, err)
	assert.Equal(t, "", p)

	p, err = d.Resolve("abc", map[string]string{}, "127.0.0.1:80", destinationString)
	assert.NoError(t, err)
	assert.Equal(t, "80", p)

	dr := resolver.GetDestinationResolver("grpc")

	p, err = dr.Resolve("abc", map[string]string{}, "127.0.0.1:80", destinationString)
	assert.NoError(t, err)
	assert.Equal(t, "80", p)

}

func BenchmarkDefaultDestinationResolver_Resolve(b *testing.B) {
	lager.Initialize("", "DEBUG", "",
		"size", true, 1, 10, 7)
	d := &authority.GPRCDefaultDestinationResolver{}
	mystring := "Server"
	var destinationString = &mystring
	for i := 0; i < b.N; i++ {
		d.Resolve("abc", map[string]string{}, "127.0.0.1:80", destinationString)
	}
}

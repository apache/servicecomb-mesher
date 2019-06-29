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

package util_test

import (
	"github.com/go-chassis/go-chassis/core/invocation"
	"github.com/go-chassis/go-chassis/pkg/util/tags"
	"github.com/go-mesh/mesher/proxy/config"
	"github.com/go-mesh/mesher/proxy/util"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEqualPolicy(t *testing.T) {

	i := &invocation.Invocation{
		MicroServiceName: "ShoppingCart",
	}

	i.RouteTags = utiltags.NewDefaultTag("0.1", "default")
	tags := make(map[string]string)
	tags["app"] = "default"
	tags["version"] = "0.0.1"
	p := &config.Policy{
		Destination:   "ShoppingCart1",
		Tags:          tags,
		LoadBalancing: nil,
	}
	value := util.EqualPolicy(i, p)
	assert.Equal(t, value, false)

	i = &invocation.Invocation{
		MicroServiceName: "ShoppingCart",
	}

	i.RouteTags = utiltags.NewDefaultTag("0.1", "default1")
	tags = make(map[string]string)
	tags["app"] = "default"
	tags["version"] = "0.0.1"
	p = &config.Policy{
		Destination:   "ShoppingCart",
		Tags:          tags,
		LoadBalancing: nil,
	}
	value = util.EqualPolicy(i, p)
	assert.Equal(t, value, false)
	i = &invocation.Invocation{
		MicroServiceName: "ShoppingCart",
	}

	i.RouteTags = utiltags.NewDefaultTag("0.1", "default")
	tags = make(map[string]string)
	tags["app"] = "default"

	tags["version"] = "0.0.1"
	p = &config.Policy{
		Destination:   "ShoppingCart",
		Tags:          tags,
		LoadBalancing: nil,
	}
	value = util.EqualPolicy(i, p)
	assert.Equal(t, value, false)

	i = &invocation.Invocation{
		MicroServiceName: "ShoppingCart",
	}

	i.RouteTags = utiltags.NewDefaultTag("0.1", "default")
	tags = make(map[string]string)
	tags["app1"] = "default"
	tags["version1"] = "0.1"
	p = &config.Policy{
		Destination:   "ShoppingCart",
		Tags:          tags,
		LoadBalancing: nil,
	}
	value = util.EqualPolicy(i, p)
	assert.Equal(t, value, false)

	i = &invocation.Invocation{
		MicroServiceName: "ShoppingCart",
	}
	inv := make(map[string]interface{})
	inv["app1"] = 1
	i.Metadata = inv
	i.RouteTags = utiltags.NewDefaultTag("0.1", "default")
	tags = make(map[string]string)
	tags["app1"] = "default"
	tags["version1"] = "0.1"
	p = &config.Policy{
		Destination:   "ShoppingCart",
		Tags:          tags,
		LoadBalancing: nil,
	}
	value = util.EqualPolicy(i, p)
	assert.Equal(t, value, false)

	i = &invocation.Invocation{
		MicroServiceName: "ShoppingCart",
	}
	inv = make(map[string]interface{})
	inv["app"] = "default"
	i.Metadata = inv
	i.RouteTags = utiltags.NewDefaultTag("0.1", "default")
	tags = make(map[string]string)
	tags["app"] = "default1"
	tags["version"] = "0.1"
	p = &config.Policy{
		Destination:   "ShoppingCart",
		Tags:          tags,
		LoadBalancing: nil,
	}
	value = util.EqualPolicy(i, p)
	assert.Equal(t, value, false)

	i = &invocation.Invocation{
		MicroServiceName: "ShoppingCart",
	}

	i.RouteTags = utiltags.NewDefaultTag("0.1", "default")
	tags = make(map[string]string)
	tags["app"] = "default"
	tags["version"] = "0.1"
	p = &config.Policy{
		Destination:   "ShoppingCart",
		Tags:          tags,
		LoadBalancing: nil,
	}
	value = util.EqualPolicy(i, p)
	assert.Equal(t, value, true)

}

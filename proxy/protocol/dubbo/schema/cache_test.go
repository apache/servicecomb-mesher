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

package schema

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFreshTicker(t *testing.T) {
	r1 := newRefresher(time.Second * 2)
	r1.Run()

	a1, a2 := 0, 0

	r1.Add(Job{Fn: func() {
		a1++
	}})
	time.Sleep(time.Second)
	r1.Add(Job{Fn: func() {
		a2++
	}})

	select {
	case <-time.After(time.Second * 5):
	}

	assert.Equal(t, a1, 2)
	assert.Equal(t, a2, 2)

}

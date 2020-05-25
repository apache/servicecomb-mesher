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

package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLong2bytes_Bytes2long(t *testing.T) {
	bytes := make([]byte, 8)
	v := int64(12345)
	Long2bytes(v, bytes, 0)

	v1 := Bytes2long(bytes, 0)
	assert.Equal(t, v, v1)
}

func TestShort2bytes_Bytes2short(t *testing.T) {
	bytes := make([]byte, 2)
	v := int(11)
	Short2bytes(v, bytes, 0)
	t.Log(bytes)

	v1 := Bytes2short(bytes, 0)
	assert.Equal(t, uint16(v), v1)
}

func TestInt2bytes_Bytes2int(t *testing.T) {
	bytes := make([]byte, 4)
	v := int(11)
	Int2bytes(v, bytes, 0)
	t.Log(bytes)

	v1 := Bytes2int(bytes, 0)
	assert.Equal(t, int32(v), v1)
}

func TestS2ByteSlice(t *testing.T) {
	bytes := []string{"Hello", " ", "World"}
	ss := S2ByteSlice(bytes)
	t.Log(ss)
	assert.Equal(t, 3, len(ss))
}

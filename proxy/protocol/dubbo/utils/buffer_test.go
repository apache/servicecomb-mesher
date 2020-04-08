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

func TestWriteBuffe(t *testing.T) {
	var buffer WriteBuffer
	// case size = 0
	buffer.Init(0)
	// case size 1024
	buffer.Init(DefaultBufferSize)

	// Write []byte
	n, err := buffer.Write([]byte("byteDate"))
	assert.NoError(t, err)
	assert.Equal(t, len("byteDate"), n)

	// Write byte
	err = buffer.WriteByte(byte(12))
	assert.NoError(t, err)

	// Write bytes
	n = buffer.WriteBytes([]byte("byteDate"))
	assert.Equal(t, len("byteDate"), n)

	// Write Object
	m := make(map[string]string)
	m["key_01"] = "value_01"
	err = buffer.WriteObject(m)
	assert.NoError(t, err)

	// Written Bytes
	buffer.Init(24)
	n = buffer.WrittenBytes()
	assert.Equal(t, 0, n)

	buffer.WriteBytes([]byte("byteDate"))
	n = buffer.WrittenBytes()
	assert.Equal(t, len("byteDate"), n)

	// Get Buf
	b := buffer.GetBuf()
	assert.Less(t, 0, len(b))

}

func TestReadBuffe(t *testing.T) {
	var buffer WriteBuffer
	buffer.Init(DefaultBufferSize)

	// Write byte
	err := buffer.WriteByte(byte(12))
	assert.NoError(t, err)

	var readBuffer ReadBuffer
	readBuffer.SetBuffer(buffer.GetBuf())
	b := readBuffer.ReadByte()
	assert.Equal(t, byte(12), b)

	// Write Object
	m := make(map[string]string)
	m["key_01"] = "value_01"
	err = buffer.WriteObject(m)
	assert.NoError(t, err)
	m, err = readBuffer.ReadMap()
	assert.NoError(t, err)
	assert.Equal(t, "value_01", m["key_01"])

	// Write bytes
	n := buffer.WriteBytes([]byte("byteDate"))
	assert.Equal(t, len("byteDate"), n)
	bs := readBuffer.ReadBytes(len("byteDate"))
	assert.Equal(t, "byteDate", string(bs))

	err = buffer.WriteObject("string01")
	assert.NoError(t, err)

	str := readBuffer.ReadString()
	assert.Equal(t, "string01", str)
}

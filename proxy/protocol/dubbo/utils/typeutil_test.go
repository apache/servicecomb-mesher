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
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestArrayToQueryString(t *testing.T) {
	key := "key_01"
	value := []interface{}{"value01", "value02"}
	t.Log(ArrayToQueryString(key, value))

}

func TestObjectToString(t *testing.T) {
	// case java string
	str, err := ObjectToString(JavaString, "string01")
	assert.NoError(t, err)
	assert.Equal(t, "string01", str)

	str, err = ObjectToString(JavaString, "")
	assert.NoError(t, err)
	assert.Equal(t, "", str)

	// case java byte
	str, err = ObjectToString(JavaByte, "9")
	assert.NoError(t, err)
	assert.Equal(t, "9", str)

	// case java short
	str, err = ObjectToString(JavaShort, "9")
	assert.NoError(t, err)
	assert.Equal(t, "9", str)

	// case java Integer
	str, err = ObjectToString(JavaInteger, "9")
	assert.NoError(t, err)
	assert.Equal(t, "9", str)

	// case java Long
	str, err = ObjectToString(JavaLong, "9")
	assert.NoError(t, err)
	assert.Equal(t, "9", str)

	// case java Float
	str, err = ObjectToString(JavaFloat, "9.01")
	assert.NoError(t, err)
	assert.Equal(t, "9.01", str)

	// case java Float
	str, err = ObjectToString(JavaDouble, "9.01")
	assert.NoError(t, err)
	assert.Equal(t, "9.01", str)

	// case java Boolean
	str, err = ObjectToString(JavaBoolean, "false")
	assert.NoError(t, err)
	assert.Equal(t, "false", str)

	// case java Array
	str, err = ObjectToString(JavaArray, "[1,2,3,4]")
	assert.NoError(t, err)

	str, err = ObjectToString(SchemaArray, "[1,2,3,4]")
	assert.NoError(t, err)

	// case java Array
	data := struct {
		SvcName string
		Version int
	}{"name", 100}

	str, err = ObjectToString(JavaObject, data)
	assert.NoError(t, err)

	str, err = ObjectToString(SchemaObject, data)
	assert.NoError(t, err)

	// Default type
	str, err = ObjectToString("Default type", "Hi")
	assert.NoError(t, err)

	// case value == nil
	str, err = ObjectToString(SchemaObject, nil)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(str))
}

func TestRestByteToValue(t *testing.T) {
	// case java string
	str, err := RestByteToValue(JavaString, []byte("string01"))
	assert.NoError(t, err)
	assert.Equal(t, "string01", str)

	// case java short
	bytes := make([]byte, 2)
	v := int(11)
	Short2bytes(v, bytes, 0)

	str, err = RestByteToValue(JavaShort, bytes)
	assert.NoError(t, err)
	assert.Equal(t, int16(v), str)

	// case java Integer
	bytes = make([]byte, 4)
	v = int(11)
	Int2bytes(v, bytes, 0)

	str, err = RestByteToValue(JavaInteger, bytes)
	assert.NoError(t, err)
	assert.Equal(t, int32(v), str)

}

func TestTypeDesToArgsObjArry(t *testing.T) {
	arg := TypeDesToArgsObjArry(JavaString)
	t.Log(arg)
}

func TestArgumen(t *testing.T) {
	arg := Argument{JavaType: JavaString, Value: "string01"}
	assert.Equal(t, JavaString, arg.GetJavaType())
	assert.Equal(t, "string01", arg.GetValue())

	arg.SetJavaType(JavaShort)
	arg.SetValue(99)
	assert.Equal(t, JavaShort, arg.GetJavaType())
	assert.Equal(t, 99, arg.GetValue())

}

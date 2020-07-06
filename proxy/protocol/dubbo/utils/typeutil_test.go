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
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestArrayToQueryString(t *testing.T) {
	key := "key_01"
	value := []interface{}{"value01", "value02"}
	s := ArrayToQueryString(key, value)
	t.Log(s)
	assert.Equal(t, `key_01=value01&key_01=value02`, s)

	// case not []interface{} type
	s = ArrayToQueryString(key, "")
	assert.Equal(t, "", s)
}

func TestObjectToString(t *testing.T) {
	// case java string
	str, err := ObjectToString(JavaString, "string01")
	assert.NoError(t, err)
	assert.Equal(t, "string01", str)

	str, err = ObjectToString(JavaString, "")
	assert.NoError(t, err)
	assert.Equal(t, "", str)

	// case Unsurported Type
	str, err = ObjectToString(JavaChar, "string01")
	assert.Error(t, err)

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

	// case Java Char byte
	bytes = make([]byte, 4)
	bytes[0] = byte(0)
	bytes[1] = byte(1)
	bytes[2] = byte(2)
	bytes[3] = byte(3)
	str, err = RestByteToValue(JavaChar, bytes)
	str, err = RestByteToValue(JavaByte, bytes)
	assert.NoError(t, err)
	assert.Equal(t, bytes, str)

	// case Java long
	bytes = make([]byte, 4)
	bytes[0] = byte(0)
	bytes[1] = byte(1)
	bytes[2] = byte(2)
	bytes[3] = byte(3)
	str, err = RestByteToValue(JavaLong, bytes)
	assert.NoError(t, err)

	// case Java float
	bytes = make([]byte, 4)
	bytes[0] = byte(0)
	bytes[1] = byte(1)
	bytes[2] = byte(2)
	bytes[3] = byte(3)
	str, err = RestByteToValue(JavaFloat, bytes)
	assert.NoError(t, err)

	// case Java double
	bytes = make([]byte, 8)
	bytes[0] = byte(0)
	bytes[1] = byte(1)
	bytes[2] = byte(2)
	bytes[3] = byte(3)
	str, err = RestByteToValue(JavaDouble, bytes)
	assert.NoError(t, err)

	// case Java boolean
	bytes = make([]byte, 8)
	bytes[0] = byte(0)
	bytes[1] = byte(1)
	bytes[2] = byte(2)
	bytes[3] = byte(3)
	str, err = RestByteToValue(JavaBoolean, bytes)
	assert.Error(t, err)

	// case Java avaArray, SchemaArray:
	str, err = RestByteToValue(JavaArray, bytes)
	assert.Error(t, err)
	str, err = RestByteToValue(SchemaArray, bytes)
	assert.Error(t, err)

	// case  JavaObject, SchemaObject:
	str, err = RestByteToValue(JavaObject, bytes)
	assert.Error(t, err)
	str, err = RestByteToValue(SchemaObject, bytes)
	assert.Error(t, err)

	type jsonObj struct {
		Name string `json:"name"`
	}

	bs, err := json.Marshal(jsonObj{"name"})
	assert.NoError(t, err)

	str, err = RestByteToValue(SchemaObject, bs)
	assert.NoError(t, err)
	m, ok := str.(map[string]interface{})
	assert.Equal(t, true, ok)
	assert.Equal(t, m["name"], "name")

	str, err = RestByteToValue(JavaObject, bs)
	assert.NoError(t, err)
	m, ok = str.(map[string]interface{})
	assert.Equal(t, true, ok)
	assert.Equal(t, m["name"], "name")

	str, err = RestByteToValue("not fount tag", bs)
	assert.Error(t, err)
}

func TestTypeDesToArgsObjArry(t *testing.T) {
	arg := TypeDesToArgsObjArry(JavaString)
	t.Log(arg)

	// case empty
	arg = TypeDesToArgsObjArry("")
	t.Log(arg)
	assert.Equal(t, 1, len(arg))
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

func TestGetJavaDesc(t *testing.T) {
	arg := make([]Argument, 1)
	arg = append(arg, Argument{JavaType: JavaString, Value: "string"})
	str := GetJavaDesc(arg)
	assert.Equal(t, JavaString, str)
}

func TestRestBytesToLstValue(t *testing.T) {
	arg := &Argument{}
	var err error
	bytesTmp := S2ByteSlice([]string{"v1"})

	// case type error
	arg.Value, err = RestBytesToLstValue("Not a type", bytesTmp)
	assert.Error(t, err)

	// case array empty
	arg.Value, err = RestBytesToLstValue(JavaString, make([][]byte, 0))
	assert.NoError(t, err)

	arg.Value, err = RestBytesToLstValue(JavaString, bytesTmp)
	assert.NoError(t, err)
}

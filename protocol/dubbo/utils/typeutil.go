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
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math"
	"net/url"
	"regexp"
)

const (
	JavaString  = "Ljava/lang/String;"
	JavaChar    = "Ljava/lang/Character;"
	JavaByte    = "Ljava/lang/Byte;"
	JavaShort   = "Ljava/lang/Short;"
	JavaInteger = "Ljava/lang/Integer;"
	JavaLong    = "Ljava/lang/Long;"
	JavaFloat   = "Ljava/lang/Float;"
	JavaDouble  = "Ljava/lang/Double;"
	JavaBoolean = "Ljava/lang/Boolean;"
	JavaArray   = "Ljava.util.Arrays;"
	JavaObject  = "Ljava/lang/Object;"
	JavaList    = "Ljava/util/List;"
	JavaMap     = "Ljava.util.Map;"
	JavaSplit   = ";"
)

//Constants ..
const (
	//vod
	JvmVoid = byte('V')
	//boolean(Z).
	JvmBool = byte('Z')
	// byte(B)
	JvmByte = byte('B')
	// char(C).
	JvmChar = byte('C')
	// double(D).
	JvmDouble = byte('D')
	//float(F).
	JvmFloat = byte('F')
	// int(I).
	JvmInt = byte('I')
	// long(J).
	JvmLong = byte('J')
	// short(S).
	JvmShort = byte('S')

	//正则表达式
	JavaIdentRegex = "(?:[_$a-zA-Z][_$a-zA-Z0-9]*)"
	ClassDesc      = "(?:L" + JavaIdentRegex + "(?:\\/" + JavaIdentRegex + ")*;)"
	ArrayDesc      = "(?:\\[+(?:(?:[VZBCDFIJS])|" + ClassDesc + "))"
	ArrayRegex     = "(?:(?:[VZBCDFIJS])|" + ClassDesc + "|" + ArrayDesc + ")"
)

const (
	SchemaString  = "string"
	SchemaArray   = "array" //对应java  ARRAY和List
	SchemaMap     = "object"
	SchemaObject  = "object"
	SchemaNumber  = "number"
	SchemaInteger = "integer"
	SchemaInt32   = "int32"
	SchemaInt64   = "int64"
	SchemaFloat   = "float"
	SchemaDouble  = "double"
	SchemaByte    = "byte"
	SchemaBin     = "binary"
	SchemaBool    = "boolean"
	SchemaDate    = "date"
	SchemaTime    = "date-time"
	SchemaPasswd  = "password"
)

//SchemeTypeMAP is a variable of type map
var SchemeTypeMAP map[string]string

func init() {
	SchemeTypeMAP = make(map[string]string)
	SchemeTypeMAP[SchemaString] = JavaString
	SchemeTypeMAP[SchemaArray] = JavaArray
	SchemeTypeMAP[SchemaObject] = JavaObject
	SchemeTypeMAP[SchemaNumber] = JavaFloat
	SchemeTypeMAP[SchemaInteger] = JavaInteger
	SchemeTypeMAP[SchemaInt32] = JavaInteger
	SchemeTypeMAP[SchemaInt64] = JavaLong
	SchemeTypeMAP[SchemaFloat] = JavaFloat
	SchemeTypeMAP[SchemaDouble] = JavaDouble
	SchemeTypeMAP[SchemaByte] = JavaByte
	SchemeTypeMAP[SchemaBin] = JavaByte
	SchemeTypeMAP[SchemaBool] = JavaBoolean
	SchemeTypeMAP[SchemaDate] = JavaString
	SchemeTypeMAP[SchemaTime] = JavaString
	SchemeTypeMAP[SchemaPasswd] = JavaString
}

//ArrayToQueryString is a function which converts array to a string
func ArrayToQueryString(key string, inlst interface{}) string {
	lst := inlst.([]interface{})
	var retstr = ""
	for i := 0; i < len(lst); i++ {
		tmp := lst[i]
		if i == 0 {
			retstr = fmt.Sprintf("%s=%s", key, url.QueryEscape(tmp.(string)))
		} else {
			tmpstr := fmt.Sprintf("&%s=%s", key, url.QueryEscape(tmp.(string)))
			retstr = retstr + tmpstr
		}
	}
	return retstr
}

//ObjectToString is a method which converts object to string
func ObjectToString(dtype string, v interface{}) (string, error) {
	if v == nil {
		return "", nil
	}
	switch dtype {
	case JavaString, SchemaString:
		return v.(string), nil
	case JavaChar:
	case JavaByte:
		return v.(string), nil
	case JavaShort:
		return v.(string), nil
	case JavaInteger:
		return v.(string), nil
	case JavaLong:
		return v.(string), nil
	case JavaFloat:
		return v.(string), nil
	case JavaDouble:
		return v.(string), nil
	case JavaBoolean:
		return v.(string), nil
	case JavaArray, SchemaArray:

		return "", nil
	case JavaObject, SchemaObject:
		tmp, _ := json.Marshal(v)
		return string(tmp), nil
	default: //默认无法识别的类型直接使用json格式化
		tmp, _ := json.Marshal(v)
		return string(tmp), nil
	}
	return "", &BaseError{"Unsurported Type"}
}

//RestBytesToLstValue is a function
func RestBytesToLstValue(jType string, value [][]byte) (interface{}, error) {
	var tmp []interface{}
	var err error
	if len(value) > 0 {
		tmp = make([]interface{}, len(value))
		for i := 0; i < len(value); i++ {
			tmp[i], err = RestByteToValue(jType, value[i])
			if err != nil {
				return nil, err
			}
		}
	} else {
		tmp = make([]interface{}, 0)
	}
	return tmp, nil
}

//RestByteToValue is a function which converts byte to value type
func RestByteToValue(jType string, value []byte) (interface{}, error) {
	switch jType {
	case JavaString, SchemaString:
		return string(value[:]), nil
	case JavaChar:
	case JavaByte:
		return value, nil
	case JavaShort:
		return int16(binary.BigEndian.Uint16(value)), nil
	case JavaInteger:
		return int32(binary.BigEndian.Uint32(value)), nil
	case JavaLong:
		return int64(binary.BigEndian.Uint32(value)), nil
	case JavaFloat:
		bits := binary.LittleEndian.Uint32(value)
		return math.Float32frombits(bits), nil
	case JavaDouble:
		bits := binary.LittleEndian.Uint64(value)
		return math.Float64frombits(bits), nil
	case JavaBoolean:
		return nil, &BaseError{"Not supported"}
	case JavaArray, SchemaArray:
		return nil, &BaseError{"Not supported "}
	case JavaObject, SchemaObject:
		var tmp interface{}
		err := json.Unmarshal(value, &tmp) //对象类型直接使用json格式
		if err != nil {
			return nil, err
		}
		return tmp, nil
	default:
		return nil, &BaseError{"Invalid type"}
	}
	return nil, nil
}

/* //BSIJFDZC
CHARACTER当做string处理
FLOAT DOUBLE当做DOUBLE
BYTE  SHORT  INTEGER 当做整形处理
*/

//GetJavaDesc is a function
func GetJavaDesc(args []Argument) string {
	tmpDesc := ""
	for _, tmp := range args {
		tmpDesc += tmp.GetJavaType()
	}
	return tmpDesc
}

//TypeDesToArgsObjArry is a function which converts description to array object
func TypeDesToArgsObjArry(desc string) []Argument {
	if len(desc) == 0 {
		return make([]Argument, 1, 1)
	}
	var tmpArgsAarry = make([]Argument, 64, 64)
	descBytes := []byte(desc)
	reg := regexp.MustCompile(ArrayRegex)

	var i = 0
	for _, match := range reg.FindAll(descBytes, -1) {
		tmpArgsAarry[i] = Argument{string(match[:]), nil}
		i++
	}

	if i == 64 {
		return tmpArgsAarry
	} else {
		var realArry = make([]Argument, i, i)
		copy(realArry, tmpArgsAarry)
		return realArry
	}
}

//Argument is a struct
type Argument struct {
	JavaType string
	Value    interface{}
}

//SetJavaType is method which sets javatype
func (p *Argument) SetJavaType(jType string) {
	p.JavaType = jType
}

//GetJavaType is a method which returns javatype
func (p *Argument) GetJavaType() string {
	return p.JavaType
}

//SetValue is a method which sets value
func (p *Argument) SetValue(value interface{}) {
	p.Value = value
}

//GetValue is a method which gets value
func (p *Argument) GetValue() interface{} {
	return p.Value
}

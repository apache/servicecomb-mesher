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
	"github.com/go-chassis/gohessian"
	"reflect"

	"fmt"
)

//BaseError is a struct
type BaseError struct {
	ErrMsg string
}

func (p *BaseError) Error() string {
	return p.ErrMsg
}

//TypMap is a variable of type map
var TypMap map[string]reflect.Type

func init() {
	TypMap = make(map[string]reflect.Type)
}

const (
	DefaultGrowSize   = 4096
	DefaultBufferSize = 1024
)

//ReadBuffer is a struct
type ReadBuffer struct {
	buffer   []byte
	rdInd    int
	length   int
	capacity int
}

//WriteBuffer is a struct
type WriteBuffer struct {
	buffer   []byte
	wrInd    int
	capacity int
}

//Init is a method to initialize write buffer attributes
func (b *WriteBuffer) Init(size int) {
	if size == 0 {
		size = DefaultBufferSize
	}
	b.buffer = make([]byte, size)
	b.wrInd = 0
	b.capacity = size
}

//Write is a method to write into buffer
func (b *WriteBuffer) Write(p []byte) (n int, err error) {
	result := b.WriteBytes(p)
	if result > 0 {
		return result, nil
	} else {
		return result, &BaseError{"Not enough space to write"}
	}
}

func (b *WriteBuffer) grow(n int) int {
	m := b.capacity
	if b.wrInd+n > m {
		var buf []byte
		buf = make([]byte, m+n)
		copy(buf, b.buffer[0:b.wrInd])
		b.buffer = buf
		b.capacity = m + n
	}

	return n
}

//WriteBytes is a method to write bytes
func (b *WriteBuffer) WriteBytes(src []byte) int {
	size := len(src)
	if b.capacity < (b.wrInd + size) {
		reSize := DefaultGrowSize
		if DefaultGrowSize < size {
			reSize = size
		}
		b.grow(reSize)
	}
	copy(b.buffer[b.wrInd:], src)
	b.wrInd = b.wrInd + size
	return size
}

//WriteIndex is a method to write index
func (b *WriteBuffer) WriteIndex(index int) error {
	if index <= b.capacity {
		b.wrInd = index
	} else {
		return &BaseError{fmt.Sprintf("Index(%d) over the capacity(%s)", index, b.capacity)}
	}
	return nil
}

//WriteByte is a method to write particular byte
func (b *WriteBuffer) WriteByte(src byte) error {
	gh := hessian.NewGoHessian(nil, nil)
	err := gh.ToBytes2(int32(src), b)
	return err
}

//WriteObject is a method to write object
func (b *WriteBuffer) WriteObject(src interface{}) error {
	gh := hessian.NewGoHessian(nil, nil)
	err := gh.ToBytes2(src, b)
	return err
}

//WrittenBytes is a methodto get amount of bytes written
func (b *WriteBuffer) WrittenBytes() int {
	return b.wrInd
}

//GetBuf is a method to get buffer
func (b *WriteBuffer) GetBuf() []byte {
	return b.buffer
}

//GetValidData is a method to get valid data
func (b *WriteBuffer) GetValidData() []byte {
	return b.buffer[0:b.wrInd]
}

//SetBuffer is a method to set buffer data
func (b *ReadBuffer) SetBuffer(src []byte) {
	b.buffer = src
	b.rdInd = 0
	b.capacity = len(src)
	b.length = len(src)
}

//Init is a method to initialize read buffer
func (b *ReadBuffer) Init(capacity int) {
	b.buffer = make([]byte, capacity)
	b.length = 0
	b.rdInd = 0
	b.capacity = capacity
}

//ReadByte is a method to read particular byte from buffer
func (b *ReadBuffer) ReadByte() byte {
	var tmp interface{}
	tmp, _ = b.ReadObject()
	return byte(tmp.(int32))
}

//ReadBytes is a method to read data from buffer
func (b *ReadBuffer) ReadBytes(len int) []byte {
	start := b.rdInd
	b.rdInd = b.rdInd + len
	return b.buffer[start:b.rdInd]
}

//ReadObject is a method to read buffer and return object
func (b *ReadBuffer) ReadObject() (interface{}, error) {
	gh := hessian.NewGoHessian(TypMap, nil)
	obj, err := gh.ToObject2(b)
	return obj, err
}

//ReadString is a method to read buffer and return as string
func (b *ReadBuffer) ReadString() string {
	gh := hessian.NewGoHessian(nil, nil)
	obj, _ := gh.ToObject2(b)
	return obj.(string)
}

//ReadMap is a method to read buffer and return as a map
func (b *ReadBuffer) ReadMap() (map[string]string, error) {
	gh := hessian.NewGoHessian(nil, nil)
	obj, err := gh.ToObject2(b)
	if err != nil {
		return nil, err
	} else {
		tmpMap := obj.(map[string]interface{})
		var strMap = make(map[string]string)
		for k, v := range tmpMap {
			strMap[k] = v.(string)
		}
		return strMap, nil
	}
}

//Read 实现io.Reader
func (b *ReadBuffer) Read(p []byte) (n int, err error) {
	size := len(p)
	if b.length > (b.rdInd + size) {
		copy(p, b.buffer[b.rdInd:b.rdInd+size])
		b.rdInd = b.rdInd + size
		return size, nil
	} else if b.length == b.rdInd {
		return 0, nil
	} else {
		cpysize := b.length - b.rdInd
		copy(p, b.buffer[b.rdInd:b.length])
		b.rdInd = b.length - 1
		return cpysize, nil
	}

}

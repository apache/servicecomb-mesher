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

//Bytes2long is a function which converts byte to int64
func Bytes2long(b []byte, off int) int64 {
	return int64(b[off+0])<<56 + int64(b[off+1])<<48 + int64(b[off+2]) + int64(b[off+3]) +
		int64(b[off+4])<<24 + int64(b[off+5])<<16 + int64(b[off+6])<<8 + int64(b[off+7])
}

//Long2bytes is a function which converts int64 to byte
func Long2bytes(i int64, b []byte, off int) {
	v := uint64(i)
	b[off+7] = byte(v)
	b[off+6] = byte(v >> 8)
	b[off+5] = byte(v >> 16)
	b[off+4] = byte(v >> 24)
	b[off+3] = byte(v >> 32)
	b[off+2] = byte(v >> 40)
	b[off+1] = byte(v >> 48)
	b[off+0] = byte(v >> 56)
}

//Short2bytes is a function which converts int to byte
func Short2bytes(value int, b []byte, off int) {
	v := value
	b[off+1] = byte(v)
	b[off+0] = byte(v >> 8)
}

//Bytes2short is a function which converts byte to int
func Bytes2short(b []byte, off int) uint16 {
	return (uint16)(((b[off+1] & 0xFF) << 0) +
		((b[off+0]) << 8))
}

//Int2bytes is a function which converts int to byte
func Int2bytes(value int, b []byte, off int) {
	v := uint32(value)
	b[off+3] = byte(v)
	b[off+2] = byte(v >> 8)
	b[off+1] = byte(v >> 16)
	b[off+0] = byte(v >> 24)
}

//Bytes2int is a function which converts byte to int
func Bytes2int(b []byte, off int) int32 {
	return int32(b[0+off])<<24 + int32(b[1+off])<<16 + int32(b[2+off])<<8 + int32(b[3+off])
}

//S2ByteSlice is a function which converts string slice to byte slice
func S2ByteSlice(str []string) [][]byte {
	b := make([][]byte, len(str))
	for i, s := range str {
		b[i] = []byte(s)
	}
	return b
}

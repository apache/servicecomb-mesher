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

package grpc

import (
	"github.com/gogo/protobuf/proto"
	"google.golang.org/grpc/encoding"
)

type frame struct {
	payload []byte
}

type codec struct {
	protoCodec encoding.Codec
}

//Marshal encode struct into bytes
func (c *codec) Marshal(f interface{}) ([]byte, error) {
	b, ok := f.(*frame)
	if !ok {
		return c.protoCodec.Marshal(f)
	}
	return b.payload, nil

}

//Unmarshal decode bytes into struct
func (c *codec) Unmarshal(d []byte, f interface{}) error {
	b, ok := f.(*frame)
	if !ok {
		return c.protoCodec.Unmarshal(d, f)
	}
	b.payload = d
	return nil
}

//String return name
func (c *codec) String() string {
	return "proxy_codec"
}

type protobufCodec struct {
}

//Marshal encode struct into bytes
func (c *protobufCodec) Marshal(f interface{}) ([]byte, error) {
	return proto.Marshal(f.(proto.Message))

}

//Unmarshal decode bytes into struct
func (c *protobufCodec) Unmarshal(d []byte, f interface{}) error {
	return proto.UnmarshalMerge(d, f.(proto.Message))
}

//Name return name
func (c *protobufCodec) Name() string {
	return "protobuf_codec"
}

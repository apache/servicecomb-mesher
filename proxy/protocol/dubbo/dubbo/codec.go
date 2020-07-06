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

package dubbo

import (
	"github.com/apache/servicecomb-mesher/proxy/protocol/dubbo/utils"
)

//Constants for request and response attributes
const (
	HeaderLength      = 16
	Magic             = 0xdabb
	MagicHigh         = byte(0xda)
	MagicLow          = byte(0xbb)
	FlagRequest       = byte(0x80)
	FlagTwoWay        = byte(0x40)
	FlagEvent         = byte(0x20)
	SerializationMask = byte(0x1f)
	HeartBeatEvent    = ""
)

//Constants for dubbo attributes
const (
	DubboVersionKey    string = "dubbo"
	DubboVersion       string = "2.0.0"
	PathKey            string = "path"
	InterfaceKey       string = "interface"
	VersionKey         string = "version"
	CommaSeparator     string = ","
	FileSeparator      string = "/"
	SemicolonSeparator string = ";"
)

//Constants
const (
	Success              = 0
	NeedMore             = -1
	InvalidFragement     = -2
	InvalidSerialization = -3
)

//serialise type
const (
	Hessian2 = byte(2)
)

//DubboCodec is a struct
type DubboCodec struct {
}

//GetContentTypeID is a method which returns content type id
func (p *DubboCodec) GetContentTypeID() byte {
	return Hessian2
}

//EncodeDubboRsp is a method which encodes dubbo response
func (p *DubboCodec) EncodeDubboRsp(rsp *DubboRsp, buffer *util.WriteBuffer) int {
	// set Magic number.
	header := make([]byte, HeaderLength)
	// set Magic number.
	util.Short2bytes(Magic, header, 0)
	// set request and serialization flag.
	header[2] = p.GetContentTypeID()
	if rsp.IsHeartbeat() {
		header[2] |= FlagEvent
	}
	// set response status.
	status := rsp.GetStatus()
	header[3] = status
	// set request id.
	util.Long2bytes(rsp.GetID(), header, 4)
	buffer.WriteIndex(HeaderLength)
	if status == Ok {
		if rsp.IsHeartbeat() {
			//encodeHeartbeatData
			ret := rsp.GetValue()
			buffer.WriteObject(ret)
		} else {
			//encodeResponseData
			except := rsp.GetException()
			if except == nil {
				ret := rsp.GetValue()
				if ret == nil {
					buffer.WriteByte(ResponseNullValue)
				} else {
					//buffer.WriteByte(ResponseValue)
					if v, ok := ret.([]byte); ok {
						buffer.WriteBytes(v)
					} else {
						buffer.WriteObject(ret)
					}

				}
			} else {
				buffer.WriteByte(ResponseWithException)
				buffer.WriteObject(except)
			}
		}
	} else {
		if rsp.GetErrorMsg() == "" {
			buffer.WriteByte(ResponseNullValue)
		} else {
			buffer.WriteObject(rsp.GetErrorMsg())
		}

	}

	len := buffer.WrittenBytes() - HeaderLength
	util.Int2bytes(len, header, 12)

	buffer.WriteIndex(0)
	buffer.WriteBytes(header)
	buffer.WriteIndex(HeaderLength + len)

	return 0
}

//DecodeDubboRsqHead is a method which decodes dubbo response header
func (p *DubboCodec) DecodeDubboRsqHead(rsp *DubboRsp, header []byte, bodyLen *int) int {
	if header[0] != MagicHigh || header[1] != MagicLow {
		return InvalidFragement
	}
	//读取请求ID
	var id int64 = util.Bytes2long(header, 4)
	rsp.SetID(id)

	var flag = header[2]
	if (flag & FlagEvent) != 0 {
		rsp.SetEvent(true)
	}
	proto := byte(flag & SerializationMask)

	if proto != Hessian2 { //当前只支持hessian2编码
		return InvalidSerialization
	}
	status := header[3]
	rsp.SetStatus(status)
	//读取长度
	*bodyLen = int(util.Bytes2int(header, 12))
	return Success
}

//DecodeDubboRspBody is a method which decodes dubbo response body
func (p *DubboCodec) DecodeDubboRspBody(buffer *util.ReadBuffer, rsp *DubboRsp) int {
	var obj interface{}
	var err error
	if rsp.IsHeartbeat() {
		rsp.SetValue(HeartBeatEvent)
	}

	if rsp.GetStatus() != Ok {
		obj, err = buffer.ReadObject()
		if err != nil {
			rsp.SetErrorMsg(err.Error())
			return 0
		}
		if s, ok := obj.(string); !ok {
			rsp.SetErrorMsg("unknown error")
		} else {
			rsp.SetErrorMsg(s)
		}
		return 0
	}

	switch rsp.IsHeartbeat() {
	case true:
		//decodeHeartbeatData
		obj, err = buffer.ReadObject()
		if err != nil {
			rsp.SetStatus(ServerError)
			rsp.SetErrorMsg(err.Error())
			return 0
		}
	case false:
		//decodeResult
		valueType, err := buffer.ReadByte()
		if err != nil {
			rsp.SetStatus(ServerError)
			rsp.SetErrorMsg(err.Error())
			return 0
		}
		switch valueType {
		case ResponseNullValue:
			//do nothing
			rsp.SetValue(nil)
			return 0
		case ResponseValue:
			obj, err = buffer.ReadObject()
			if err != nil {
				rsp.SetStatus(ServerError)
				rsp.SetErrorMsg(err.Error())
				return -1
			}
		case ResponseWithException:
			//readObject,设置异常
			rsp.SetStatus(ServiceError)
			obj, err = buffer.ReadObject()
			if err != nil {
				rsp.SetStatus(ServerError)
				rsp.SetErrorMsg(err.Error())
				return 0
			}
		}
		rsp.SetValue(buffer.GetBuf())
	}

	return 0
}

//EncodeDubboReq is a method which encodes dubbo request
func (p *DubboCodec) EncodeDubboReq(req *Request, buffer *util.WriteBuffer) int {
	// set Magic number.
	header := make([]byte, HeaderLength)
	util.Short2bytes(Magic, header, 0)
	// set request and serialization flag.
	header[2] = (byte)(FlagRequest | p.GetContentTypeID())
	if req.IsHeartbeat() {
		header[2] |= FlagEvent
	}
	if req.IsEvent() {
		header[2] |= FlagEvent
	}
	if req.IsTwoWay() {
		header[2] |= FlagTwoWay
	}

	status := req.GetStatus()
	header[3] = status
	// set request id.
	util.Long2bytes(req.GetMsgID(), header, 4)
	if buffer.WriteIndex(HeaderLength) != nil {
		return -1
	}

	if v, ok := req.GetData().([]byte); ok {
		buffer.WriteBytes(v)
	}
	buffer.WriteObject(req.GetAttachments())

	len := buffer.WrittenBytes() - HeaderLength
	util.Int2bytes(len, header, 12)
	buffer.WriteIndex(0)
	buffer.WriteBytes(header)
	buffer.WriteIndex(HeaderLength + len)

	return 0
}

//DecodeDubboReqBodyForRegstry is a method which decodes dubbo request body from registry
func (p *DubboCodec) DecodeDubboReqBodyForRegstry(req *Request, bodyBuf *util.ReadBuffer) int {
	var obj interface{}
	var err error
	if req.IsHeartbeat() {
		//decodeHeartbeatData
		obj, err = bodyBuf.ReadObject()
		if err != nil {
			req.SetData(err.Error())
			req.SetBroken(true)
			return -1
		}
	} else if req.IsEvent() {
		//decodeEventData
		obj, err = bodyBuf.ReadObject()
		if err != nil {
			req.SetData(err.Error())
			req.SetBroken(true)
			return -1
		}
	} else {
		req.SetAttachment(DubboVersionKey, bodyBuf.ReadString())
		req.SetAttachment(PathKey, bodyBuf.ReadString())
		req.SetAttachment(VersionKey, bodyBuf.ReadString())
		req.SetVersion(req.GetAttachment(VersionKey, ""))
		req.SetMethodName(bodyBuf.ReadString())

		//解析参数
		typeDesc := string(bodyBuf.ReadString())
		agrsArry := util.TypeDesToArgsObjArry(typeDesc)
		if typeDesc == "" {
			agrsArry = nil
		} else {
			size := len(agrsArry)
			if req.GetMethodName() == "subscribe" {
				size = 1
			}
			for i := 0; i < size; i++ {
				val, err := bodyBuf.ReadObject()
				if err != nil {
					req.SetBroken(true)
					req.SetData(err.Error())
					return -1
				} else {
					agrsArry[i].SetValue(val)
				}
			}
			req.SetArguments(agrsArry)
		}

		if err == nil {
			req.SetAttachments(nil)
		} else {
			req.SetBroken(true)
			req.SetData(err.Error())
			return -1
		}
		req.SetBroken(false)
		req.SetData(obj)
	}

	return 0
}

//DecodeDubboReqBody is a method which decodes dobbo request body
func (p *DubboCodec) DecodeDubboReqBody(req *Request, bodyBuf *util.ReadBuffer) int {
	var obj interface{}
	var err error
	if req.IsHeartbeat() {
		//decodeHeartbeatData
		obj, err = bodyBuf.ReadObject()
		if err != nil {
			req.SetData(err.Error())
			req.SetBroken(true)
			return -1
		}
	} else if req.IsEvent() {
		//decodeEventData
		obj, err = bodyBuf.ReadObject()
		if err != nil {
			req.SetData(err.Error())
			req.SetBroken(true)
			return -1
		}
	} else {
		req.SetAttachment(DubboVersionKey, bodyBuf.ReadString())
		req.SetAttachment(PathKey, bodyBuf.ReadString())
		req.SetAttachment(VersionKey, bodyBuf.ReadString())
		req.SetVersion(req.GetAttachment(VersionKey, ""))
		req.SetMethodName(bodyBuf.ReadString())

		obj = bodyBuf.GetBuf()
		req.SetBroken(false)
		req.SetData(obj)
	}

	return 0
}

//DecodeDubboReqHead is a method which decodes dubbo request header
func (p *DubboCodec) DecodeDubboReqHead(req *Request, header []byte, bodyLen *int) int {
	if len(header) < HeaderLength {
		return NeedMore
	}
	//读取Magic
	if header[0] != MagicHigh || header[1] != MagicLow {
		return InvalidFragement
	}
	//读取请求ID
	var id = util.Bytes2long(header, 4)

	var flag = header[2]
	proto := byte(flag & SerializationMask)

	if proto != Hessian2 { //当前只支持hessian2编码
		return InvalidSerialization
	}

	if (flag & FlagRequest) == 0 {
		return InvalidFragement
	}
	req.SetMsgID(id)
	req.SetVersion(DubboVersion)
	req.SetTwoWay((flag & FlagTwoWay) != 0)
	if (flag & FlagEvent) != 0 {
		req.SetEvent(HeartBeatEvent)
	}
	//读取长度
	*bodyLen = int(util.Bytes2int(header, 12))

	return Success
}

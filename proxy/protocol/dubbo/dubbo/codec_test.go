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
	"testing"

	"github.com/apache/servicecomb-mesher/proxy/protocol/dubbo/utils"
	"github.com/go-chassis/gohessian"

	"github.com/stretchr/testify/assert"
)

func TestDubboCodec_DecodeDubboReqBody(t *testing.T) {
	t.Log("If returns of rbf.ReadObject() is nil, should not panic")
	d := &DubboCodec{}

	req := NewDubboRequest()
	resp := &DubboRsp{}
	resp.Init()

	wbf := &util.WriteBuffer{}
	rbf := &util.ReadBuffer{}
	rbf.SetBuffer([]byte{hessian.BC_NULL})
	c := make([]byte, 10)
	_, err := rbf.Read(c)
	assert.NoError(t, err)
	assert.NotNil(t, c)
	assert.Equal(t, hessian.BC_NULL, c[0])
	obj, err := rbf.ReadObject()
	assert.Nil(t, err)
	assert.Nil(t, obj)

	assert.Equal(t, Hessian2, d.GetContentTypeID())

	// case EncodeDubboRsp status is ERROR
	t.Run("Test status error", func(t *testing.T) {
		resp.SetStatus(ServerError)
		d.EncodeDubboRsp(resp, wbf)
		d.DecodeDubboRspBody(rbf, resp)
	})

	// =====OK===============
	t.Run("Test status ok", func(t *testing.T) {
		// case event
		var buffer util.WriteBuffer
		buffer.Init(0)

		resp.SetStatus(Ok)
		resp.SetEvent(true)
		rbf.SetBuffer(append(buffer.GetBuf()[:buffer.WrittenBytes()], []byte{0x34, 0x02}...))
		d.EncodeDubboRsp(resp, wbf)
		d.DecodeDubboRspBody(rbf, resp)

		//case
		resp.SetStatus(Ok)
		resp.SetEvent(false)
		//resp.mEvent = true

		// case ResponseValue
		buffer.WriteIndex(0)
		buffer.WriteByte(ResponseValue)
		rbf.SetBuffer(append(buffer.GetBuf()[:buffer.WrittenBytes()], []byte{0x34, 0x02}...))

		d.EncodeDubboRsp(resp, wbf)
		d.DecodeDubboRspBody(rbf, resp)

		// case ResponseNullValue
		resp.SetStatus(Ok)
		resp.SetEvent(false)
		buffer.WriteIndex(0)
		buffer.WriteByte(ResponseNullValue)
		rbf.SetBuffer(append(buffer.GetBuf()[:buffer.WrittenBytes()], []byte{0x34, 0x02}...))

		d.EncodeDubboRsp(resp, wbf)
		d.DecodeDubboRspBody(rbf, resp)

		// case ResponseWithException
		resp.SetStatus(Ok)
		resp.SetEvent(false)
		buffer.WriteIndex(0)
		buffer.WriteByte(ResponseWithException)

		rbf.SetBuffer(append(buffer.GetBuf()[:buffer.WrittenBytes()], []byte{0x34, 0x02}...))

		d.EncodeDubboRsp(resp, wbf)
		d.DecodeDubboRspBody(rbf, resp)
	})

	// ResponseNullValue
	t.Run("Test ResponseNullValue", func(t *testing.T) {
		var buffer util.WriteBuffer
		buffer.Init(0)

		buffer.WriteByte(ResponseNullValue)
		rbf.SetBuffer(buffer.GetValidData())
		d.DecodeDubboRspBody(rbf, resp)

		// ResponseValue
		buffer.WriteByte(ResponseValue)
		buffer.WriteObject("Hello")
		rbf.SetBuffer(buffer.GetValidData())
		d.DecodeDubboRspBody(rbf, resp)

		// ResponseWithException
		buffer.WriteByte(ResponseWithException)
		rbf.SetBuffer(buffer.GetValidData())
		d.DecodeDubboRspBody(rbf, resp)
	})

	// case DecodeDubboReqBodyForRegstry
	t.Run("Test DecodeDubboReqBody", func(t *testing.T) {
		var buffer util.WriteBuffer
		buffer.Init(0)

		rbf := &util.ReadBuffer{}

		req.SetAttachment(DubboVersionKey, "dubbov1")
		req.SetAttachment(PathKey, "rest")
		req.SetAttachment(VersionKey, "1.0.0")
		req.SetVersion(req.GetAttachment(VersionKey, ""))
		req.SetMethodName("POST")

		buffer.WriteIndex(0)

		buffer.WriteObject("dubbov1")
		buffer.WriteObject("rest")
		buffer.WriteObject("1.0.0")
		buffer.WriteObject("1.0.0")

		buffer.WriteIndex(buffer.WrittenBytes())

		rbf.SetBuffer(buffer.GetValidData()[:])
		d.DecodeDubboReqBody(req, rbf)

		// case IsHeartbeat
		req.SetEvent("")
		buffer.WriteObject("Hello")
		d.DecodeDubboReqBody(req, rbf)

		rbf.SetBuffer([]byte{0x34, 0x02}) // tag not found
		d.DecodeDubboReqBody(req, rbf)

		//case IsEvent
		req.SetEvent("envent")
		buffer.WriteObject("Hello")
		d.DecodeDubboReqBody(req, rbf)

		rbf.SetBuffer([]byte{0x34, 0x02}) // tag not found
		d.DecodeDubboReqBody(req, rbf)
	})

	t.Run("Test DecodeDubboReqBodyForRegstry", func(t *testing.T) {
		var buffer util.WriteBuffer
		buffer.Init(0)
		req := NewDubboRequest()

		rbf := &util.ReadBuffer{}

		req.SetAttachment(DubboVersionKey, "dubbov1")
		req.SetAttachment(PathKey, "rest")
		req.SetAttachment(VersionKey, "1.0.0")
		req.SetVersion(req.GetAttachment(VersionKey, ""))
		req.SetMethodName("POST")

		buffer.WriteIndex(0)
		buffer.WriteObject("dubbov1")
		buffer.WriteObject("rest")
		buffer.WriteObject("1.0.0")
		buffer.WriteObject("1.0.0")

		//处理参数
		dubboArgs := make([]util.Argument, 2)
		for i := 0; i < 2; i++ {
			arg := &util.Argument{}
			bytesTmp := util.S2ByteSlice([]string{"v1"})
			arg.Value, err = util.RestBytesToLstValue("string", bytesTmp)
			arg.JavaType = "Ljava/lang/String;"
			dubboArgs[i] = *arg
		}
		req.SetArguments(dubboArgs)

		//buffer.WriteObject(req.GetAttachments())
		s := util.GetJavaDesc(dubboArgs)
		buffer.WriteObject(s)

		rbf.SetBuffer(buffer.GetValidData()[:])
		d.DecodeDubboReqBodyForRegstry(req, rbf)

		// case IsHeartbeat
		req.SetEvent("")
		buffer.WriteObject("Hello")
		d.DecodeDubboReqBodyForRegstry(req, rbf)

		rbf.SetBuffer([]byte{0x34, 0x02}) // tag not found
		d.DecodeDubboReqBodyForRegstry(req, rbf)

		//case IsEvent
		req.SetEvent("envent")
		buffer.WriteObject("Hello")
		d.DecodeDubboReqBodyForRegstry(req, rbf)

		rbf.SetBuffer([]byte{0x34, 0x02}) // tag not found
		d.DecodeDubboReqBodyForRegstry(req, rbf)

	})

	t.Run("Test DecodeDubboReqBodyForRegstry", func(t *testing.T) {
		var buffer util.WriteBuffer
		buffer.Init(0)
		d.EncodeDubboReq(req, &buffer)

		// case IsHeartbeat
		buffer.WriteIndex(0)
		req.SetEvent("")
		d.EncodeDubboReq(req, &buffer)

		//case IsEvent
		buffer.WriteIndex(0)
		req.SetEvent("envent")
		d.EncodeDubboReq(req, &buffer)

	})

	t.Run("Test DecodeDubboReqHead and DecodeDubboRsqHead", func(t *testing.T) {
		headBuf := make([]byte, HeaderLength)
		bodyLen := 0

		// case DecodeDubboReqHead
		d.DecodeDubboReqHead(req, headBuf, &bodyLen)
		bodyLen = 0
		// case DecodeDubboRsqHead
		d.DecodeDubboRsqHead(resp, headBuf, &bodyLen)

		// init for other test case
		bodyLen = 0
		util.Short2bytes(Magic, headBuf, 0)

		// case DecodeDubboReqHead
		d.DecodeDubboReqHead(req, headBuf, &bodyLen)
		bodyLen = 0
		// case DecodeDubboRsqHead
		d.DecodeDubboRsqHead(resp, headBuf, &bodyLen)

		// init for other test case
		bodyLen = 0
		util.Short2bytes(Magic, headBuf, 0)
		headBuf[2] = (byte)(FlagRequest | Hessian2)

		// case DecodeDubboReqHead
		d.DecodeDubboReqHead(req, headBuf, &bodyLen)
		bodyLen = 0
		// case DecodeDubboRsqHead
		d.DecodeDubboRsqHead(resp, headBuf, &bodyLen)
	})
	// case EncodeDubboRsp
	d.EncodeDubboRsp(resp, wbf)

	// case EncodeDubboReq
	d.EncodeDubboReq(req, wbf)

	// case DecodeDubboRspBody
	d.DecodeDubboRspBody(rbf, resp)

	// case DecodeDubboReqBodyForRegstry
	d.DecodeDubboReqBodyForRegstry(req, rbf)
	headBuf := make([]byte, HeaderLength)
	bodyLen := 0

	// case DecodeDubboReqHead
	d.DecodeDubboReqHead(req, headBuf, &bodyLen)
	bodyLen = 0

	// case DecodeDubboRsqHead
	d.DecodeDubboRsqHead(resp, headBuf, &bodyLen)

	GCurMSGID = 0
}

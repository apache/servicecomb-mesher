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

const (
	Ok                             = byte(20)
	ClientTimeout                  = byte(30)
	ServerTimeout                  = byte(31)
	BadRequest                     = byte(40)
	BadResponse                    = byte(50)
	ServiceNotFound                = byte(60)
	ServiceError                   = byte(70)
	ServerError                    = byte(80)
	ClentError                     = byte(90)
	ServerThreadPoolExhaustedError = byte(100)
)
const (
	ResponseWithException = byte(0)
	ResponseValue         = byte(1)
	ResponseNullValue     = byte(2)
)

//DubboRsp is a struct which has attributes for dubbo response
type DubboRsp struct {
	DubboRPCResult
	mID       int64
	mVersion  string
	mStatus   byte
	mEvent    bool
	mErrorMsg string
}

//Init method initializes value
func (p *DubboRsp) Init() {
	p.mID = 0
	p.mVersion = "0.0.0"
	p.mStatus = Ok
	p.mEvent = false
	p.mErrorMsg = ""
	//p.mResult = nil
}

//IsHeartbeat is a method which checks for heartbeat
func (p *DubboRsp) IsHeartbeat() bool {
	return p.mEvent
}

//SetEvent is a method which sets event
func (p *DubboRsp) SetEvent(bEvt bool) {
	p.mEvent = bEvt
}

//GetStatus is a method which gets status
func (p *DubboRsp) GetStatus() byte {
	return p.mStatus
}

//SetStatus is a method which sets status
func (p *DubboRsp) SetStatus(status byte) {
	p.mStatus = status
}

//GetID is a method which gets ID
func (p *DubboRsp) GetID() int64 {
	return p.mID
}

//SetID is a method which sets ID
func (p *DubboRsp) SetID(reqID int64) {
	p.mID = reqID
}

//GetErrorMsg is a method which gets error message
func (p *DubboRsp) GetErrorMsg() string {
	return p.mErrorMsg
}

//SetErrorMsg is a method which sets error message
func (p *DubboRsp) SetErrorMsg(err string) {
	p.mErrorMsg = err
}

//DubboRPCResult is a struct which has attibutes for dubbo rpc result
type DubboRPCResult struct {
	attchments map[string]string
	exception  interface{}
	value      interface{}
}

//NewDubboRPCResult is a function which create new dubbo rpc result
func NewDubboRPCResult() *DubboRPCResult {
	return &DubboRPCResult{make(map[string]string), nil, nil}
}

//GetValue is a method which gets value
func (p *DubboRPCResult) GetValue() interface{} {
	return p.value
}

//SetValue is a method which sets value
func (p *DubboRPCResult) SetValue(v interface{}) {
	p.value = v
}

//GetException is a method which gets exception
func (p *DubboRPCResult) GetException() interface{} {
	return p.exception
}

//SetException is a method which sets exception
func (p *DubboRPCResult) SetException(e interface{}) {
	p.exception = e
}

//GetAttachments is a method which gets attachment
func (p *DubboRPCResult) GetAttachments() map[string]string {
	return p.attchments
}

//SetAttachments is a method which sets attachment
func (p *DubboRPCResult) SetAttachments(attach map[string]string) {
	p.attchments = attach
}

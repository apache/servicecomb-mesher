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
	"sync"
)

//GCurMSGID is a variable of type int64
var GCurMSGID int64
var msgIDMtx = sync.Mutex{}

//GenerateMsgID is a function which generates message ID
func GenerateMsgID() int64 {
	msgIDMtx.Lock()
	defer msgIDMtx.Unlock()
	GCurMSGID++
	return GCurMSGID
}

//Request is a struct
type Request struct {
	DubboRPCInvocation
	msgID    int64
	status   byte
	event    bool
	twoWay   bool
	isBroken bool
	data     interface{}
}

//NewDubboRequest is a function which creates new dubbo request
func NewDubboRequest() *Request {
	tmp := &Request{}
	tmp.SetMsgID(GenerateMsgID())
	tmp.methodName = ""
	tmp.mVersion = DubboVersion
	tmp.status = Ok
	tmp.event = false
	tmp.twoWay = true
	tmp.isBroken = false
	tmp.arguments = nil
	tmp.attachments = make(map[string]string)
	tmp.urlPath = ""
	return tmp
}

//IsBroken check whether the connection is broken
func (p *Request) IsBroken() bool {
	return p.isBroken
}

//SetBroken sets connection as broken
func (p *Request) SetBroken(broken bool) {
	p.isBroken = broken

}

//SetEvent sets event to be true
func (p *Request) SetEvent(event string) {
	p.event = true
	p.data = event
}

//GetMsgID gets message ID
func (p *Request) GetMsgID() int64 {
	return p.msgID
}

//SetMsgID sets message ID
func (p *Request) SetMsgID(id int64) {
	p.msgID = id
}

//GetStatus gets the status
func (p *Request) GetStatus() byte {
	return p.status
}

//IsHeartbeat is method
func (p *Request) IsHeartbeat() bool {
	return p.event && HeartBeatEvent == p.data
}

//IsEvent checks whether event is present
func (p *Request) IsEvent() bool {
	return p.event
}

//SetTwoWay is a method which set the connection to two-way
func (p *Request) SetTwoWay(is bool) {
	p.twoWay = is
}

//IsTwoWay is a method which checks whether it is two-way connection
func (p *Request) IsTwoWay() bool {
	return p.twoWay
}

//SetData is a method which sets data
func (p *Request) SetData(data interface{}) {
	p.data = data
}

//GetData is a method which gets data
func (p *Request) GetData() interface{} {
	return p.data
}

//DubboRPCInvocation is a struct
type DubboRPCInvocation struct {
	methodName  string
	mVersion    string
	arguments   []util.Argument
	attachments map[string]string
	urlPath     string
}

//SetVersion is a method which sets version
func (p *DubboRPCInvocation) SetVersion(ver string) {
	p.mVersion = ver
}

//GetAttachment is a method which gets particular attachment
func (p *DubboRPCInvocation) GetAttachment(key string, defaultValue string) string {
	if _, ok := p.attachments[key]; ok {
		return p.attachments[key]
	} else {
		return defaultValue
	}
}

//GetAttachments which gets all attachments
func (p *DubboRPCInvocation) GetAttachments() map[string]string {
	return p.attachments
}

//GetMethodName is a method which will get method name
func (p *DubboRPCInvocation) GetMethodName() string {
	return p.methodName
}

//SetMethodName is a method which sets method name
func (p *DubboRPCInvocation) SetMethodName(name string) {
	p.methodName = name
}

//SetAttachment is a method which sets attachment
func (p *DubboRPCInvocation) SetAttachment(key string, value string) {
	if p.attachments == nil {
		p.attachments = make(map[string]string)
	}
	if value == "" { //is empty, remove the key
		delete(p.attachments, value)
	} else {
		p.attachments[key] = value
	}
}

//SetAttachments is a method which sets multiple attachment
func (p *DubboRPCInvocation) SetAttachments(attachs map[string]string) {
	p.attachments = attachs
}

//GetArguments is a method which gets arguments
func (p *DubboRPCInvocation) GetArguments() []util.Argument {
	return p.arguments
}

//SetArguments is a method which sets arguments
func (p *DubboRPCInvocation) SetArguments(agrs []util.Argument) {
	p.arguments = agrs
}

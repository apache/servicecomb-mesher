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

	"github.com/stretchr/testify/assert"
)

func TestGenerateMsgID(t *testing.T) {
	assert.Equal(t, int64(1), GenerateMsgID())
	assert.Equal(t, int64(2), GenerateMsgID())
	assert.Equal(t, int64(3), GenerateMsgID())
}

func Test_Request(t *testing.T) {
	req := NewDubboRequest()

	// broken
	req.SetBroken(true)
	assert.Equal(t, true, req.IsBroken())

	// data
	req.SetData("info")
	assert.Equal(t, "info", req.GetData())

	// msg id
	req.SetMsgID(int64(101))
	assert.Equal(t, int64(101), req.GetMsgID())

	// status
	assert.Equal(t, Ok, req.GetStatus())

	// event
	req.SetEvent("event happend")
	assert.Equal(t, true, req.IsEvent())

	// twoway
	req.SetTwoWay(false)
	assert.Equal(t, false, req.IsTwoWay())

	// version
	req.SetVersion("1.0.0")
	assert.Equal(t, "1.0.0", req.mVersion)

	// Attachments
	m := make(map[string]string)
	m["key_01"] = "value_01"
	m["key_02"] = "value_02"
	m["key_03"] = "value_03"

	req.SetAttachments(m)
	attch := req.GetAttachments()
	assert.NotNil(t, attch)
	assert.Equal(t, "value_03", attch["key_03"])
	assert.Equal(t, "value_03", req.GetAttachment("key_03", "defaultValue"))
	assert.Equal(t, "defaultValue", req.GetAttachment("key_04", "defaultValue"))

	// method name
	req.SetMethodName("methodname")
	assert.Equal(t, "methodname", req.GetMethodName())

}

package dubbo

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDubboDubboRsp(t *testing.T) {
	rsp := DubboRsp{}
	rsp.Init()
	assert.Equal(t, "0.0.0", rsp.mVersion)
	assert.Equal(t, false, rsp.mEvent)
	assert.Equal(t, "", rsp.mErrorMsg)
	assert.Equal(t, int64(0), rsp.mID)
	assert.Equal(t, Ok, rsp.mStatus)

	// Value
	rsp.SetValue("test")
	assert.NotNil(t, rsp.GetValue())

	// Attachments
	m := make(map[string]string)
	m["key_01"] = "value_01"
	m["key_02"] = "value_02"
	m["key_03"] = "value_03"
	rsp.SetAttachments(m)
	attch := rsp.GetAttachments()
	assert.NotNil(t, attch)
	assert.Equal(t, "value_03", attch["key_03"])

	// ID
	rsp.SetID(12345)
	assert.Equal(t, int64(12345), rsp.GetID())

	// Exception
	rsp.SetException("Java Throw Exception")
	assert.Equal(t, "Java Throw Exception", rsp.GetException())

	// Msg
	rsp.SetErrorMsg("Test error msg")
	assert.Equal(t, "Test error msg", rsp.GetErrorMsg())

	// IsHeartbeat
	assert.Equal(t, false, rsp.IsHeartbeat())
}

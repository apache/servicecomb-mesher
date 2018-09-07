package handler

import (
	"testing"

	"github.com/go-chassis/go-chassis/core/common"
	"github.com/go-chassis/go-chassis/core/config"
	"github.com/go-chassis/go-chassis/core/config/model"
	"github.com/go-chassis/go-chassis/core/handler"
	"github.com/go-chassis/go-chassis/core/invocation"
	"github.com/stretchr/testify/assert"
)

func TestPortRewriteHandler_ValidEndpoint(t *testing.T) {
	t.Log("testing port rewrite handler with valid endpoint")

	c := handler.Chain{}
	c.AddHandler(&PortSelectionHandler{})

	config.GlobalDefinition = &model.GlobalCfg{}
	config.GlobalDefinition.Cse.Handler.Chain.Consumer = make(map[string]string)
	config.GlobalDefinition.Cse.Handler.Chain.Consumer["outgoing"] = PortMapForPilot
	i := &invocation.Invocation{
		MicroServiceName: "service1",
		SchemaID:         "schema1",
		OperationID:      "SayHello",
		Endpoint:         "127.0.0.1:5555",
	}

	c.Next(i, func(r *invocation.Response) error {
		assert.NoError(t, r.Err)
		return r.Err
	})
	c.Reset()
}

func TestPortRewriteHandler_InValidEndpoint(t *testing.T) {
	t.Log("testing port rewrite handler with empty endpoint")

	c := handler.Chain{}
	c.AddHandler(&PortSelectionHandler{})

	config.GlobalDefinition = &model.GlobalCfg{}
	config.GlobalDefinition.Cse.Handler.Chain.Consumer = make(map[string]string)
	config.GlobalDefinition.Cse.Handler.Chain.Consumer["outgoing"] = PortMapForPilot
	i := &invocation.Invocation{
		MicroServiceName: "service1",
		SchemaID:         "schema1",
		OperationID:      "SayHello",
		Endpoint:         "",
	}

	c.Next(i, func(r *invocation.Response) error {
		assert.Error(t, r.Err)
		return r.Err
	})

	c.Reset()
}

func TestPortRewriteHandler_Names(t *testing.T) {
	handlerObject := &PortSelectionHandler{}
	name := handlerObject.Name()
	assert.Equal(t, PortMapForPilot, name)
}

func TestReplacePort_InvalidEndpoint(t *testing.T) {
	output, err := replacePort("grpc", "")
	assert.Error(t, err)
	assert.Equal(t, "", output)
}

func TestReplacePort_ValidEndpoint(t *testing.T) {
	output, err := replacePort(common.ProtocolRest, "127.0.0.1:80")
	assert.Equal(t, "127.0.0.1:30101", output)
	assert.NoError(t, err)
}

func BenchmarkReplacePort(b *testing.B) {
	for i := 0; i < b.N; i++ {
		replacePort(common.ProtocolRest, "127.0.0.1:80")
	}
}

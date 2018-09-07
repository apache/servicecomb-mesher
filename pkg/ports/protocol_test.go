package ports_test

import (
	"github.com/go-mesh/mesher/pkg/ports"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSetFixedPort(t *testing.T) {
	ports.SetFixedPort("rpc", "9090")
	assert.Equal(t, "9090", ports.GetFixedPort("rpc"))
}

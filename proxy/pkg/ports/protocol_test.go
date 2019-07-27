package ports_test

import (
	"github.com/apache/servicecomb-mesher/proxy/pkg/ports"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSetFixedPort(t *testing.T) {
	ports.SetFixedPort("rpc", "9090")
	assert.Equal(t, "9090", ports.GetFixedPort("rpc"))
}

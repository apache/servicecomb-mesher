package metrics

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewPromRecorder(t *testing.T) {
	_, err := NewPromRecorder(nil)
	assert.Error(t, err)

}

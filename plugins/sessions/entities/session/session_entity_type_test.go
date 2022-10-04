package session

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSessionType_Name(t *testing.T) {
	assert.Equal(t, "session", Type.Name())
}

func TestSessionType_New(t *testing.T) {
	assert.Equal(t, map[string]interface{}{
		"flash": map[string]interface{}{},
	}, Type.New(map[string]interface{}{}))
}

func TestSessionType_Validate(t *testing.T) {
	assert.Nil(t, Type.Validate(map[string]interface{}{}))
}
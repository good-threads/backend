package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPing(t *testing.T) {
	assert.Equal(t, "pong\n", Setup().Ping())
}

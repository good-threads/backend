package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPing(t *testing.T) {
	commonLogic := Setup()
	result := commonLogic.Ping()
	assert.Equal(t, "pong\n", result)
}

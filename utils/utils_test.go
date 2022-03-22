package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetUnixMillisTimestamp(t *testing.T) {
	ts := GetUnixMillisTimestamp()
	assert.True(t, ts > 0)
}

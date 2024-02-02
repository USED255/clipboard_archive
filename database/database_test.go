package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnectDatabase(t *testing.T) {
	connectDatabase("file::memory:?cache=shared")

	assert.NotNil(t, Orm)

	Close()
}

func TestCloseDatabase(t *testing.T) {
	connectDatabase("file::memory:?cache=shared")
	Close()
}

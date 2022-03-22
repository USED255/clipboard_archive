package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCloseDatabase(t *testing.T) {
	ConnectDatabase("file::memory:?cache=shared")
	CloseDatabase()
}

func TestConnectDatabase(t *testing.T) {
	ConnectDatabase("file::memory:?cache=shared")

	assert.NotNil(t, Orm)

	CloseDatabase()
}

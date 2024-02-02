package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateVersion0Database(t *testing.T) {
	var config Config
	connectDatabase("file::memory:?cache=shared")

	createVersion0Database()

	Orm.First(&config, "key = ?", "version")
	assert.Equal(t, "", config.Value)

	Close()
}

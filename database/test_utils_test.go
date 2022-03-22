package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateVersion0Database(t *testing.T) {
	var config Config
	ConnectDatabase("file::memory:?cache=shared")

	createVersion0Database()

	Orm.First(&config, "key = ?", "version")
	assert.Equal(t, "", config.Value)

	CloseDatabase()
}
func TestCreateVersion1Database(t *testing.T) {
	var config Config
	ConnectDatabase("file::memory:?cache=shared")

	createVersion1Database()

	Orm.First(&config, "key = ?", "version")
	assert.Equal(t, "1.0.0", config.Value)

	CloseDatabase()
}

func TestCreateVersion2Database(t *testing.T) {
	var config Config
	ConnectDatabase("file::memory:?cache=shared")

	createVersion2Database()

	Orm.First(&config, "key = ?", "version")
	assert.Equal(t, "2.0.0", config.Value)

	CloseDatabase()
}

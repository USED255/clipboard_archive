package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMigrateVersion(t *testing.T) {
	var config Config
	ConnectDatabase("file::memory:?cache=shared")

	MigrateVersion()

	Orm.First(&config, "key = ?", "version")
	assert.Equal(t, version, config.Value)

	CloseDatabase()
}

func TestMigrateVersion0Database(t *testing.T) {
	var config Config
	ConnectDatabase("file::memory:?cache=shared")
	createVersion0Database()

	MigrateVersion()

	Orm.First(&config, "key = ?", "version")
	assert.Equal(t, version, config.Value)

	CloseDatabase()
}

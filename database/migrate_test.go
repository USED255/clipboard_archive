package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMigrateVersion(t *testing.T) {
	var config Config
	connectDatabase("file::memory:?cache=shared")

	migrateVersion()

	Orm.First(&config, "key = ?", "version")
	assert.Equal(t, version, config.Value)

	Close()
}

func TestMigrateVersion0Database(t *testing.T) {
	var config Config
	connectDatabase("file::memory:?cache=shared")
	createVersion0Database()

	migrateVersion()

	Orm.First(&config, "key = ?", "version")
	assert.Equal(t, version, config.Value)

	Close()
}

package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMigrateVersion(t *testing.T) {
	var config Config
	ConnectDatabase("file::memory:?cache=shared")
	MigrateVersion()

	MigrateVersion()

	Orm.First(&config, "key = ?", "version")
	assert.Equal(t, version, config.Value)

	CloseDatabase()
}

func TestMigrateVersionInitializingDatabase(t *testing.T) {
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

func TestMigrateVersion1Database(t *testing.T) {
	var config Config
	ConnectDatabase("file::memory:?cache=shared")
	createVersion1Database()

	MigrateVersion()

	Orm.First(&config, "key = ?", "version")
	assert.Equal(t, version, config.Value)

	CloseDatabase()
}

func TestMigrateVersion2Database(t *testing.T) {
	var config Config
	ConnectDatabase("file::memory:?cache=shared")
	createVersion2Database()

	MigrateVersion()

	Orm.First(&config, "key = ?", "version")
	assert.Equal(t, version, config.Value)

	CloseDatabase()
}

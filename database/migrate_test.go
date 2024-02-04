package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMigrateVersion(t *testing.T) {
	connectDatabase("file::memory:?cache=shared")
	defer Close()

	migrateVersion()

	v, _ := getDatabaseVersion()
	assert.Equal(t, version, v)
}

func TestMigrateVersion0Database(t *testing.T) {
	connectDatabase("file::memory:?cache=shared")
	defer Close()

	createVersion0Database()

	migrateVersion()

	v, _ := getDatabaseVersion()
	assert.Equal(t, version, v)
}

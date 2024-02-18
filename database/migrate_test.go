package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMigrateVersion(t *testing.T) {
	OpenNoDatabase()
	defer Close()

	migrateVersion()

	v, _ := getDatabaseVersion()
	assert.Equal(t, version, v)
}

func TestMigrateVersion0Database(t *testing.T) {
	OpenNoDatabase()
	defer Close()

	createVersion0Database()

	migrateVersion()

	v, _ := getDatabaseVersion()
	assert.Equal(t, version, v)
}

func TestMigrateVersion4Database(t *testing.T) {
	OpenMemoryDatabase()
	defer Close()

	Orm.Save(&Config{Key: "version", Value: "4"})

	migrateVersion()

	v, _ := getDatabaseVersion()
	assert.Equal(t, version, v)
}

func TestMigrateInvalidVersion(t *testing.T) {
	OpenNoDatabase()
	defer Close()

	Orm.Exec(createConfigsTableQuery)
	Orm.Save(&Config{Key: "version", Value: "999"})
	err = migrateVersion()

	assert.Error(t, err)
}

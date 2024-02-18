package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnectDatabase(t *testing.T) {
	OpenNoDatabase()
	defer Close()

	assert.NotNil(t, Orm)
}

func TestGetDatabaseVersion(t *testing.T) {
	OpenNoDatabase()
	defer Close()

	InitializeDatabase()

	v, _ := getDatabaseVersion()
	assert.Equal(t, version, v)
}

func TestGetDatabaseVersion0(t *testing.T) {
	OpenNoDatabase()
	defer Close()

	v, _ := getDatabaseVersion()
	assert.Equal(t, int64(0), v)
}

func TestGetDatabaseVersion1(t *testing.T) {
	OpenNoDatabase()
	defer Close()

	Orm.AutoMigrate(&Config{})
	migrateVersion0To1()

	v, _ := getDatabaseVersion()
	assert.Equal(t, int64(1), v)
}

func TestGetDatabaseVersionError(t *testing.T) {
	OpenNoDatabase()
	defer Close()

	Orm.AutoMigrate(&Config{})
	Orm.Create(&Config{Key: "version", Value: "a"})

	v, err := getDatabaseVersion()
	assert.Error(t, err)
	assert.Equal(t, int64(0), v)
}

func TestGetMajorVersion(t *testing.T) {
	v, err := getMajorVersion("1.2.3")
	assert.NoError(t, err)
	assert.Equal(t, int64(1), v)

	v, err = getMajorVersion("0.0.0")
	assert.NoError(t, err)
	assert.Equal(t, int64(0), v)

	v, err = getMajorVersion("65535.0.0")
	assert.NoError(t, err)
	assert.Equal(t, int64(65535), v)
}

func TestGetMajorVersionError(t *testing.T) {
	v, err := getMajorVersion("a")
	assert.Error(t, err)
	assert.Equal(t, int64(0), v)

	v, err = getMajorVersion("1.1.1.1")
	assert.Error(t, err)
	assert.Equal(t, int64(0), v)

	v, err = getMajorVersion("-1.0.0")
	assert.Error(t, err)
	assert.Equal(t, int64(0), v)

	v, err = getMajorVersion("184467440737095516150.0.0")
	assert.Error(t, err)
	assert.Equal(t, int64(0), v)
}

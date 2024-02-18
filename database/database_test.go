package database

import (
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func TestOpenDatabase(t *testing.T) {
	Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
	defer Close()

	assert.NotNil(t, Orm)
}

func TestOpenDatabaseError(t *testing.T) {
	OpenMemoryDatabase()
	err = Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
	assert.Error(t, err)
}

func TestCloseDatabase(t *testing.T) {
	OpenMemoryDatabase()
	Close()
}

func TestCloseDatabaseError(t *testing.T) {
	err = Close()
	assert.Error(t, err)
}

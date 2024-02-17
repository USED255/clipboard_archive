package database

import (
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func TestOpenDatabase(t *testing.T) {
	OrmConfig = &gorm.Config{Logger: logger.Default.LogMode(logger.Info)}
	Open("file::memory:?cache=shared")
	defer Close()

	assert.NotNil(t, Orm)
}

func TestOpenDatabaseError(t *testing.T) {
	Open("")
	err = Open("")
	assert.Error(t, err)
}

func TestOpenDatabase2(t *testing.T) {
	Open2(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
	defer Close()

	assert.NotNil(t, Orm)
}

func TestOpenDatabaseError2(t *testing.T) {
	Open2(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
	err = Open2(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
	assert.Error(t, err)
}

func TestCloseDatabase(t *testing.T) {
	connectDatabase("file::memory:?cache=shared")
	Close()
}

func TestCloseDatabaseError(t *testing.T) {
	err = Close()
	assert.Error(t, err)
}

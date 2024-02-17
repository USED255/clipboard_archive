package database

import (
	"testing"

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

func TestCloseDatabase(t *testing.T) {
	connectDatabase("file::memory:?cache=shared")
	Close()
}

func TestCloseDatabaseError(t *testing.T) {
	err = Close()
	assert.Error(t, err)
}

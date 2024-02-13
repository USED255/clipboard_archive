package database

import (
	"errors"

	"gorm.io/gorm"
)

const Version = version

var OrmConfig *gorm.Config
var Orm *gorm.DB

func Open(dsn string) error {
	err = connectDatabase(dsn)
	if err != nil {
		return err
	}
	err = migrateVersion()
	if err != nil {
		return err
	}
	return nil
}

func Close() error {
	if Orm == nil {
		return errors.New("database not connected")
	}

	sqlDB, err := Orm.DB()
	if err != nil {
		return err
	}

	err = sqlDB.Close()
	if err != nil {
		return err
	}

	Orm = nil
	return nil
}

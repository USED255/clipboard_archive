package database

import (
	"gorm.io/gorm"
)

const Version = version

var OrmConfig *gorm.Config
var Orm *gorm.DB

func Open(dns string) error {
	err = connectDatabase(dns)
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

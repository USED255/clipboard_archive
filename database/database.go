package database

import (
	"errors"

	"gorm.io/gorm"
)

const Version = version

var OrmConfig *gorm.Config
var Orm *gorm.DB

func Open(dialector gorm.Dialector, config *gorm.Config) error {
	err = connectDatabase(dialector, config)
	if err != nil {
		return err
	}
	err = migrateVersion()
	if err != nil {
		//迁移失败
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
		//?
		return err
	}

	err = sqlDB.Close()
	if err != nil {
		//无法关闭?
		return err
	}

	Orm = nil
	return nil
}

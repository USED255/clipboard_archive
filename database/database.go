package database

import (
	"log"

	"gorm.io/gorm"
)

const Version = version

var Orm *gorm.DB
var err error

func Open(dns string) {
	connectDatabase(dns)
	migrateVersion()
}

func Close() {
	sqlDB, err := Orm.DB()
	if err != nil {
		log.Fatalf("Failed to get database handle: %s", err)
	}
	sqlDB.Close()
	Orm = nil
}

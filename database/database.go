package database

import (
	"log"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

const Version = version

var Orm *gorm.DB
var err error

func ConnectDatabase(dns string) {
	if Orm != nil {
		log.Fatalf("Database already connected")
	}
	Orm, err = gorm.Open(sqlite.Open(dns), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
}
func CloseDatabase() {
	sqlDB, err := Orm.DB()
	if err != nil {
		log.Fatalf("Failed to get database handle: %s", err)
	}
	sqlDB.Close()
	Orm = nil
}

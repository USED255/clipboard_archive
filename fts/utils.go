package fts

import (
	"errors"
	"log"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var orm *gorm.DB
var err error

func connect() error {
	dns := "fts.sqlite3"
	if orm != nil {
		return errors.New("database already connected")
	}
	orm, err = gorm.Open(sqlite.Open(dns), &gorm.Config{})
	if err != nil {
		return err
	}
	return nil
}
func init() {
	tx := orm.Begin()

	err = tx.Exec(createFts5TableQuery).Error
	if err != nil {
		tx.Rollback()
		log.Fatal(err)
	}
	tx.Commit()
}

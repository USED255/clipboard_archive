package database

import (
	"log"
)

const version = "3.0.0"

func MigrateVersion() {
	var _count int64
	var config Config
	var configMajorVersion uint64

	currentMajorVersion, err := getMajorVersion(version)
	if err != nil {
		log.Fatal(err)
	}

	err = Orm.AutoMigrate(&ClipboardItem{}, &Config{})
	if err != nil {
		log.Fatal(err)
	}

	Orm.Model(&ClipboardItem{}).Count(&_count)
	count := _count
	Orm.Model(&Config{}).Count(&_count)
	count = count + _count
	if count == 0 {
		initializingDatabase()
	}

migrate:
	Orm.First(&config, "key = ?", "version")
	if config.Key == "" {
		configMajorVersion = 0
	} else {
		configMajorVersion, err = getMajorVersion(config.Value)
		if err != nil {
			log.Fatal(err)
		}
	}
	log.Println("Current version: ", config.Value)

	switch configMajorVersion {
	case currentMajorVersion:
		return
	case 2:
		MigrateVersion2To3()
		goto migrate
	case 1:
		MigrateVersion1To2()
		goto migrate
	case 0:
		MigrateVersion0To1()
		goto migrate
	default:
		log.Fatal("Unsupported version: ", config.Value)
	}

}

func initializingDatabase() {
	log.Println("No data in database, initializing")

	tx := Orm.Begin()
	err = tx.Exec(CreateFts5TableQuery).Error
	if err != nil {
		tx.Rollback()
		log.Fatal(err)
	}
	err = tx.Create(&Config{Key: "version", Value: version}).Error
	if err != nil {
		tx.Rollback()
		log.Fatal(err)
	}
	tx.Commit()
}

func MigrateVersion2To3() {
	log.Println("Migrating to 3.0.0")

	tx := Orm.Begin()
	err := tx.Migrator().DropColumn(&ClipboardItem{}, "index")
	if err != nil {
		tx.Rollback()
		log.Fatal("Migration failed: ", err)
	}
	err = tx.Migrator().RenameColumn(&ClipboardItem{}, "id", "index")
	if err != nil {
		tx.Rollback()
		log.Fatal("Migration failed: ", err)
	}
	err = tx.Save(&Config{Key: "version", Value: "3.0.0"}).Error
	if err != nil {
		tx.Rollback()
		log.Fatal("Migration failed: ", err)
	}
	tx.Commit()
}

func MigrateVersion1To2() {
	log.Println("Migrating to 2.0.0")

	Query := `
	INSERT INTO clipboard_items_fts (
		rowid, 
		clipboard_item_text
	)
	SELECT clipboard_items.clipboard_item_time, clipboard_items.clipboard_item_text 
	FROM clipboard_items;
	`
	tx := Orm.Begin()
	err := tx.Exec(CreateFts5TableQuery).Error
	if err != nil {
		tx.Rollback()
		log.Fatal("Migration failed: ", err)
	}
	err = tx.Exec(Query).Error
	if err != nil {
		tx.Rollback()
		log.Fatal("Migration failed: ", err)
	}
	err = tx.Save(&Config{Key: "version", Value: "2.0.0"}).Error
	if err != nil {
		tx.Rollback()
		log.Fatal("Migration failed: ", err)
	}
	tx.Commit()
}

func MigrateVersion0To1() {
	log.Println("Migrating to 1.0.0")

	tx := Orm.Begin()
	err := tx.Create(&Config{Key: "version", Value: "1.0.0"}).Error
	if err != nil {
		tx.Rollback()
		log.Fatal("Migration failed: ", err)
	}
	tx.Commit()
}

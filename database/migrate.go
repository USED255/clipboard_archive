package database

import (
	"log"
)

const version = "3.0.0"

func getDatabaseVersion() uint64 {
	var config Config
	var configMajorVersion uint64

	Orm.First(&config, "key = ?", "version")
	if config.Key == "" {
		return 0
	}

	configMajorVersion, err = getMajorVersion(config.Value)
	if err != nil {
		log.Fatal(err)
	}
	return configMajorVersion
}

func migrateVersion() {
	var config Config
	var databaseVersion uint64

	if !Orm.Migrator().HasTable(&ClipboardItem{}) {
		initializingDatabase()
	}
	currentMajorVersion, err := getMajorVersion(version)
	if err != nil {
		log.Fatal(err)
	}

	for {
		databaseVersion = getDatabaseVersion()

		switch databaseVersion {
		case currentMajorVersion:
			return
		case 2:
			migrateVersion2To3()
			continue
		case 1:
			migrateVersion1To2()
			continue
		case 0:
			migrateVersion0To1()
			continue
		default:
			log.Fatal("Unsupported version: ", config.Value)
		}
	}
}

func initializingDatabase() {
	log.Println("No data in database, initializing")

	tx := Orm.Begin()

	err = tx.AutoMigrate(&ClipboardItem{}, &Config{})
	if err != nil {
		log.Fatal(err)
	}
	err = tx.Exec(createFts5TableQuery).Error
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

func migrateVersion2To3() {
	log.Println("Migrating to version 3")
	tx := Orm.Begin()
	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
			log.Fatal("Migration failed: ", err)
		}
	}()

	err = tx.Migrator().DropColumn(&ClipboardItem{}, "index")
	if err != nil {
		panic(err)
	}
	err = tx.Migrator().RenameColumn(&ClipboardItem{}, "id", "index")
	if err != nil {
		panic(err)
	}
	err = tx.Save(&Config{Key: "version", Value: "3.0.0"}).Error
	if err != nil {
		panic(err)
	}

	tx.Commit()
}

func migrateVersion1To2() {
	tx := Orm.Begin()
	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
			log.Fatal("Migration failed: ", err)
		}
	}()

	err := tx.Exec(createFts5TableQuery).Error
	if err != nil {
		panic(err)
	}
	err = tx.Exec(insertFts5TableQuery).Error
	if err != nil {
		panic(err)
	}
	err = tx.Save(&Config{Key: "version", Value: "2.0.0"}).Error
	if err != nil {
		panic(err)
	}
	tx.Commit()
}

func migrateVersion0To1() {
	tx := Orm.Begin()
	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
			log.Fatal("Migration failed: ", err)
		}
	}()

	err := tx.Create(&Config{Key: "version", Value: "1.0.0"}).Error
	if err != nil {
		panic(err)
	}
	tx.Commit()
}

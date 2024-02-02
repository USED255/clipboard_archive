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

func databaseIsNew() bool {
	var _count int64
	Orm.Model(&ClipboardItem{}).Count(&_count)
	count := _count
	Orm.Model(&Config{}).Count(&_count)
	count = count + _count
	return count == 0
}

func MigrateVersion() {
	var config Config
	var databaseVersion uint64

	currentMajorVersion, err := getMajorVersion(version)
	if err != nil {
		log.Fatal(err)
	}

	err = Orm.AutoMigrate(&ClipboardItem{}, &Config{})
	if err != nil {
		log.Fatal(err)
	}

	if databaseIsNew() {
		initializingDatabase()
	}

	for {
		databaseVersion = getDatabaseVersion()
		log.Println("Current version: ", config.Value)

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

func migrateVersion2To3() {
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

func migrateVersion1To2() {
	log.Println("Migrating to 2.0.0")

	tx := Orm.Begin()
	err := tx.Exec(CreateFts5TableQuery).Error
	if err != nil {
		tx.Rollback()
		log.Fatal("Migration failed: ", err)
	}
	err = tx.Exec(InsertFts5TableQuery).Error
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

func migrateVersion0To1() {
	log.Println("Migrating to 1.0.0")

	tx := Orm.Begin()
	err := tx.Create(&Config{Key: "version", Value: "1.0.0"}).Error
	if err != nil {
		tx.Rollback()
		log.Fatal("Migration failed: ", err)
	}
	tx.Commit()
}

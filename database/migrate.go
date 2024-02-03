package database

import (
	"errors"
	"log"
	"strconv"
)

func getDatabaseVersion() (uint64, error) {
	var config Config
	var databaseVersion uint64

	Orm.First(&config, "key = ?", "version")
	if config.Key == "" {
		return 0, nil
	}
	databaseVersion, _ = strconv.ParseUint(config.Value, 10, 64)
	if databaseVersion != 0 {
		return databaseVersion, nil
	}
	databaseVersion, _ = getMajorVersion(config.Value)
	if databaseVersion != 0 {
		return databaseVersion, nil
	}
	return 0, errors.New("invalid version")
}

func migrateVersion() error {
	var databaseVersion uint64

	if !Orm.Migrator().HasTable(&Config{}) {
		initializingDatabase()
	}

	for {
		databaseVersion, err = getDatabaseVersion()
		if err != nil {
			return err
		}

		switch databaseVersion {
		case version:
			return nil
		case 3:
			err = migrateVersion3To4()
			if err != nil {
				return err
			}
			continue
		case 2:
			err = migrateVersion2To3()
			if err != nil {
				return err
			}
			continue
		case 1:
			err = migrateVersion1To2()
			if err != nil {
				return err
			}
			continue
		case 0:
			err = migrateVersion0To1()
			if err != nil {
				return err
			}
			continue
		default:
			return errors.New("invalid version")
		}
	}
}

func initializingDatabase() {
	log.Println("No data in database, initializing")

	tx := Orm.Begin()

	err = tx.AutoMigrate(&Item{}, &Config{})
	if err != nil {
		log.Fatal(err)
	}
	err = tx.Create(&Config{Key: "version", Value: strconv.Itoa(version)}).Error
	if err != nil {
		tx.Rollback()
		log.Fatal(err)
	}
	tx.Commit()
}

func migrateVersion3To4() error {
	log.Println("Migrating to version 4")
	tx := Orm.Begin()
	err = tx.Migrator().RenameColumn(&Item{}, "ItemTime", "Time")
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Migrator().RenameColumn(&Item{}, "ItemData", "Data")
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Migrator().RenameTable(&Item{}, &Item{})
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func migrateVersion2To3() error {
	tx := Orm.Begin()
	err = tx.Migrator().DropColumn(&Item{}, "index")
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Migrator().RenameColumn(&Item{}, "id", "index")
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Save(&Config{Key: "version", Value: "3.0.0"}).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func migrateVersion1To2() error {
	tx := Orm.Begin()

	err := tx.Exec(createFts5TableQuery).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Exec(insertFts5TableQuery).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Save(&Config{Key: "version", Value: "2.0.0"}).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func migrateVersion0To1() error {
	tx := Orm.Begin()

	err := tx.Create(&Config{Key: "version", Value: "1.0.0"}).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

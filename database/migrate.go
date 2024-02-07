package database

import (
	"encoding/base64"
	"errors"
	"log"
	"strconv"
)

type ClipboardItem struct {
	Index             int64  `gorm:"primaryKey"`
	ClipboardItemTime int64  `json:"ItemTime" binding:"required"` // unix milliseconds timestamp
	ClipboardItemText string `json:"ItemText"`
	ClipboardItemHash string `gorm:"unique" json:"ItemHash"`
	ClipboardItemData string `json:"ItemData"`
}

func migrateVersion() error {
	var databaseVersion int64

	if !Orm.Migrator().HasTable(&Config{}) {
		err = initializingDatabase()
		if err != nil {
			return err
		}
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

func initializingDatabase() error {
	log.Println("No data in database, initializing")

	tx := Orm.Begin()
	err = tx.AutoMigrate(&Item{}, &Config{})
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Create(&Config{Key: "version", Value: strconv.FormatInt(version, 10)}).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func migrateVersion3To4() error {
	log.Println("Migrating to version 4")

	rows, err := Orm.Model(&ClipboardItem{}).Rows()
	if err != nil {
		return err
	}
	defer rows.Close()

	tx := Orm.Begin()
	err = tx.Migrator().CreateTable(&Item{})
	if err != nil {
		tx.Rollback()
		return err
	}

	for rows.Next() {
		var item ClipboardItem
		err = Orm.ScanRows(rows, &item)
		if err != nil {
			tx.Rollback()
			return err
		}
		data, err := base64.StdEncoding.DecodeString(item.ClipboardItemData)
		if err != nil {
			tx.Rollback()
			return err
		}
		err = tx.Create(&Item{
			Time: item.ClipboardItemTime,
			Data: data,
		}).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	err = tx.Save(&Config{Key: "version", Value: "4"}).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func migrateVersion2To3() error {
	tx := Orm.Begin()

	err = tx.Migrator().DropColumn(&ClipboardItem{}, "index")
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Migrator().RenameColumn(&ClipboardItem{}, "id", "index")
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

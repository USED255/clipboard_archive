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
		log.Println("Initialize database")
		err = InitializeDatabase()
		if err != nil {
			//初始化失败
			return err
		}
	}

	for {
		databaseVersion, err = getDatabaseVersion()
		if err != nil {
			//获取版本失败
			return err
		}

		switch databaseVersion {
		case version:
			return nil
		case 4:
			log.Println("Migrate to version 5")
			err = migrateVersion4To5()
			if err != nil {
				//迁移失败
				return err
			}
		case 3:
			log.Println("Migrate to version 5")
			err = migrateVersion3To5()
			if err != nil {
				//迁移失败
				return err
			}
			continue
		case 2:
			err = migrateVersion2To3()
			if err != nil {
				//迁移失败
				return err
			}
			continue
		case 1:
			err = migrateVersion1To2()
			if err != nil {
				//迁移失败
				return err
			}
			continue
		case 0:
			err = migrateVersion0To1()
			if err != nil {
				//迁移失败
				return err
			}
			continue
		default:
			return errors.New("invalid version")
		}
	}
}

func InitializeDatabase() error {
	tx := Orm.Begin()
	err = tx.AutoMigrate(&Item{}, &Config{})
	if err != nil {
		//建表失败
		tx.Rollback()
		return err
	}
	err = tx.Create(&Config{Key: "version", Value: strconv.FormatInt(version, 10)}).Error
	if err != nil {
		//插入失败
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func migrateVersion4To5() error {
	err = Orm.Save(&Config{Key: "version", Value: "5"}).Error
	if err != nil {
		//更新失败
		return err
	}
	return nil
}

func migrateVersion3To5() error {
	tx := Orm.Begin()
	rows, err := tx.Model(&ClipboardItem{}).Rows()
	if err != nil {
		//Rows失败
		return err
	}
	defer rows.Close()

	err = tx.Migrator().CreateTable(&Item{})
	if err != nil {
		//建表失败
		tx.Rollback()
		return err
	}

	for rows.Next() {
		var item ClipboardItem
		err = tx.ScanRows(rows, &item)
		if err != nil {
			//Scan失败
			tx.Rollback()
			return err
		}
		data, err := base64.StdEncoding.DecodeString(item.ClipboardItemData)
		if err != nil {
			//解码失败
			tx.Rollback()
			return err
		}
		err = tx.Create(&Item{
			Time: item.ClipboardItemTime,
			Data: data,
		}).Error
		if err != nil {
			//插入失败
			tx.Rollback()
			return err
		}
	}

	err = tx.Save(&Config{Key: "version", Value: "5"}).Error
	if err != nil {
		//更新失败
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
		//删列失败
		tx.Rollback()
		return err
	}
	err = tx.Migrator().RenameColumn(&ClipboardItem{}, "id", "index")
	if err != nil {
		//重命名列失败
		tx.Rollback()
		return err
	}
	err = tx.Save(&Config{Key: "version", Value: "3.0.0"}).Error
	if err != nil {
		//更新失败
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
		//查询失败
		tx.Rollback()
		return err
	}
	err = tx.Exec(insertFts5TableQuery).Error
	if err != nil {
		//查询失败
		tx.Rollback()
		return err
	}
	err = tx.Save(&Config{Key: "version", Value: "2.0.0"}).Error
	if err != nil {
		//更新失败
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
		//插入失败
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

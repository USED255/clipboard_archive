package database

import (
	"errors"
	"regexp"
	"strconv"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var err error

const version int64 = 5

func getMajorVersion(version string) (int64, error) {
	var _majorVersion string

	re := regexp.MustCompile(`^(\d+)\.(\d+)\.(\d+)$`)
	if re.MatchString(version) {
		_majorVersion = re.FindStringSubmatch(version)[1]
	} else {
		return 0, errors.New("invalid version")
	}

	majorVersion, err := strconv.ParseUint(_majorVersion, 10, 64)
	if err != nil {
		return 0, err
	}

	return int64(majorVersion), nil
}

func connectDatabase(dns string) error {
	if Orm != nil {
		return errors.New("database already connected")
	}
	Orm, err = gorm.Open(sqlite.Open(dns), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
	//Orm, err = gorm.Open(sqlite.Open(dns), &gorm.Config{})
	if err != nil {
		return err
	}
	return nil
}

type ClipboardItem struct {
	Index             int64  `gorm:"primaryKey"`
	ClipboardItemTime int64  `json:"ItemTime" binding:"required"` // unix milliseconds timestamp
	ClipboardItemText string `json:"ItemText"`
	ClipboardItemHash string `gorm:"unique" json:"ItemHash"`
	ClipboardItemData string `json:"ItemData"`
}

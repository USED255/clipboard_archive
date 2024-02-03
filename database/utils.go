package database

import (
	"errors"
	"regexp"
	"strconv"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var err error

const version = 5

func getMajorVersion(version string) (uint64, error) {
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

	return majorVersion, nil
}

func connectDatabase(dns string) error {
	if Orm != nil {
		return errors.New("database already connected")
	}
	Orm, err = gorm.Open(sqlite.Open(dns), &gorm.Config{})
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

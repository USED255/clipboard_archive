package database

import (
	"errors"
	"regexp"
	"strconv"

	"gorm.io/gorm"
)

var err error

const version int64 = 5

func connectDatabase(dialector gorm.Dialector, config *gorm.Config) error {
	if Orm != nil {
		return errors.New("database already connected")
	}
	Orm, err = gorm.Open(dialector, config)
	if err != nil {
		//比如数据库损坏
		return err
	}
	return nil
}

func getDatabaseVersion() (int64, error) {
	var config Config
	var databaseVersion int64

	Orm.First(&config, "key = ?", "version")
	if config.Key == "" {
		return 0, nil
	}
	databaseVersion, _ = strconv.ParseInt(config.Value, 10, 64)
	if databaseVersion != 0 {
		return databaseVersion, nil
	}
	databaseVersion, _ = getMajorVersion(config.Value)
	if databaseVersion != 0 {
		return databaseVersion, nil
	}
	return 0, errors.New("invalid version")
}

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

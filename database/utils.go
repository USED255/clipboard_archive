package database

import (
	"errors"
	"regexp"
	"strconv"
)

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

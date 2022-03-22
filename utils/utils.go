package utils

import (
	"time"
)

func GetUnixMillisTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

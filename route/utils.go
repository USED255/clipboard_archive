package route

import (
	"time"
)

func GetUnixMillisTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

type jsonItem struct {
	Time int64  `json:"Time" binding:"required"` // unix milliseconds timestamp
	Data string `json:"Data"  binding:"required"`
}

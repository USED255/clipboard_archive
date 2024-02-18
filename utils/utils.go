package utils

import (
	"io"
	"log"
)

var DebugLog = log.New(io.Discard, "", 0)

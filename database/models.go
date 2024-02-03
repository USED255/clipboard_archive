package database

import (
	"time"

	"gorm.io/gorm"
)

type Config struct {
	Key       string `gorm:"primary_key"`
	Value     string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Item struct {
	Time      int64  `gorm:"primaryKey" json:"Time" binding:"required"` // unix milliseconds timestamp
	Data      []byte `json:"Data"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

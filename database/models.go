package database

type Config struct {
	Key   string `gorm:"primary_key"`
	Value string
}

type Item struct {
	Time int64  `gorm:"primaryKey" json:"Time" binding:"required"` // unix milliseconds timestamp
	Data string `json:"Data"`
}

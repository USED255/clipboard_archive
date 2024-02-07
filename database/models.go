package database

type Config struct {
	Key   string `gorm:"primary_key"`
	Value string
}

type Item struct {
	Time int64  `gorm:"primaryKey" json:"Time"` // unix milliseconds timestamp
	Data []byte `json:"Data" binding:"required"`
}

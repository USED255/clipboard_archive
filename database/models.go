package database

type Config struct {
	Key   string `gorm:"primary_key"`
	Value string
}

type ClipboardItem struct {
	Index             int64  `gorm:"primaryKey"`
	ClipboardItemTime int64  `json:"ClipboardItemTime" binding:"required"` // unix milliseconds timestamp
	ClipboardItemText string `json:"ClipboardItemText"`
	ClipboardItemHash string `gorm:"unique" json:"ClipboardItemHash"`
	ClipboardItemData string `json:"ClipboardItemData"`
}

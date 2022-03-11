package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

const version = "1.1.8"

var db *gorm.DB
var err error

type Config struct {
	Key   string `gorm:"primary_key"`
	Value string
}

type ClipboardItem struct {
	gorm.Model               // 弃用或者是可有可无
	ClipboardItemTime int64  `gorm:"unique" json:"ClipboardItemTime"` // unix milliseconds timestamp
	ClipboardItemText string `json:"ClipboardItemText"`
	ClipboardItemHash string `gorm:"unique" json:"ClipboardItemHash"`
	ClipboardItemData string `json:"ClipboardItemData"`
}

func main() {
	connectDatabase()
	migrateVersion()
	//	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	// Private network
	// IPv4 CIDR
	r.SetTrustedProxies([]string{"192.168.0.0/24", "172.16.0.0/12", "10.0.0.0/8"})
	api := r.Group("/api/v1")
	api.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": "pong",
		})
	})
	api.GET("/version", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"version": version,
			"message": fmt.Sprintf("version %s", version),
		})
	})
	api.POST("/ClipboardItem", insertClipboardItem)
	api.GET("/ClipboardItem", getClipboardItem)
	r.Run() // 监听并在 0.0.0.0:8080 上启动服务
}

func insertClipboardItem(c *gin.Context) {
	var item ClipboardItem
	err = c.BindJSON(&item)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Invalid JSON", "error": err.Error()})
		return
	}
	tx := db.Create(&item)
	if tx.Error != nil {
		if tx.Error.Error() == "constraint failed: UNIQUE constraint failed: clipboard_items.clipboard_item_hash (2067)" {
			c.JSON(http.StatusConflict, gin.H{"status": http.StatusConflict, "message": "ClipboardItem already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "Error inserting ClipboardItem", "error": tx.Error.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"status": http.StatusCreated, "message": "ClipboardItem created successfully", "ClipboardItem": item})
}

func getClipboardItem(c *gin.Context) {

	var startTimestamp int64
	var endTimestamp int64
	var limit int

	_start_timestamp := c.Query("startTimestamp")
	_end_timestamp := c.Query("endTimestamp")
	_limit := c.Query("limit")

	items := []ClipboardItem{}

	if _limit == "" {
		limit = 100
	} else {
		limit, err = strconv.Atoi(_limit)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Invalid limit", "error": err.Error()})
			return
		}
	}

	tx := db.Limit(limit).Order("clipboard_item_time desc")

	if _start_timestamp != "" {
		t, err := strconv.Atoi(_end_timestamp)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Invalid startTimestamp", "error": err.Error()})
			return
		}
		startTimestamp = int64(t)
		tx = tx.Where("clipboard_item_time > ?", startTimestamp)
	}

	if _end_timestamp != "" {
		t, err := strconv.Atoi(_end_timestamp)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Invalid endTimestamp", "error": err.Error()})
			return
		}
		endTimestamp = int64(t)
		tx.Where("clipboard_item_time <= ?", endTimestamp)
	}

	tx.Find(&items)
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "Error getting ClipboardItem", "error": tx.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "ClipboardItem found successfully", "ClipboardItem": items})
}

func connectDatabase() {
	var count int64
	db, err = gorm.Open(sqlite.Open("clipboard_archive_backend.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	db.AutoMigrate(&ClipboardItem{}, &Config{})
	db.Model(&Config{}).Count(&count)
	if count == 0 {
		db.Create(&Config{Key: "version", Value: version})
	}
}

func migrateVersion() {
	var config Config
	err = db.First(&config, "version").Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Fatal("数据不一致")
		} else {
			log.Fatal(err)
		}
	}
	switch config.Value {
	case version:
		return
	case "1.1.7":
		log.Fatal("You are kidding me ?")
	default:
		log.Fatal("数据不一致")
	}
}

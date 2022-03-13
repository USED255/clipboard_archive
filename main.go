package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

const version = "1.1.9"

var db *gorm.DB
var err error

type Config struct {
	Key   string `gorm:"primary_key"`
	Value string
}

type ClipboardItem struct {
	gorm.Model               // 可有可无
	ClipboardItemTime int64  `gorm:"unique" json:"ClipboardItemTime"` // unix milliseconds timestamp
	ClipboardItemText string `json:"ClipboardItemText"`
	ClipboardItemHash string `gorm:"unique" json:"ClipboardItemHash"`
	ClipboardItemData string `json:"ClipboardItemData"`
}

func main() {
	log.Println("Welcome 🐱‍🏍")
	connectDatabase()
	migrateVersion()
	go webServer()
	awaitSignalAndExit()
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
	search := c.Query("search")

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
		startTimestamp, err = strconv.ParseInt(_start_timestamp, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Invalid startTimestamp", "error": err.Error()})
			return
		}
		tx.Where("clipboard_item_time <= ?", startTimestamp)
	}

	if _end_timestamp != "" {
		endTimestamp, err = strconv.ParseInt(_end_timestamp, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Invalid endTimestamp", "error": err.Error()})
			return
		}
		tx.Where("clipboard_item_time >= ?", endTimestamp)
	}
	if search != "" {
		log.Println("Searching for: " + search)
		tx.Debug().Table("clipboard_items_fts").Where("clipboard_items_fts MATCH ?", search).Joins("NATURAL JOIN clipboard_items").Scan(&items)
	} else {
		tx.Find(&items)
	}
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
	db.Model(&ClipboardItem{}).Count(&count)
	if count == 0 {
		db.Create(&Config{Key: "version", Value: version})
		Query := `
CREATE VIRTUAL TABLE clipboard_items_fts USING fts5(clipboard_item_time, clipboard_item_text, content = clipboard_items, content_rowid = clipboard_item_time);

CREATE TRIGGER clipboard_items_ai AFTER INSERT ON clipboard_items BEGIN
    INSERT INTO clipboard_items_fts(rowid, clipboard_item_text) VALUES (new.clipboard_item_time, new.clipboard_item_text);
END;

CREATE TRIGGER clipboard_items_ad AFTER DELETE ON clipboard_items BEGIN
    INSERT INTO clipboard_items_fts(clipboard_items_fts, rowid, clipboard_item_text) VALUES('delete', old.clipboard_item_time, old.clipboard_item_text);
END;

CREATE TRIGGER clipboard_items_au AFTER UPDATE ON clipboard_items BEGIN
    INSERT INTO clipboard_items_fts(clipboard_items_fts, rowid, clipboard_item_text) VALUES('delete', old.clipboard_item_time, old.clipboard_item_text);
    INSERT INTO clipboard_items_fts(rowid, clipboard_item_text) VALUES (new.clipboard_item_time, new.clipboard_item_text);
END;
`
		err := db.Exec(Query).Error
		if err != nil {
			log.Fatal(err)
		}
	}
}

func migrateVersion() {
	var config Config
	err = db.First(&config, "key = ?", "version").Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Fatal("数据不一致")
		} else {
			log.Fatal(err)
		}
	}
migrate:
	switch config.Value {
	case version:
		return
	case "1.1.8":
		log.Println("Migrating to 1.1.9")
		Query := `
CREATE VIRTUAL TABLE clipboard_items_fts USING fts5(
clipboard_item_time, clipboard_item_text, content = clipboard_items, content_rowid = clipboard_item_time);

CREATE TRIGGER clipboard_items_ai AFTER INSERT ON clipboard_items BEGIN
    INSERT INTO clipboard_items_fts(rowid, clipboard_item_text) VALUES (new.clipboard_item_time, new.clipboard_item_text);
END;

CREATE TRIGGER clipboard_items_ad AFTER DELETE ON clipboard_items BEGIN
    INSERT INTO clipboard_items_fts(clipboard_items_fts, rowid, clipboard_item_text) VALUES('delete', old.clipboard_item_time, old.clipboard_item_text);
END;

CREATE TRIGGER clipboard_items_au AFTER UPDATE ON clipboard_items BEGIN
    INSERT INTO clipboard_items_fts(clipboard_items_fts, rowid, clipboard_item_text) VALUES('delete', old.clipboard_item_time, old.clipboard_item_text);
    INSERT INTO clipboard_items_fts(rowid, clipboard_item_text) VALUES (new.clipboard_item_time, new.clipboard_item_text);
END;

INSERT INTO clipboard_items_fts (rowid, clipboard_item_text)
SELECT clipboard_items.clipboard_item_time, clipboard_items.clipboard_item_text FROM clipboard_items;

UPDATE configs SET value = '1.1.9' WHERE key = 'version';
`
		tx := db.Begin()
		err := tx.Exec(Query).Error
		if err != nil {
			tx.Rollback()
			log.Fatal("Migration failed: ", err)
		}
		tx.Commit()
		goto migrate
	case "1.1.7":
		log.Fatal("Are you kidding me ?")
	default:
		log.Fatal("数据不一致")
	}
}

func webServer() {
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
	//r.Run(":8080") // 监听并在 0.0.0.0:8080 上启动服务
	r.Run(":8888") // 监听并在 0.0.0.0:8888 上启动服务
}

func awaitSignalAndExit() {
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT)
	<-s
	log.Println("Bey 🐱‍👤")
	os.Exit(0)
}

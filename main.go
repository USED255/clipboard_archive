package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

const version = "3.0.0"
const CreateFts5TableQuery = `
CREATE VIRTUAL TABLE clipboard_items_fts USING fts5(
	clipboard_item_time, 
	clipboard_item_text, 
	content = clipboard_items, 
	content_rowid = clipboard_item_time
);

CREATE TRIGGER clipboard_items_ai AFTER INSERT ON clipboard_items BEGIN
	INSERT INTO clipboard_items_fts(rowid, clipboard_item_text) 
		VALUES (new.clipboard_item_time, new.clipboard_item_text);
END;

CREATE TRIGGER clipboard_items_ad AFTER DELETE ON clipboard_items BEGIN
	INSERT INTO clipboard_items_fts(clipboard_items_fts, rowid, clipboard_item_text) 
		VALUES('delete', old.clipboard_item_time, old.clipboard_item_text);
END;

CREATE TRIGGER clipboard_items_au AFTER UPDATE ON clipboard_items BEGIN
	INSERT INTO clipboard_items_fts(clipboard_items_fts, rowid, clipboard_item_text) 
		VALUES('delete', old.clipboard_item_time, old.clipboard_item_text);
	INSERT INTO clipboard_items_fts(rowid, clipboard_item_text) 
		VALUES (new.clipboard_item_time, new.clipboard_item_text);
END;
`

var db *gorm.DB
var err error

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

func main() {
	bindFlagPtr := flag.String("bind", ":8080", "bind address")
	versionFlagPtr := flag.Bool("v", false, "show version")
	flag.Parse()
	if *versionFlagPtr {
		fmt.Println(version)
		os.Exit(0)
	}

	log.Println("Welcome 🐱‍🏍")
	log.Println("Clipboard Archive Version: ", version)
	connectDatabase("clipboard_archive.db")
	migrateVersion()
	go func() {
		err := setupRouter().Run(*bindFlagPtr)
		if err != nil {
			log.Fatal(err)
		}
	}()
	awaitSignalAndExit()
}

func connectDatabase(dns string) {
	if db != nil {
		log.Fatalf("Database already connected")
	}
	db, err = gorm.Open(sqlite.Open(dns), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
}

func migrateVersion() {
	var _count int64
	var config Config
	var configMajorVersion uint64

	currentMajorVersion, err := getMajorVersion(version)
	if err != nil {
		log.Fatal(err)
	}

	err = db.AutoMigrate(&ClipboardItem{}, &Config{})
	if err != nil {
		log.Fatal(err)
	}

	db.Model(&ClipboardItem{}).Count(&_count)
	count := _count
	db.Model(&Config{}).Count(&_count)
	count = count + _count
	if count == 0 {
		initializingDatabase()
	}

migrate:
	db.First(&config, "key = ?", "version")
	if config.Key == "" {
		configMajorVersion = 0
	} else {
		configMajorVersion, err = getMajorVersion(config.Value)
		if err != nil {
			log.Fatal(err)
		}
	}
	log.Println("Current version: ", config.Value)

	switch configMajorVersion {
	case currentMajorVersion:
		return
	case 2:
		migrateVersion2To3()
		goto migrate
	case 1:
		migrateVersion1To2()
		goto migrate
	case 0:
		migrateVersion0To1()
		goto migrate
	default:
		log.Fatal("Unsupported version: ", config.Value)
	}

}

func initializingDatabase() {
	log.Println("No data in database, initializing")
	tx := db.Begin()
	err = tx.Exec(CreateFts5TableQuery).Error
	if err != nil {
		tx.Rollback()
		log.Fatal(err)
	}
	err = tx.Create(&Config{Key: "version", Value: version}).Error
	if err != nil {
		tx.Rollback()
		log.Fatal(err)
	}
	tx.Commit()
}

func migrateVersion2To3() {
	log.Println("Migrating to 3.0.0")
	tx := db.Begin()
	err := tx.Migrator().DropColumn(&ClipboardItem{}, "index")
	if err != nil {
		tx.Rollback()
		log.Fatal("Migration failed: ", err)
	}
	err = tx.Migrator().RenameColumn(&ClipboardItem{}, "id", "index")
	if err != nil {
		tx.Rollback()
		log.Fatal("Migration failed: ", err)
	}
	err = tx.Save(&Config{Key: "version", Value: "3.0.0"}).Error
	if err != nil {
		tx.Rollback()
		log.Fatal("Migration failed: ", err)
	}
	tx.Commit()
}

func migrateVersion1To2() {
	log.Println("Migrating to 2.0.0")
	Query := `
INSERT INTO clipboard_items_fts (rowid, clipboard_item_text)
SELECT clipboard_items.clipboard_item_time, clipboard_items.clipboard_item_text 
FROM clipboard_items;
`
	tx := db.Begin()
	err := tx.Exec(CreateFts5TableQuery).Error
	if err != nil {
		tx.Rollback()
		log.Fatal("Migration failed: ", err)
	}
	err = tx.Exec(Query).Error
	if err != nil {
		tx.Rollback()
		log.Fatal("Migration failed: ", err)
	}
	err = tx.Save(&Config{Key: "version", Value: "2.0.0"}).Error
	if err != nil {
		tx.Rollback()
		log.Fatal("Migration failed: ", err)
	}
	tx.Commit()
}

func migrateVersion0To1() {
	log.Println("Migrating to 1.0.0")
	tx := db.Begin()
	err := tx.Create(&Config{Key: "version", Value: "1.0.0"}).Error
	if err != nil {
		tx.Rollback()
		log.Fatal("Migration failed: ", err)
	}
	tx.Commit()
}

func setupRouter() *gin.Engine {
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
	api.DELETE("/ClipboardItem/:id", deleteClipboardItem)
	api.GET("/ClipboardItem", getClipboardItem)
	api.GET("/ClipboardItem/:id", takeClipboardItem)
	api.PUT("/ClipboardItem/:id", updateClipboardItem)
	api.GET("/ClipboardItem/count", getClipboardItemCount)
	return r
}

func insertClipboardItem(c *gin.Context) {
	var item ClipboardItem

	err := c.BindJSON(&item)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Invalid JSON",
			"error":   err.Error(),
		})
		return
	}

	tx := db.Create(&item)
	UniqueError := "constraint failed: UNIQUE constraint failed: clipboard_items.clipboard_item_hash (2067)"
	if tx.Error != nil {
		if tx.Error.Error() == UniqueError {
			c.JSON(http.StatusConflict, gin.H{
				"status":  http.StatusConflict,
				"message": "ClipboardItem already exists",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Error inserting ClipboardItem",
			"error":   tx.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":        http.StatusCreated,
		"message":       "ClipboardItem created successfully",
		"ClipboardItem": item,
	})
}

func deleteClipboardItem(c *gin.Context) {
	var item ClipboardItem

	_id := c.Params.ByName("id")
	id, err := strconv.Atoi(_id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Invalid ID",
			"error":   err.Error(),
		})
		return
	}
	err = db.Where("clipboard_item_time = ?", id).First(&item).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  http.StatusNotFound,
				"message": "ClipboardItem not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Error deleting ClipboardItem",
			"error":   err.Error(),
		})
		return
	}
	err = db.Delete(&item, item.Index).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  http.StatusNotFound,
				"message": "ClipboardItem not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Error deleting ClipboardItem",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":            http.StatusOK,
		"message":           "ClipboardItem deleted successfully",
		"ClipboardItemTime": id,
	})
}

func getClipboardItem(c *gin.Context) {
	var startTimestamp int64
	var endTimestamp int64
	var limit int
	var count int64

	_startTimestamp := c.Query("startTimestamp")
	_endTimestamp := c.Query("endTimestamp")
	_limit := c.Query("limit")
	search := c.Query("search")

	requestedForm := gin.H{
		"startTimestamp": _startTimestamp,
		"endTimestamp":   _endTimestamp,
		"limit":          _limit,
		"search":         search,
	}

	items := []ClipboardItem{}

	functionStartTime := getUnixMillisTimestamp()

	if _limit == "" {
		limit = 100
	} else {
		limit, err = strconv.Atoi(_limit)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Invalid limit",
				"error":   err.Error(),
			})
			return
		}
	}

	tx := db.Order("clipboard_item_time desc")

	if _startTimestamp != "" {
		startTimestamp, err = strconv.ParseInt(_startTimestamp, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Invalid startTimestamp",
				"error":   err.Error(),
			})
			return
		}
		tx.Where("clipboard_item_time >= ?", startTimestamp)
	}

	if _endTimestamp != "" {
		endTimestamp, err = strconv.ParseInt(_endTimestamp, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Invalid endTimestamp",
				"error":   err.Error(),
			})
			return
		}
		tx.Where("clipboard_item_time <= ?", endTimestamp)
	}

	if search != "" {
		//log.Println("Searching for: " + search)
		tx.
			Table("clipboard_items_fts").
			Where("clipboard_items_fts MATCH ?", search).
			Joins("NATURAL JOIN clipboard_items").
			Count(&count)
		if tx.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  http.StatusInternalServerError,
				"message": "Error getting ClipboardItem",
				"error":   tx.Error.Error(),
			})
			return
		}
		//tx.Debug()
		tx.Limit(limit).Scan(&items)
	} else {
		tx.Model(&items).Count(&count)
		if tx.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  http.StatusInternalServerError,
				"message": "Error getting ClipboardItem",
				"error":   tx.Error.Error(),
			})
			return
		}
		tx.Limit(limit).Find(&items)
	}

	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Error getting ClipboardItem",
			"error":   tx.Error.Error(),
		})
		return
	}

	functionEndTime := getUnixMillisTimestamp()

	c.JSON(http.StatusOK, gin.H{
		"status":              http.StatusOK,
		"requested_form":      requestedForm,
		"count":               count,
		"function_start_time": functionStartTime,
		"function_end_time":   functionEndTime,
		"message":             "ClipboardItem found successfully",
		"ClipboardItem":       items,
	})
}

func takeClipboardItem(c *gin.Context) {
	var item ClipboardItem

	_id := c.Params.ByName("id")
	id, err := strconv.Atoi(_id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Invalid ID",
			"error":   err.Error(),
		})
		return
	}
	err = db.Where("clipboard_item_time = ?", id).First(&item).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  http.StatusNotFound,
				"message": "ClipboardItem not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Error taking ClipboardItem",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":        http.StatusOK,
		"message":       "ClipboardItem taken successfully",
		"ClipboardItem": item,
	})
}

func updateClipboardItem(c *gin.Context) {
	var item ClipboardItem

	_id := c.Params.ByName("id")
	id, err := strconv.ParseInt(_id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Invalid id",
			"error":   err.Error(),
		})
		return
	}

	err = db.Where("clipboard_item_time = ?", id).First(&item).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  http.StatusNotFound,
				"message": "ClipboardItem not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Error updating ClipboardItem",
			"error":   err.Error(),
		})
		return
	}

	if c.BindJSON(&item) == nil {
		item.ClipboardItemTime = id
		err = db.Save(&item).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  http.StatusInternalServerError,
				"message": "Error updating ClipboardItem",
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":        http.StatusOK,
			"message":       "ClipboardItem updated successfully",
			"ClipboardItem": item,
		})
	}
}

func getClipboardItemCount(c *gin.Context) {
	var count int64

	db.Model(&ClipboardItem{}).Count(&count)
	if db.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": db.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"count":   count,
		"message": fmt.Sprintf("%d items in clipboard", count),
	})
}

func awaitSignalAndExit() {
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT)
	<-s
	log.Println("Bey 🐱‍👤")
	os.Exit(0)
}

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

func getUnixMillisTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

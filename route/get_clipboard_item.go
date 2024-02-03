package route

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/used255/clipboard_archive/v5/database"
	"github.com/used255/clipboard_archive/v5/utils"
)

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

	functionStartTime := utils.GetUnixMillisTimestamp()

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

	tx := database.Orm.Order("clipboard_item_time desc")

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
		//tx.Debug()
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

	functionEndTime := utils.GetUnixMillisTimestamp()

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

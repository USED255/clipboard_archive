package route

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/used255/clipboard_archive/v3/database"
	"gorm.io/gorm"
)

func takeClipboardItem(c *gin.Context) {
	var item ClipboardItem

	_id := c.Params.ByName("id")
	id, err := strconv.ParseInt(_id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Invalid ID",
			"error":   err.Error(),
		})
		return
	}

	err = database.Orm.Where("clipboard_item_time = ?", id).First(&item).Error
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

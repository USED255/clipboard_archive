package route

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/used255/clipboard_archive/v3/database"
)

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

	tx := database.Orm.Create(&item)
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

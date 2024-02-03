package route

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/used255/clipboard_archive/v5/database"
)

func insertItem(c *gin.Context) {
	var item Item

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
				"message": "Item already exists",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Error inserting Item",
			"error":   tx.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  http.StatusCreated,
		"message": "Item created successfully",
		"Item":    item,
	})
}

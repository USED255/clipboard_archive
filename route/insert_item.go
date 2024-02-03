package route

import (
	"encoding/base64"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/used255/clipboard_archive/v5/database"
)

func insertItem(c *gin.Context) {
	var jsonItem jsonItem

	err := c.BindJSON(&jsonItem)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Invalid JSON",
			"error":   err.Error(),
		})
		return
	}
	data, err := base64.StdEncoding.DecodeString(jsonItem.Data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Invalid Data",
			"error":   err.Error(),
		})
		return
	}
	tx := database.Orm.Create(Item{
		Time: jsonItem.Time,
		Data: data,
	})
	if tx.Error != nil {
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
		"Item":    jsonItem,
	})
}

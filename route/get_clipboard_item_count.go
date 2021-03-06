package route

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/used255/clipboard_archive/v3/database"
)

func getClipboardItemCount(c *gin.Context) {
	var count int64

	database.Orm.Model(&ClipboardItem{}).Count(&count)
	if database.Orm.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": database.Orm.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"count":   count,
		"message": fmt.Sprintf("%d items in clipboard", count),
	})
}

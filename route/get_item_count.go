package route

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/used255/clipboard_archive/v5/database"
	"github.com/used255/clipboard_archive/v5/utils"
)

func getItemCount(c *gin.Context) {
	var count int64

	err = database.Orm.Model(&Item{}).Count(&count).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Error getting item count",
			"error":   err.Error(),
		})
		utils.DebugLog.Println(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"count":   count,
		"message": fmt.Sprintf("%d items in clipboard", count),
	})
}

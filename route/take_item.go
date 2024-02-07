package route

import (
	"encoding/base64"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/used255/clipboard_archive/v5/database"
	"github.com/used255/clipboard_archive/v5/utils"
	"gorm.io/gorm"
)

func takeItem(c *gin.Context) {
	var item Item

	time, err := strconv.ParseInt(c.Params.ByName("time"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Invalid ItemTime",
			"error":   err.Error(),
		})
		utils.DebugLog.Println(err)
		return
	}

	err = database.Orm.First(&item, time).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  http.StatusNotFound,
				"message": "Item not found",
			})
			utils.DebugLog.Println(err)
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Error taking Item",
			"error":   err.Error(),
		})
		utils.DebugLog.Println(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Item taken successfully",
		"Item": jsonItem{
			Time: item.Time,
			Data: base64.StdEncoding.EncodeToString([]byte(item.Data)),
		},
	})
}

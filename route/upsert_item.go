package route

import (
	"encoding/base64"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/used255/clipboard_archive/v5/database"
	"github.com/used255/clipboard_archive/v5/utils"
)

func upsertItem(c *gin.Context) {
	var jsonItem struct {
		Data string `json:"Data" binding:"required"`
	}
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
	err = c.BindJSON(&jsonItem)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Invalid JSON",
			"error":   err.Error(),
		})
		utils.DebugLog.Println(err)
		return
	}
	data, err := base64.StdEncoding.DecodeString(jsonItem.Data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Invalid Data",
			"error":   err.Error(),
		})
		utils.DebugLog.Println(err)
		return
	}

	tx := database.Orm.Create(&Item{
		Time: time,
		Data: data,
	})
	if tx.Error != nil {
		if strings.HasPrefix(tx.Error.Error(), "constraint failed") {
			tx := database.Orm.Save(&Item{
				Time: time,
				Data: data,
			})
			if tx.Error != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"status":  http.StatusInternalServerError,
					"message": "Error upserting Item",
					"error":   tx.Error.Error(),
				})
				utils.DebugLog.Println(err)
				return
			}
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Error upserting Item",
			"error":   tx.Error.Error(),
		})
		utils.DebugLog.Println(err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":   http.StatusCreated,
		"message":  "Item created successfully",
		"ItemTime": time,
	})
}

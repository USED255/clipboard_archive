package route

import (
	"encoding/base64"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/used255/clipboard_archive/v5/database"
	"github.com/used255/clipboard_archive/v5/utils"
	"gorm.io/gorm/clause"
)

func upsertItem(c *gin.Context) {
	var json struct {
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
	err = c.BindJSON(&json)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Invalid JSON",
			"error":   err.Error(),
		})
		utils.DebugLog.Println(err)
		return
	}
	data, err := base64.StdEncoding.DecodeString(json.Data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Invalid Data",
			"error":   err.Error(),
		})
		utils.DebugLog.Println(err)
		return
	}

	err = database.Orm.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "time"}},
		DoUpdates: clause.AssignmentColumns([]string{"data"}),
	}).Create(&Item{
		Time: time,
		Data: data,
	}).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Error upserting Item",
			"error":   err.Error(),
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

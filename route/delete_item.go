package route

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/used255/clipboard_archive/v5/database"
	"github.com/used255/clipboard_archive/v5/utils"
)

func deleteItem(c *gin.Context) {
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

	err = database.Orm.Delete(&Item{}, time).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Error deleting Item",
			"error":   err.Error(),
		})
		utils.DebugLog.Println(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   http.StatusOK,
		"message":  "Item deleted successfully",
		"ItemTime": time,
	})
}

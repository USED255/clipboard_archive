package route

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/used255/clipboard_archive/v5/database"
)

func getItem(c *gin.Context) {
	var times []int64

	limit := 10
	request := gin.H{}

	_startTimestamp := c.Query("startTimestamp")
	_endTimestamp := c.Query("endTimestamp")
	_limit := c.Query("limit")

	tx := database.Orm.Begin()

	tx.Order("time desc")

	if _startTimestamp != "" {
		request["startTimestamp"] = _startTimestamp
		startTime, err := strconv.ParseInt(_startTimestamp, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Invalid startTime",
				"error":   err.Error(),
			})
			return
		}
		tx.Where("time >= ?", startTime)
	}

	if _endTimestamp != "" {
		request["endTimestamp"] = _endTimestamp
		endTime, err := strconv.ParseInt(_endTimestamp, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Invalid endTime",
				"error":   err.Error(),
			})
			return
		}
		tx.Where("time <= ?", endTime)
	}

	if _limit != "" {
		request["limit"] = _limit
		limit, err = strconv.Atoi(_limit)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Invalid limit",
				"error":   err.Error(),
			})
			return
		}
	}
	tx.Limit(limit)

	tx.Pluck("time", &times)
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Error getting Item",
			"error":   tx.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":         http.StatusOK,
		"requested_form": request,
		"message":        "Item found successfully",
		"Items":          times,
	})
}

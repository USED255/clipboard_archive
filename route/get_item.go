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

	_startTime := c.Query("startTime")
	_endTime := c.Query("endTime")
	_limit := c.Query("limit")

	tx := database.Orm.Model(&Item{})

	tx.Order("time desc")

	if _startTime != "" {
		request["startTime"] = _startTime
		startTime, err := strconv.ParseInt(_startTime, 10, 64)
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

	if _endTime != "" {
		request["endTime"] = _endTime
		endTime, err := strconv.ParseInt(_endTime, 10, 64)
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
			"message": "Error getting Items",
			"error":   tx.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":         http.StatusOK,
		"requested_form": request,
		"message":        "Items found successfully",
		"Items":          times,
	})
}

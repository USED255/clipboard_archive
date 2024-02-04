package route

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func searchItem(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "search item",
	})
}

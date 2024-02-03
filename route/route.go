package route

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/used255/clipboard_archive/v5/database"
)

var err error

type Item database.Item

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.SetTrustedProxies([]string{"192.168.0.0/24", "172.16.0.0/12", "10.0.0.0/8"}) // Private network

	api := r.Group("/api/v2")
	api.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": "pong",
		})
	})
	api.GET("/version", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"version": database.Version,
			"message": fmt.Sprintf("version %s", database.Version),
		})
	})
	api.POST("/Item", insertItem)
	api.DELETE("/Item/:time", deleteItem)
	api.GET("/Item", getItem)
	api.GET("/Item/:time", takeItem)
	api.PUT("/Item/:time", updateItem)
	api.GET("/Item/count", getItemCount)
	api.GET("/Item/search", searchItem)

	return r
}

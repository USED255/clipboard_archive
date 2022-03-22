package route

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/used255/clipboard_archive/v3/database"
)

var err error

type ClipboardItem database.ClipboardItem

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.SetTrustedProxies([]string{"192.168.0.0/24", "172.16.0.0/12", "10.0.0.0/8"}) // Private network

	api := r.Group("/api/v1")
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
	api.POST("/ClipboardItem", insertClipboardItem)
	api.DELETE("/ClipboardItem/:id", deleteClipboardItem)
	api.GET("/ClipboardItem", getClipboardItem)
	api.GET("/ClipboardItem/:id", takeClipboardItem)
	api.PUT("/ClipboardItem/:id", updateClipboardItem)
	api.GET("/ClipboardItem/count", getClipboardItemCount)

	return r
}

package route

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/used255/clipboard_archive/v3/database"
)

func TestGetClipboardItemCount(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	database.ConnectDatabase("file::memory:?cache=shared")
	database.MigrateVersion()
	r := SetupRouter()
	item := preparationClipboardItem()
	database.Orm.Create(&item)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/ClipboardItem/count", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	expected := gin.H{
		"status":  http.StatusOK,
		"message": "1 items in clipboard",
		"count":   1,
	}
	expected = reloadJSON(expected)
	got := loadJSON(w.Body.String())
	assert.Equal(t, expected, got)

	database.CloseDatabase()
}

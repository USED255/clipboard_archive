package route

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/used255/clipboard_archive/v3/database"
)

func TestDeleteClipboardItem(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	database.ConnectDatabase("file::memory:?cache=shared")
	database.MigrateVersion()
	r := SetupRouter()

	item := preparationClipboardItem()
	database.Orm.Create(&item)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/api/v1/ClipboardItem/%d", item.ClipboardItemTime), nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	expected := gin.H{
		"status":            http.StatusOK,
		"message":           "ClipboardItem deleted successfully",
		"ClipboardItemTime": item.ClipboardItemTime,
	}
	expected = reloadJSON(expected)
	got := loadJSON(w.Body.String())
	assert.Equal(t, expected, got)

	err := database.Orm.Where("clipboard_item_time = ?", item.ClipboardItemTime).First(&item).Error
	assert.Error(t, err)

	database.CloseDatabase()
}

func TestDeleteClipboardItemParamsError(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	database.ConnectDatabase("file::memory:?cache=shared")
	database.MigrateVersion()
	r := SetupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/v1/ClipboardItem/a", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	expected := gin.H{
		"status":  http.StatusBadRequest,
		"message": "Invalid ID",
	}
	expected = reloadJSON(expected)
	got := loadJSON(w.Body.String())
	delete(got, "error")
	assert.Equal(t, expected, got)

	database.CloseDatabase()
}

func TestDeleteClipboardItemNotFoundError(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	database.ConnectDatabase("file::memory:?cache=shared")
	database.MigrateVersion()
	r := SetupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/v1/ClipboardItem/0", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	expected := gin.H{
		"status":  http.StatusNotFound,
		"message": "ClipboardItem not found",
	}
	expected = reloadJSON(expected)
	got := loadJSON(w.Body.String())
	assert.Equal(t, expected, got)

	database.CloseDatabase()
}

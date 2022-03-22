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

func TestTakeClipboardItems(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	database.ConnectDatabase("file::memory:?cache=shared")
	database.MigrateVersion()
	r := SetupRouter()

	item := preparationClipboardItem()
	database.Orm.Create(&item)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", fmt.Sprintf("/api/v1/ClipboardItem/%d", item.ClipboardItemTime), nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	expected := gin.H{
		"status":        http.StatusOK,
		"message":       "ClipboardItem taken successfully",
		"ClipboardItem": item,
	}
	expected = reloadJSON(expected)
	got := loadJSON(w.Body.String())
	assert.Equal(t, expected, got)

	database.CloseDatabase()
}

func TestTakeClipboardItemsParamsError(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	database.ConnectDatabase("file::memory:?cache=shared")
	database.MigrateVersion()
	r := SetupRouter()

	item := preparationClipboardItem()
	database.Orm.Create(&item)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/ClipboardItem/a", nil)
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

func TestTakeClipboardItemsNotFoundError(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	database.ConnectDatabase("file::memory:?cache=shared")
	database.MigrateVersion()
	r := SetupRouter()

	item := preparationClipboardItem()
	database.Orm.Create(&item)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/ClipboardItem/1", nil)
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

func TestTakeClipboardItemsDatabaseError(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	database.ConnectDatabase("file::memory:?cache=shared")
	r := SetupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/ClipboardItem/1", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	expected := gin.H{
		"status":  http.StatusInternalServerError,
		"message": "Error taking ClipboardItem",
	}
	expected = reloadJSON(expected)
	got := loadJSON(w.Body.String())
	delete(got, "error")
	assert.Equal(t, expected, got)

	database.CloseDatabase()
}

package route

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/used255/clipboard_archive/v3/database"
)

func TestUpdateClipboardItem(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	database.Open("file::memory:?cache=shared")
	r := SetupRouter()

	item := preparationClipboardItem()
	database.Orm.Create(&item)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", fmt.Sprintf("/api/v1/ClipboardItem/%d", item.ClipboardItemTime), strings.NewReader(`{"clipboardItemText": "';DROP TABLE clipboard_items;"}`))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	item.ClipboardItemText = `';DROP TABLE clipboard_items;`
	expected := gin.H{
		"status":        http.StatusOK,
		"message":       "ClipboardItem updated successfully",
		"ClipboardItem": item,
	}
	expected = reloadJSON(expected)
	got := loadJSON(w.Body.String())
	assert.Equal(t, expected, got)
	var item2 ClipboardItem
	database.Orm.First(&item2)
	assert.Equal(t, item, item2)

	database.Close()
}

func TestUpdateClipboardItemParamsError(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	database.Open("file::memory:?cache=shared")
	r := SetupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/v1/ClipboardItem/a", strings.NewReader(`{"clipboardItemText": "test"}`))
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

	database.Close()
}

func TestUpdateClipboardItemBindJsonError(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	database.Open("file::memory:?cache=shared")
	r := SetupRouter()

	item := preparationClipboardItem()
	database.Orm.Create(&item)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", fmt.Sprintf("/api/v1/ClipboardItem/%d", item.ClipboardItemTime), strings.NewReader(`a`))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	expected := gin.H{
		"status":  http.StatusBadRequest,
		"message": "Invalid JSON",
	}
	expected = reloadJSON(expected)
	got := loadJSON(w.Body.String())
	delete(got, "error")
	assert.Equal(t, expected, got)

	database.Close()
}

func TestUpdateClipboardItemNotFoundError(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	database.Open("file::memory:?cache=shared")
	r := SetupRouter()
	item := preparationClipboardItem()
	database.Orm.Create(&item)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/v1/ClipboardItem/1", strings.NewReader(`{"clipboardItemText": "test"}`))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	expected := gin.H{
		"status":  http.StatusNotFound,
		"message": "ClipboardItem not found",
	}
	expected = reloadJSON(expected)
	got := loadJSON(w.Body.String())
	assert.Equal(t, expected, got)

	database.Close()
}

func TestUpdateClipboardItemDatabaseError(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	r := SetupRouter()

	database.OpenNoDatabase()
	defer database.Close()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/v1/ClipboardItem/1", strings.NewReader(`{"clipboardItemText": "test"}`))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	expected := gin.H{
		"status":  http.StatusInternalServerError,
		"message": "Error updating ClipboardItem",
	}
	expected = reloadJSON(expected)
	got := loadJSON(w.Body.String())
	delete(got, "error")
	assert.Equal(t, expected, got)

}

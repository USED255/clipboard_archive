package route

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/used255/clipboard_archive/v3/database"
)

func TestInsertClipboardItem(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	database.ConnectDatabase("file::memory:?cache=shared")
	database.MigrateVersion()
	r := SetupRouter()

	item := preparationClipboardItem()
	item.ClipboardItemText = `'; DELETE TABLE clipboard_items; --`
	item_req := clipboardItemToGinH(item)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/ClipboardItem", strings.NewReader(dumpJSON(item_req)))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	item_req["Index"] = 1
	expected := gin.H{
		"status":        http.StatusCreated,
		"message":       "ClipboardItem created successfully",
		"ClipboardItem": item_req,
	}
	expected = reloadJSON(expected)
	got := loadJSON(w.Body.String())
	assert.Equal(t, expected, got)

	item.Index = 1
	var item2 ClipboardItem
	database.Orm.Where("clipboard_item_time = ?", item.ClipboardItemTime).First(&item2)
	assert.Equal(t, item, item2)

	database.CloseDatabase()
}

func TestInsertClipboardItemBindJsonError(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	database.ConnectDatabase("file::memory:?cache=shared")
	database.MigrateVersion()
	r := SetupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/ClipboardItem", strings.NewReader("{}"))
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

	database.CloseDatabase()
}

func TestInsertClipboardItemUniqueError(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	database.ConnectDatabase("file::memory:?cache=shared")
	database.MigrateVersion()
	r := SetupRouter()

	item := preparationClipboardItem()
	item_req := clipboardItemToGinH(item)
	database.Orm.Create(&item)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/ClipboardItem", strings.NewReader(dumpJSON(item_req)))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)

	expected := gin.H{
		"status":  http.StatusConflict,
		"message": "ClipboardItem already exists",
	}
	expected = reloadJSON(expected)
	got := loadJSON(w.Body.String())
	assert.Equal(t, expected, got)

	database.CloseDatabase()
}

func TestInsertClipboardItemDatabaseError(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	database.ConnectDatabase("file::memory:?cache=shared")
	r := SetupRouter()

	item_req := clipboardItemToGinH(preparationClipboardItem())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/ClipboardItem", strings.NewReader(dumpJSON(item_req)))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	expected := gin.H{
		"status":  http.StatusInternalServerError,
		"message": "Error inserting ClipboardItem",
	}
	delete(expected, "error")
	expected = reloadJSON(expected)
	got := loadJSON(w.Body.String())
	delete(got, "error")
	assert.Equal(t, expected, got)

	database.CloseDatabase()
}

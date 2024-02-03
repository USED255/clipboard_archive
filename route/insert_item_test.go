package route

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/used255/clipboard_archive/v5/database"
)

func TestInsertItem(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	database.Open("file::memory:?cache=shared")
	r := SetupRouter()

	item := preparationItem()
	item.ItemText = `'; DELETE TABLE clipboard_items; --`
	itemReq := ItemToGinH(item)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/Item", strings.NewReader(dumpJSON(itemReq)))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	itemReq["Index"] = 1
	expected := gin.H{
		"status":  http.StatusCreated,
		"message": "Item created successfully",
		"Item":    itemReq,
	}
	expected = reloadJSON(expected)
	got := loadJSON(w.Body.String())
	assert.Equal(t, expected, got)

	item.Index = 1
	var item2 Item
	database.Orm.Where("clipboard_item_time = ?", item.ItemTime).First(&item2)
	assert.Equal(t, item, item2)

	database.Close()
}

func TestInsertItemBindJsonError(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	database.Open("file::memory:?cache=shared")
	r := SetupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/Item", strings.NewReader("{}"))
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

func TestInsertItemUniqueError(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	database.Open("file::memory:?cache=shared")
	r := SetupRouter()

	item := preparationItem()
	itemReq := ItemToGinH(item)
	database.Orm.Create(&item)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/Item", strings.NewReader(dumpJSON(itemReq)))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)

	expected := gin.H{
		"status":  http.StatusConflict,
		"message": "Item already exists",
	}
	expected = reloadJSON(expected)
	got := loadJSON(w.Body.String())
	assert.Equal(t, expected, got)

	database.Close()
}

func TestInsertItemDatabaseError(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	r := SetupRouter()

	database.OpenNoDatabase()
	defer database.Close()

	itemReq := ItemToGinH(preparationItem())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/Item", strings.NewReader(dumpJSON(itemReq)))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	expected := gin.H{
		"status":  http.StatusInternalServerError,
		"message": "Error inserting Item",
	}
	delete(expected, "error")
	expected = reloadJSON(expected)
	got := loadJSON(w.Body.String())
	delete(got, "error")
	assert.Equal(t, expected, got)

}

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

func TestGetClipboardItems(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	database.Open("file::memory:?cache=shared")
	r := SetupRouter()

	item := preparationClipboardItem()
	database.Orm.Create(&item)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/ClipboardItem", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	items := []ClipboardItem{}
	items = append(items, item)
	requestedForm := gin.H{
		"startTimestamp": "",
		"endTimestamp":   "",
		"limit":          "",
		"search":         "",
	}
	expected := gin.H{
		"status":         http.StatusOK,
		"requested_form": requestedForm,
		"count":          1,
		"message":        "ClipboardItem found successfully",
		"ClipboardItem":  items,
	}
	expected = reloadJSON(expected)
	got := loadJSON(w.Body.String())
	delete(got, "function_start_time")
	delete(got, "function_end_time")
	assert.Equal(t, expected, got)

	database.Close()
}

func TestGetClipboardItemsStartTimestampQuery(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	database.Open("file::memory:?cache=shared")
	r := SetupRouter()

	item := preparationClipboardItem()
	database.Orm.Create(&item)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/ClipboardItem?startTimestamp=1", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	items := []ClipboardItem{}
	items = append(items, item)
	requestedForm := gin.H{
		"startTimestamp": "1",
		"endTimestamp":   "",
		"limit":          "",
		"search":         "",
	}
	expected := gin.H{
		"status":         http.StatusOK,
		"requested_form": requestedForm,
		"count":          1,
		"message":        "ClipboardItem found successfully",
		"ClipboardItem":  items,
	}
	expected = reloadJSON(expected)
	got := loadJSON(w.Body.String())
	delete(got, "function_start_time")
	delete(got, "function_end_time")
	assert.Equal(t, expected, got)

	database.Close()
}

func TestGetClipboardItemsStartTimestampQueryError(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	database.Open("file::memory:?cache=shared")
	r := SetupRouter()

	item := preparationClipboardItem()
	database.Orm.Create(&item)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/ClipboardItem?startTimestamp=a", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	expected := gin.H{
		"status":  http.StatusBadRequest,
		"message": "Invalid startTimestamp",
	}
	expected = reloadJSON(expected)
	got := loadJSON(w.Body.String())
	delete(got, "error")
	assert.Equal(t, expected, got)

	database.Close()
}

func TestGetClipboardItemsEndTimeStampQuery(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	database.Open("file::memory:?cache=shared")
	r := SetupRouter()

	item := preparationClipboardItem()
	database.Orm.Create(&item)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/ClipboardItem?endTimestamp=1844674407370955161", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	items := []ClipboardItem{}
	items = append(items, item)
	requestedForm := gin.H{
		"startTimestamp": "",
		"endTimestamp":   "1844674407370955161",
		"limit":          "",
		"search":         "",
	}
	expected := gin.H{
		"status":         http.StatusOK,
		"requested_form": requestedForm,
		"count":          1,
		"message":        "ClipboardItem found successfully",
		"ClipboardItem":  items,
	}
	expected = reloadJSON(expected)
	got := loadJSON(w.Body.String())
	delete(got, "function_start_time")
	delete(got, "function_end_time")
	assert.Equal(t, expected, got)

	database.Close()
}

func TestGetClipboardItemsEndTimeStampQueryError(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	database.Open("file::memory:?cache=shared")
	r := SetupRouter()

	item := preparationClipboardItem()
	database.Orm.Create(&item)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/ClipboardItem?endTimestamp=a", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	expected := gin.H{
		"status":  http.StatusBadRequest,
		"message": "Invalid endTimestamp",
	}
	expected = reloadJSON(expected)
	got := loadJSON(w.Body.String())
	delete(got, "error")
	assert.Equal(t, expected, got)

	database.Close()
}

func TestGetClipboardItemsLimitQuery(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	database.Open("file::memory:?cache=shared")
	r := SetupRouter()

	item := preparationClipboardItem()
	database.Orm.Create(&item)
	item2 := preparationClipboardItem()
	item2.ClipboardItemTime = 1
	database.Orm.Create(&item2)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/ClipboardItem?limit=1", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	items := []ClipboardItem{}
	items = append(items, item)
	requestedForm := gin.H{
		"startTimestamp": "",
		"endTimestamp":   "",
		"limit":          "1",
		"search":         "",
	}
	expected := gin.H{
		"status":         http.StatusOK,
		"requested_form": requestedForm,
		"count":          2,
		"message":        "ClipboardItem found successfully",
		"ClipboardItem":  items,
	}
	expected = reloadJSON(expected)
	got := loadJSON(w.Body.String())
	delete(got, "function_start_time")
	delete(got, "function_end_time")
	assert.Equal(t, expected, got)

	database.Close()
}

func TestGetClipboardItemsLimitQueryError(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	database.Open("file::memory:?cache=shared")
	r := SetupRouter()

	item := preparationClipboardItem()
	database.Orm.Create(&item)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/ClipboardItem?limit=a", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	expected := gin.H{
		"status":  http.StatusBadRequest,
		"message": "Invalid limit",
	}
	expected = reloadJSON(expected)
	got := loadJSON(w.Body.String())
	delete(got, "error")
	assert.Equal(t, expected, got)

	database.Close()
}

func TestGetClipboardItemsSearchQuery(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	database.Open("file::memory:?cache=shared")
	r := SetupRouter()

	item := preparationClipboardItem()
	database.Orm.Create(&item)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", fmt.Sprintf("/api/v1/ClipboardItem?search=%s", item.ClipboardItemText), nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	items := []ClipboardItem{}
	items = append(items, item)
	requestedForm := gin.H{
		"startTimestamp": "",
		"endTimestamp":   "",
		"limit":          "",
		"search":         item.ClipboardItemText,
	}
	expected := gin.H{
		"status":         http.StatusOK,
		"requested_form": requestedForm,
		"count":          1,
		"message":        "ClipboardItem found successfully",
		"ClipboardItem":  items,
	}
	expected = reloadJSON(expected)
	got := loadJSON(w.Body.String())
	delete(got, "function_start_time")
	delete(got, "function_end_time")
	assert.Equal(t, expected, got)

	database.Close()
}

func TestGetClipboardItemsAllQuery(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	database.Open("file::memory:?cache=shared")
	r := SetupRouter()

	item := preparationClipboardItem()
	database.Orm.Create(&item)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", fmt.Sprintf("/api/v1/ClipboardItem?startTimestamp=1&endTimestamp=1844674407370955161&limit=1&search=%s", item.ClipboardItemText), nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	items := []ClipboardItem{}
	items = append(items, item)
	requestedForm := gin.H{
		"startTimestamp": "1",
		"endTimestamp":   "1844674407370955161",
		"limit":          "1",
		"search":         item.ClipboardItemText,
	}
	expected := gin.H{
		"status":         http.StatusOK,
		"requested_form": requestedForm,
		"count":          1,
		"message":        "ClipboardItem found successfully",
		"ClipboardItem":  items,
	}
	expected = reloadJSON(expected)
	got := loadJSON(w.Body.String())
	delete(got, "function_start_time")
	delete(got, "function_end_time")
	assert.Equal(t, expected, got)

	database.Close()
}

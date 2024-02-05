package route

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/used255/clipboard_archive/v5/database"
)

func TestGetItems(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	r := SetupRouter()

	database.Open("file::memory:?cache=shared")
	defer database.Close()

	item := preparationItemReflect()
	database.Orm.Create(&item)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v2/Item", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	items := []int64{item.Time}
	requestedForm := gin.H{}
	expected := gin.H{
		"status":         http.StatusOK,
		"requested_form": requestedForm,
		"message":        "Items found successfully",
		"Items":          items,
	}
	expected = reloadJSON(expected)
	got := loadJSON(w.Body.String())

	assert.Equal(t, expected, got)
}

func TestGetItemsStartTimestampQuery(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	database.Open("file::memory:?cache=shared")
	r := SetupRouter()

	item := preparationJsonItem()
	database.Orm.Create(&item)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/Item?startTimestamp=1", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	items := []jsonItem{}
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
		"message":        "Item found successfully",
		"Item":           items,
	}
	expected = reloadJSON(expected)
	got := loadJSON(w.Body.String())
	delete(got, "function_start_time")
	delete(got, "function_end_time")
	assert.Equal(t, expected, got)

	database.Close()
}

func TestGetItemsStartTimestampQueryError(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	database.Open("file::memory:?cache=shared")
	r := SetupRouter()

	item := preparationJsonItem()
	database.Orm.Create(&item)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/Item?startTimestamp=a", nil)
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

func TestGetItemsEndTimeStampQuery(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	database.Open("file::memory:?cache=shared")
	r := SetupRouter()

	item := preparationJsonItem()
	database.Orm.Create(&item)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/Item?endTimestamp=1844674407370955161", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	items := []jsonItem{}
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
		"message":        "Item found successfully",
		"Item":           items,
	}
	expected = reloadJSON(expected)
	got := loadJSON(w.Body.String())
	delete(got, "function_start_time")
	delete(got, "function_end_time")
	assert.Equal(t, expected, got)

	database.Close()
}

func TestGetItemsEndTimeStampQueryError(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	database.Open("file::memory:?cache=shared")
	r := SetupRouter()

	item := preparationJsonItem()
	database.Orm.Create(&item)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/Item?endTimestamp=a", nil)
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

func TestGetItemsLimitQuery(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	database.Open("file::memory:?cache=shared")
	r := SetupRouter()

	item := preparationJsonItem()
	database.Orm.Create(&item)
	item2 := preparationJsonItem()
	item2.Time = 1
	database.Orm.Create(&item2)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/Item?limit=1", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	items := []jsonItem{}
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
		"message":        "Item found successfully",
		"Item":           items,
	}
	expected = reloadJSON(expected)
	got := loadJSON(w.Body.String())
	delete(got, "function_start_time")
	delete(got, "function_end_time")
	assert.Equal(t, expected, got)

	database.Close()
}

func TestGetItemsLimitQueryError(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	database.Open("file::memory:?cache=shared")
	r := SetupRouter()

	item := preparationJsonItem()
	database.Orm.Create(&item)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/Item?limit=a", nil)
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

func TestGetItemsAllQuery(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	database.Open("file::memory:?cache=shared")
	r := SetupRouter()

	item := preparationJsonItem()
	database.Orm.Create(&item)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/Item?startTimestamp=1&endTimestamp=1844674407370955161&limit=1", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	items := []jsonItem{}
	items = append(items, item)
	requestedForm := gin.H{
		"startTimestamp": "1",
		"endTimestamp":   "1844674407370955161",
		"limit":          "1",
	}
	expected := gin.H{
		"status":         http.StatusOK,
		"requested_form": requestedForm,
		"count":          1,
		"message":        "Item found successfully",
		"Item":           items,
	}
	expected = reloadJSON(expected)
	got := loadJSON(w.Body.String())
	delete(got, "function_start_time")
	delete(got, "function_end_time")
	assert.Equal(t, expected, got)

	database.Close()
}

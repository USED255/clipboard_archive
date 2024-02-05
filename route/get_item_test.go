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

func TestGetItemsStartTimeQuery(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	r := SetupRouter()

	database.Open("file::memory:?cache=shared")
	defer database.Close()

	item := preparationItemReflect()
	database.Orm.Create(&item)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v2/Item?startTime=1", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	items := []int64{item.Time}
	requestedForm := gin.H{
		"startTime": "1",
	}
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

func TestGetItemsStartTimeQueryError(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	r := SetupRouter()

	database.Open("file::memory:?cache=shared")
	defer database.Close()

	item := preparationItemReflect()
	database.Orm.Create(&item)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v2/Item?startTime=a", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	expected := gin.H{
		"status":  http.StatusBadRequest,
		"message": "Invalid startTime",
	}
	expected = reloadJSON(expected)
	got := loadJSON(w.Body.String())
	delete(got, "error")

	assert.Equal(t, expected, got)
}

func TestGetItemsEndTimeQuery(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	r := SetupRouter()

	database.Open("file::memory:?cache=shared")
	defer database.Close()

	item := preparationItemReflect()
	database.Orm.Create(&item)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v2/Item?endTime=1844674407370955161", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	items := []int64{item.Time}
	requestedForm := gin.H{
		"endTime": "1844674407370955161",
	}
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

func TestGetItemsEndTimeQueryError(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	r := SetupRouter()

	database.Open("file::memory:?cache=shared")
	defer database.Close()

	item := preparationItemReflect()
	database.Orm.Create(&item)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v2/Item?endTime=a", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	expected := gin.H{
		"status":  http.StatusBadRequest,
		"message": "Invalid endTime",
	}
	expected = reloadJSON(expected)
	got := loadJSON(w.Body.String())
	delete(got, "error")

	assert.Equal(t, expected, got)
}

func TestGetItemsLimitQuery(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	r := SetupRouter()

	database.Open("file::memory:?cache=shared")
	defer database.Close()

	item := preparationItemReflect()
	database.Orm.Create(&item)

	item2 := preparationItemReflect()
	item2.Time = 1
	database.Orm.Create(&item2)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v2/Item?limit=1", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	items := []int64{item.Time}
	requestedForm := gin.H{
		"limit": "1",
	}
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

func TestGetItemsLimitQueryError(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	r := SetupRouter()

	database.Open("file::memory:?cache=shared")
	defer database.Close()

	item := preparationItemReflect()
	database.Orm.Create(&item)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v2/Item?limit=a", nil)
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
}

func TestGetItemsAllQuery(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	r := SetupRouter()

	database.Open("file::memory:?cache=shared")
	defer database.Close()

	item := preparationItemReflect()
	database.Orm.Create(&item)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v2/Item?startTime=1&endTime=1844674407370955161&limit=1", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	items := []int64{item.Time}
	requestedForm := gin.H{
		"startTime": "1",
		"endTime":   "1844674407370955161",
		"limit":     "1",
	}
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

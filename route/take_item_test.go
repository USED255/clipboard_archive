package route

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/used255/clipboard_archive/v5/database"
)

func TestTakeItems(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	r := SetupRouter()

	database.Open("file::memory:?cache=shared")
	defer database.Close()

	item := newJsonItem()
	data, _ := base64.StdEncoding.DecodeString(item.Data)
	database.Orm.Create(&Item{
		Time: item.Time,
		Data: data,
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", fmt.Sprintf("/api/v2/Item/%d", item.Time), nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	expected := gin.H{
		"status":  http.StatusOK,
		"message": "Item taken successfully",
		"Item":    item,
	}
	expected = ginHToGinH(expected)
	got := stringToJson(w.Body.String())

	assert.Equal(t, expected, got)
}

func TestTakeItemsParamsError(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	database.Open("file::memory:?cache=shared")
	r := SetupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v2/Item/a", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	expected := gin.H{
		"status":  http.StatusBadRequest,
		"message": "Invalid ItemTime",
	}
	expected = ginHToGinH(expected)
	got := stringToJson(w.Body.String())
	delete(got, "error")
	assert.Equal(t, expected, got)

	database.Close()
}

func TestTakeItemsNotFoundError(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	r := SetupRouter()

	database.Open("file::memory:?cache=shared")
	defer database.Close()

	item := newJsonItem()
	data, _ := base64.StdEncoding.DecodeString(item.Data)
	database.Orm.Create(&Item{
		Time: item.Time,
		Data: data,
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v2/Item/1", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	expected := gin.H{
		"status":  http.StatusNotFound,
		"message": "Item not found",
	}
	expected = ginHToGinH(expected)
	got := stringToJson(w.Body.String())

	assert.Equal(t, expected, got)
}

func TestTakeItemsDatabaseError(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	r := SetupRouter()

	database.OpenNoDatabase()
	defer database.Close()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v2/Item/1", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	expected := gin.H{
		"status":  http.StatusInternalServerError,
		"message": "Error taking Item",
	}
	expected = ginHToGinH(expected)
	got := stringToJson(w.Body.String())
	delete(got, "error")

	assert.Equal(t, expected, got)
}

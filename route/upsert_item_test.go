package route

import (
	"encoding/base64"
	"fmt"
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
	r := SetupRouter()

	database.Open("file::memory:?cache=shared")
	defer database.Close()

	item := newJsonItem()
	itemReq := jsonItemToGinH(item)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", fmt.Sprintf("/api/v2/Item/%d", item.Time), strings.NewReader(ginHToJson(itemReq)))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	expected := gin.H{
		"status":   http.StatusCreated,
		"message":  "Item created successfully",
		"ItemTime": item.Time,
	}
	expected = ginHToGinH(expected)
	got := stringToJson(w.Body.String())
	assert.Equal(t, expected, got)

	var item2 Item
	database.Orm.Where(&Item{Time: item.Time}).First(&item2)
	data := base64.StdEncoding.EncodeToString(item2.Data)
	item3 := jsonItem{
		Time: item2.Time,
		Data: data,
	}
	assert.Equal(t, item, item3)
}

func TestInsertItemBindJsonError(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	r := SetupRouter()

	database.Open("file::memory:?cache=shared")
	defer database.Close()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", fmt.Sprintf("/api/v2/Item/%d", 1), strings.NewReader("{}"))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	expected := gin.H{
		"status":  http.StatusBadRequest,
		"message": "Invalid JSON",
	}
	expected = ginHToGinH(expected)
	got := stringToJson(w.Body.String())
	delete(got, "error")

	assert.Equal(t, expected, got)
}

func TestInsertItemDatabaseError(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	r := SetupRouter()

	database.OpenNoDatabase()
	defer database.Close()

	itemReq := jsonItemToGinH(newJsonItem())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", fmt.Sprintf("/api/v2/Item/%d", 1), strings.NewReader(ginHToJson(itemReq)))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	expected := gin.H{
		"status":  http.StatusInternalServerError,
		"message": "Error upserting Item",
	}
	expected = ginHToGinH(expected)
	got := stringToJson(w.Body.String())
	delete(got, "error")

	assert.Equal(t, expected, got)
}

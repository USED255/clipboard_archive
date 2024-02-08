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

func TestUpsertItem(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	r := SetupRouter()

	database.Open("file::memory:?cache=shared")
	defer database.Close()

	item := newJsonItem()
	time := item.Time

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", fmt.Sprintf("/api/v2/Item/%d", time), strings.NewReader(ginHToJson(jsonItemToGinH(item))))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	expected := ginHToGinH(gin.H{
		"status":   http.StatusCreated,
		"message":  "Item created successfully",
		"ItemTime": time,
	})
	got := stringToJson(w.Body.String())
	assert.Equal(t, expected, got)

	item.Data = stringToBase64(randString(5))
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("PUT", fmt.Sprintf("/api/v2/Item/%d", time), strings.NewReader(ginHToJson(jsonItemToGinH(item))))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	expected = ginHToGinH(gin.H{
		"status":   http.StatusCreated,
		"message":  "Item created successfully",
		"ItemTime": time,
	})
	got = stringToJson(w.Body.String())

	assert.Equal(t, expected, got)

	var item2 Item
	database.Orm.First(&item2, time)
	gotItem := jsonItem{
		Time: time,
		Data: base64.StdEncoding.EncodeToString(item2.Data),
	}
	assert.Equal(t, item, gotItem)
}

func TestUpsertItemParamsError(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	r := SetupRouter()

	database.Open("file::memory:?cache=shared")
	defer database.Close()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/v2/Item/a", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	expected := ginHToGinH(gin.H{
		"status":  http.StatusBadRequest,
		"message": "Invalid ItemTime",
	})
	got := stringToJson(w.Body.String())
	delete(got, "error")

	assert.Equal(t, expected, got)
}

func TestUpsertItemBindJsonError(t *testing.T) {
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

func TestUpsertItemDecodeError(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	r := SetupRouter()

	database.Open("file::memory:?cache=shared")
	defer database.Close()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", fmt.Sprintf("/api/v2/Item/%d", 1), strings.NewReader("{\"data\":\"a\"}"))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	expected := gin.H{
		"status":  http.StatusBadRequest,
		"message": "Invalid Data",
	}
	expected = ginHToGinH(expected)
	got := stringToJson(w.Body.String())
	delete(got, "error")

	assert.Equal(t, expected, got)
}

func TestUpsertItemDatabaseError(t *testing.T) {
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

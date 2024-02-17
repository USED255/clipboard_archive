package route

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/used255/clipboard_archive/v5/database"
	"github.com/used255/clipboard_archive/v5/utils"
)

func TestDeleteItem(t *testing.T) {
	utils.DebugLog = log.Default()
	gin.SetMode(gin.ReleaseMode)
	r := SetupRouter()

	database.Open("file::memory:?cache=shared")
	defer database.Close()

	item := newItemReflect()
	time := item.Time

	database.Orm.Create(item)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/api/v2/Item/%d", time), nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	expected := ginHToGinH(gin.H{
		"status":   http.StatusOK,
		"message":  "Item deleted successfully",
		"ItemTime": time,
	})
	got := stringToJson(w.Body.String())
	assert.Equal(t, expected, got)

	err := database.Orm.First(&item, time).Error
	assert.Error(t, err)
}

func TestDeleteItemParamsError(t *testing.T) {
	utils.DebugLog = log.Default()
	gin.SetMode(gin.ReleaseMode)
	r := SetupRouter()

	database.Open("file::memory:?cache=shared")
	defer database.Close()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/v2/Item/a", nil)
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

func TestDeleteItemDatabaseError(t *testing.T) {
	utils.DebugLog = log.Default()
	gin.SetMode(gin.ReleaseMode)
	r := SetupRouter()

	database.OpenNoDatabase()
	defer database.Close()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/v2/Item/1", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	expected := gin.H{
		"status":  http.StatusInternalServerError,
		"message": "Error deleting Item",
	}
	expected = ginHToGinH(expected)
	got := stringToJson(w.Body.String())
	delete(got, "error")

	assert.Equal(t, expected, got)
}

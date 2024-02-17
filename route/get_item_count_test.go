package route

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/used255/clipboard_archive/v5/database"
	"github.com/used255/clipboard_archive/v5/utils"
)

func TestGetItemCount(t *testing.T) {
	utils.DebugLog = log.Default()
	gin.SetMode(gin.ReleaseMode)
	r := SetupRouter()

	database.Open("file::memory:?cache=shared")
	defer database.Close()

	database.Orm.Create(newItemReflect())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v2/Item/count", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	expected := ginHToGinH(gin.H{
		"status":  http.StatusOK,
		"message": "1 items in clipboard",
		"count":   1,
	})
	got := stringToJson(w.Body.String())

	assert.Equal(t, expected, got)
}

func TestGetItemCountDatabaseError(t *testing.T) {
	utils.DebugLog = log.Default()
	gin.SetMode(gin.ReleaseMode)
	r := SetupRouter()

	database.OpenNoDatabase()
	defer database.Close()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v2/Item/count", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	expected := gin.H{
		"status":  http.StatusInternalServerError,
		"message": "Error getting item count",
	}
	expected = ginHToGinH(expected)
	got := stringToJson(w.Body.String())
	delete(got, "error")

	assert.Equal(t, expected, got)
}

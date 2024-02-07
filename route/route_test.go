package route

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/used255/clipboard_archive/v5/database"
)

func TestGetPing(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	r := SetupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v2/ping", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	expected := ginHToGinH(gin.H{
		"status":  http.StatusOK,
		"message": "pong",
	})
	got := stringToJson(w.Body.String())

	assert.Equal(t, expected, got)
}

func TestGetVersion(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	r := SetupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v2/version", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	expected := ginHToGinH(gin.H{
		"status":  http.StatusOK,
		"version": database.Version,
		"message": fmt.Sprintf("version %d", database.Version),
	})
	got := stringToJson(w.Body.String())

	assert.Equal(t, expected, got)
}

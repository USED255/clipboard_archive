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

func TestGetPing(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	r := SetupRouter()
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/api/v1/ping", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	expected := gin.H{
		"status":  http.StatusOK,
		"message": "pong",
	}
	expected = reloadJSON(expected)
	got := loadJSON(w.Body.String())
	assert.Equal(t, expected, got)
}

func TestGetVersion(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	r := SetupRouter()
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/api/v1/version", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	expected := gin.H{
		"status":  http.StatusOK,
		"version": database.Version,
		"message": fmt.Sprintf("version %s", database.Version),
	}
	expected = reloadJSON(expected)
	got := loadJSON(w.Body.String())
	assert.Equal(t, expected, got)
}

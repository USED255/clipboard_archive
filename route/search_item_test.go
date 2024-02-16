package route

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestSearchItem(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	r := SetupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v2/Item/search", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

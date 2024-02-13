package route

import (
	"encoding/base64"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestNewJsonItem(t *testing.T) {
	newJsonItem()
}

func TestStringToBase64(t *testing.T) {
	assert.Equal(t, base64.StdEncoding.EncodeToString([]byte("test")), stringToBase64("test"))
}

func TestNewItemReflect(t *testing.T) {
	newItemReflect()
}

func TestGetUnixMillisTimestamp(t *testing.T) {
	ts := getUnixMillisTimestamp()
	assert.True(t, ts > 0)
}

func TestRandString(t *testing.T) {
	assert.Len(t, randString(5), 5)
	a := randString(5)
	b := randString(5)
	assert.NotEqual(t, a, b)
}

func TestJsonItemToGinH(t *testing.T) {
	jsonItemToGinH(newJsonItem())
}

func TestGinHToGinH(t *testing.T) {
	assert.Equal(t, gin.H{}, ginHToGinH(gin.H{}))
}

func TestStringToJson(t *testing.T) {
	assert.Equal(t, gin.H{}, stringToJson("{}"))
}

func TestGinHToJson(t *testing.T) {
	assert.Equal(t, "{}", ginHToJson(gin.H{}))
}

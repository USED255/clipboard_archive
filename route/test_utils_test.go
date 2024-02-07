package route

import (
	"log"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestDumpJSON(t *testing.T) {
	assert.Equal(t, "{}", ginHToJson(gin.H{}))
}
func TestLoadJSON(t *testing.T) {
	assert.Equal(t, gin.H{}, stringToJson("{}"))
}

func TestReloadJSON(t *testing.T) {
	assert.Equal(t, gin.H{}, ginHToGinH(gin.H{}))
}

func TestRandString(t *testing.T) {
	assert.Len(t, randString(5), 5)
	a := randString(5)
	b := randString(5)
	log.Println("a:", a)
	log.Println("b:", b)
	assert.NotEqual(t, a, b)
}

func TestItemToGinH(t *testing.T) {
	jsonItemToGinH(newJsonItem())
}

func TestPreparationItem(t *testing.T) {
	newJsonItem()
}
func TestGetUnixMillisTimestamp(t *testing.T) {
	ts := getUnixMillisTimestamp()
	assert.True(t, ts > 0)
}

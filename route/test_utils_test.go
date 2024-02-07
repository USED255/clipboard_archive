package route

import (
	"log"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestDumpJSON(t *testing.T) {
	assert.Equal(t, "{}", dumpJSON(gin.H{}))
}
func TestLoadJSON(t *testing.T) {
	assert.Equal(t, gin.H{}, loadJSON("{}"))
}

func TestReloadJSON(t *testing.T) {
	assert.Equal(t, gin.H{}, reloadJSON(gin.H{}))
}

func TestToSha256(t *testing.T) {
	s := "The quick brown fox jumps over the lazy dog"
	b := toSha256(s)

	assert.NotEmpty(t, b)
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
	itemToGinH(preparationJsonItem())
}

func TestPreparationItem(t *testing.T) {
	preparationJsonItem()
}

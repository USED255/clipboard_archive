package route

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/rand"

	"github.com/gin-gonic/gin"
)

func randString(l int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, l)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func preparationJsonItem() jsonItem {
	return jsonItem{
		Time: GetUnixMillisTimestamp(),
		Data: toBase64(randString(5)),
	}
}
func preparationItemReflect() *Item {
	return &Item{
		Time: GetUnixMillisTimestamp(),
		Data: []byte(randString(5)),
	}
}

func itemToGinH(s jsonItem) gin.H {
	var c gin.H
	b, _ := json.Marshal(&s)
	_ = json.Unmarshal(b, &c)
	return c
}

func dumpJSON(g gin.H) string {
	b, _ := json.Marshal(g)
	return string(b)
}

func loadJSON(s string) gin.H {
	var g gin.H
	json.Unmarshal([]byte(s), &g)
	return g
}

func reloadJSON(g gin.H) gin.H {
	return loadJSON(dumpJSON(g))
}

func toBase64(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}

func toSha256(s string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(s)))
}

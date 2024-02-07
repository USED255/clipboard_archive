package route

import (
	"encoding/base64"
	"encoding/json"
	"math/rand"
	"time"

	"github.com/gin-gonic/gin"
)

type jsonItem struct {
	Time int64  `json:"Time" binding:"required"` // unix milliseconds timestamp
	Data string `json:"Data"  binding:"required"`
}

func newJsonItem() jsonItem {
	return jsonItem{
		Time: getUnixMillisTimestamp(),
		Data: base64.StdEncoding.EncodeToString([]byte(randString(5))),
	}
}

func newItemReflect() *Item {
	return &Item{
		Time: getUnixMillisTimestamp(),
		Data: []byte(randString(5)),
	}
}

func getUnixMillisTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func randString(l int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, l)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func jsonItemToGinH(s jsonItem) gin.H {
	var c gin.H
	b, _ := json.Marshal(&s)
	_ = json.Unmarshal(b, &c)
	return c
}

func ginHToGinH(g gin.H) gin.H {
	return stringToJson(ginHToJson(g))
}

func stringToJson(s string) gin.H {
	var g gin.H
	json.Unmarshal([]byte(s), &g)
	return g
}

func ginHToJson(g gin.H) string {
	b, _ := json.Marshal(g)
	return string(b)
}

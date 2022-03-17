package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func dumpJson(g gin.H) string {
	b, _ := json.Marshal(g)
	return string(b)
}

func TestDumpJson(t *testing.T) {
	assert.Equal(t, dumpJson(gin.H{
		"status":  http.StatusOK,
		"message": "pong",
	}), dumpJson(gin.H{
		"status":  http.StatusOK,
		"message": "pong",
	}))
	assert.NotEqual(t, dumpJson(gin.H{
		"status":  http.StatusOK,
		"message": "pong",
	}), dumpJson(gin.H{
		"status":  http.StatusOK,
		"message": "pong2",
	}))
	assert.NotEqual(t, dumpJson(gin.H{
		"status":  http.StatusOK,
		"message": "pong",
	}), dumpJson(gin.H{
		"status":  http.StatusOK,
		"message": "pong",
		"extra":   "extra",
	}))
}

func bindJson(s string) gin.H {
	var g gin.H
	json.Unmarshal([]byte(s), &g)
	return g
}

func TestBindJson(t *testing.T) {
	assert.Equal(t, bindJson(dumpJson(gin.H{
		"ClipboardItemTime": 1,
		"ClipboardItemText": "",
		"ClipboardItemHash": "",
		"ClipboardItemData": "",
	})), bindJson(`{"ClipboardItemTime":1,"ClipboardItemText":"","ClipboardItemHash":"","ClipboardItemData":""}`))
}

func abcJson(s string) string {
	return dumpJson(bindJson(s))
}

func TestAbcJson(t *testing.T) {
	assert.Equal(t, abcJson(`{"ClipboardItemTime":1,"ClipboardItemText":"","ClipboardItemHash":"","ClipboardItemData":""}`), abcJson(`{"ClipboardItemTime":1,"ClipboardItemText":"","ClipboardItemHash":"","ClipboardItemData":""}`))
}

func TestConnectDatabase(t *testing.T) {
	connectDatabase("file::memory:?cache=shared")
	if db == nil {
		t.Fatalf("Expected database connection, got nil")
	}
}

func TestMigrateVersion(t *testing.T) {
	migrateVersion()
	var config Config
	db.First(&config, "key = ?", "version")
	if config.Value != version {
		t.Fatalf("Expected %s, got %s", version, config.Value)
	}
}

func TestRouterPing(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	r := setupRouter()
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/api/v1/ping", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, dumpJson(gin.H{
		"status":  http.StatusOK,
		"message": "pong",
	}), w.Body.String())
}

func TestRouterVersion(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	r := setupRouter()
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/api/v1/version", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, dumpJson(gin.H{
		"status":  http.StatusOK,
		"version": version,
		"message": fmt.Sprintf("version %s", version),
	}), w.Body.String())
}

func TestInsertClipboardItem(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	item := gin.H{
		"ClipboardItemTime": 1,
		"ClipboardItemText": "",
		"ClipboardItemHash": "",
		"ClipboardItemData": "",
	}
	r := setupRouter()
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("POST", "/api/v1/ClipboardItem", strings.NewReader(dumpJson(item)))
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
	item["Index"] = 1
	assert.Equal(t, dumpJson(gin.H{
		"status":        http.StatusCreated,
		"message":       "ClipboardItem created successfully",
		"ClipboardItem": item,
	}), abcJson(w.Body.String()))
}

func TestDeleteClipboardItem(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	r := setupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/v1/ClipboardItem/1", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, dumpJson(gin.H{
		"status":            http.StatusOK,
		"message":           "ClipboardItem deleted successfully",
		"ClipboardItemTime": 1,
	}), abcJson(w.Body.String()))

}

func TestGetMajorVersion(t *testing.T) {
	v, err := getMajorVersion("1.2.3")
	if err != nil {
		t.Fatalf("Error: %s", err)
		return
	}
	if v != 1 {
		t.Fatalf("Expected 1, got %d", v)
		return
	}
	v, err = getMajorVersion("0.0.0")
	if err != nil {
		t.Fatalf("Error: %s", err)
		return
	}
	if v != 0 {
		t.Fatalf("Expected 0, got %d", v)
		return
	}
	v, err = getMajorVersion("65535.0.0")
	if err != nil {
		t.Fatalf("Error: %s", err)
		return
	}
	if v != 65535 {
		t.Fatalf("Expected 65535, got %d", v)
		return
	}
	/***
	v, err = getMajorVersion("1.2.3-alpha")
	if err != nil {
		t.Fatalf("Error: %s", err)
		return
	}
	if v != 1 {
		t.Fatalf("Expected 1, got %d", v)
		return
	}
	v, err = getMajorVersion("1.2.3-alpha.1")
	if err != nil {
		t.Fatalf("Error: %s", err)
		return
	}
	if v != 1 {
		t.Fatalf("Expected 1, got %d", v)
		return
	}
	v, err = getMajorVersion("1.2.3-alpha+build")
	if err != nil {
		t.Fatalf("Error: %s", err)
		return
	}
	v, err = getMajorVersion("1.2.3+build")
	if err != nil {
		t.Fatalf("Error: %s", err)
		return
	}
	if v != 1 {
		t.Fatalf("Expected 1, got %d", v)
		return
	}
	***/
	v, err = getMajorVersion("a")
	if err == nil {
		t.Fatalf("Expected error, got %d", v)
		return
	}
	if v != 0 {
		t.Fatalf("Expected 0, got %d", v)
		return
	}
	v, err = getMajorVersion("1.1.1.1")
	if err == nil {
		t.Fatalf("Expected error, got %d", v)
		return
	}
	if v != 0 {
		t.Fatalf("Expected 0, got %d", v)
		return
	}
	v, err = getMajorVersion("-1.0.0")
	if err == nil {
		t.Fatalf("Expected error, got %d", v)
		return
	}
	if v != 0 {
		t.Fatalf("Expected 0, got %d", v)
		return
	}
}

func TestGetUnixMillisTimestamp(t *testing.T) {
	ts := getUnixMillisTimestamp()
	if ts < 0 {
		t.Fatalf("Expected positive number, got %d", ts)
		return
	}
	if ts < int64(time.Millisecond) {
		t.Fatalf("Expected at least 1ms, got %d", ts)
		return
	}
}

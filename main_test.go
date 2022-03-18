package main

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func dumpJSON(g gin.H) string {
	b, _ := json.Marshal(g)
	return string(b)
}

func TestDumpJSON(t *testing.T) {
	assert.Equal(t, `{}`, dumpJSON(gin.H{}))
}

func loadJSON(s string) gin.H {
	var g gin.H
	json.Unmarshal([]byte(s), &g)
	return g
}

func TestLoadJSON(t *testing.T) {
	assert.Equal(t, gin.H{}, loadJSON(`{}`))
}

func reloadJSON(g gin.H) gin.H {
	return loadJSON(dumpJSON(g))
}

func TestReloadJSON(t *testing.T) {
	assert.Equal(t, gin.H{}, reloadJSON(gin.H{}))
}

func toBase64(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}

func TestToBase64(t *testing.T) {
	s := "The quick brown fox jumps over the lazy dog"
	b := toBase64(s)
	if b == "" {
		t.Fatalf("Expected non-empty string, got %s", b)
		return
	}
}

func toSha256(s string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(s)))
}

func TestToSha256(t *testing.T) {
	s := "The quick brown fox jumps over the lazy dog"
	b := toSha256(s)
	if b == "" {
		t.Fatalf("Expected non-empty string, got %s", b)
		return
	}
}

func TestConnectDatabase(t *testing.T) {
	connectDatabase("file::memory:?cache=shared")
	assert.NotNil(t, db)
}

func TestMigrateVersion(t *testing.T) {
	migrateVersion()
	var config Config
	db.First(&config, "key = ?", "version")
	assert.Equal(t, version, config.Value)
}

func TestGetPing(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	r := setupRouter()
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
	r := setupRouter()
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/api/v1/version", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	expected := gin.H{
		"status":  http.StatusOK,
		"version": version,
		"message": fmt.Sprintf("version %s", version),
	}
	expected = reloadJSON(expected)
	got := loadJSON(w.Body.String())
	assert.Equal(t, expected, got)
}

func TestInsertClipboardItem(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	ClipboardItemText := "The quick brown fox jumps over the lazy dog"
	ClipboardItemData := toBase64(ClipboardItemText)
	ClipboardItemHash := toSha256(ClipboardItemData)
	item := gin.H{
		"ClipboardItemTime": 1,
		"ClipboardItemText": ClipboardItemText,
		"ClipboardItemHash": ClipboardItemHash,
		"ClipboardItemData": ClipboardItemData,
	}
	r := setupRouter()
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("POST", "/api/v1/ClipboardItem", strings.NewReader(dumpJSON(item)))
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
	item["Index"] = 1
	expected := gin.H{
		"status":        http.StatusCreated,
		"message":       "ClipboardItem created successfully",
		"ClipboardItem": item,
	}
	expected = reloadJSON(expected)
	got := loadJSON(w.Body.String())
	assert.Equal(t, expected, got)
}

func TestDeleteClipboardItem(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	r := setupRouter()
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("DELETE", "/api/v1/ClipboardItem/1", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	expected := gin.H{
		"status":            http.StatusOK,
		"message":           "ClipboardItem deleted successfully",
		"ClipboardItemTime": 1,
	}
	expected = reloadJSON(expected)
	got := loadJSON(w.Body.String())
	assert.Equal(t, expected, got)
}

func TestGetClipboardItems(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	r := setupRouter()
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/api/v1/ClipboardItem", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestTakeClipboardItems(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	r := setupRouter()
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/api/v1/ClipboardItem/1", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestUpdateClipboardItem(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	r := setupRouter()
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("PUT", "/api/v1/ClipboardItem/1", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetClipboardItemCount(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	r := setupRouter()
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/api/v1/ClipboardItem/count", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
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
	v, err = getMajorVersion("a")
	if err == nil {
		t.Fatalf("Expected error, got %d", v)
		return
	}
	v, err = getMajorVersion("1.1.1.1")
	if err == nil {
		t.Fatalf("Expected error, got %d", v)
		return
	}
	v, err = getMajorVersion("-1.0.0")
	if err == nil {
		t.Fatalf("Expected error, got %d", v)
		return
	}
}

func TestGetUnixMillisTimestamp(t *testing.T) {
	ts := getUnixMillisTimestamp()
	if ts < 0 {
		t.Fatalf("Expected positive number, got %d", ts)
		return
	}
}

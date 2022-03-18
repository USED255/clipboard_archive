package main

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func dumpJSON(g gin.H) string {
	b, _ := json.Marshal(g)
	return string(b)
}

func TestDumpJSON(t *testing.T) {
	assert.Equal(t, "{}", dumpJSON(gin.H{}))
}

func loadJSON(s string) gin.H {
	var g gin.H
	json.Unmarshal([]byte(s), &g)
	return g
}

func TestLoadJSON(t *testing.T) {
	assert.Equal(t, gin.H{}, loadJSON("{}"))
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
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get database handle: %s", err)
	}
	sqlDB.Close()
}

func TestMigrateVersion(t *testing.T) {
	connectDatabase("file::memory:?cache=shared")
	migrateVersion()
	var config Config
	db.First(&config, "key = ?", "version")
	assert.Equal(t, version, config.Value)
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get database handle: %s", err)
	}
	sqlDB.Close()
}

func TestMigrateVersion1(t *testing.T) {
	connectDatabase("file::memory:?cache=shared")

	Query := `
	CREATE TABLE "clipboard_items" ("id" integer,"created_at" datetime,"updated_at" datetime,"deleted_at" datetime,"clipboard_item_time" integer UNIQUE,"clipboard_item_text" text,"clipboard_item_hash" text UNIQUE,"clipboard_item_data" text,PRIMARY KEY ("id"))
	`
	db.Exec(Query)

	Query = `
	INSERT INTO "main"."clipboard_items" ("id", "created_at", "updated_at", "deleted_at", "clipboard_item_time", "clipboard_item_text", "clipboard_item_hash", "clipboard_item_data") VALUES ("499", "2022-03-13 13:22:43.238644233+08:00", "2022-03-13 13:22:43.238644233+08:00", "", "1647146952858", "migrate", "2cb5fed12b27c377de172eb922161838b1343adf55dbd9db39aa50391f1fc2c7", "/////gAAAAQAAAAaADkAeAAtAGMAbwBwAHkAcQAtAHQAYQBnAHMAAAAAFSwgMjAyMi0wMy0xMyAxMjo0OToxMgAAAC4AOQB4AC0AYwBvAHAAeQBxAC0AdQBzAGUAcgAtAGMAbwBwAHkALQB0AGkAbQBlAAAAAA0xNjQ3MTQ2OTUyODU4AAAACgA4AGgAdABtAGwAAAABLzxodG1sPgo8Ym9keT4KPCEtLVN0YXJ0RnJhZ21lbnQtLT48ZGl2IHN0eWxlPSJjb2xvcjogIzM1MzUzNTtiYWNrZ3JvdW5kLWNvbG9yOiAjZjhmOGY4O2ZvbnQtZmFtaWx5OiBDb25zb2xhcywgJ0NvdXJpZXIgTmV3JywgbW9ub3NwYWNlO2ZvbnQtd2VpZ2h0OiBub3JtYWw7Zm9udC1zaXplOiAxNHB4O2xpbmUtaGVpZ2h0OiAxOXB4O3doaXRlLXNwYWNlOiBwcmU7Ij48ZGl2PjxzcGFuIHN0eWxlPSJjb2xvcjogIzg0MzFjNTsiPm1pZ3JhdGU8L3NwYW4+PC9kaXY+PC9kaXY+PCEtLUVuZEZyYWdtZW50LS0+CjwvYm9keT4KPC9odG1sPgAAAAwAOABwAGwAYQBpAG4AAAAAB21pZ3JhdGU=");
	`
	db.Exec(Query)

	migrateVersion()
	var config Config
	db.First(&config, "key = ?", "version")
	assert.Equal(t, version, config.Value)
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get database handle: %s", err)
	}
	sqlDB.Close()
}

func TestMigrateVersion2(t *testing.T) {
	connectDatabase("file::memory:?cache=shared")

	Query := `
	CREATE TABLE "clipboard_items" ("id" integer,"created_at" datetime,"updated_at" datetime,"deleted_at" datetime,"clipboard_item_time" integer UNIQUE,"clipboard_item_text" text,"clipboard_item_hash" text UNIQUE,"clipboard_item_data" text,PRIMARY KEY ("id"))
	`
	db.Exec(Query)

	Query = `
	CREATE TABLE "configs" ("key" text,"value" text,PRIMARY KEY ("key"))
	`
	db.Exec(Query)

	Query = `
	INSERT INTO "main"."clipboard_items" ("id", "created_at", "updated_at", "deleted_at", "clipboard_item_time", "clipboard_item_text", "clipboard_item_hash", "clipboard_item_data") VALUES ("499", "2022-03-13 13:22:43.238644233+08:00", "2022-03-13 13:22:43.238644233+08:00", "", "1647146952858", "migrate", "2cb5fed12b27c377de172eb922161838b1343adf55dbd9db39aa50391f1fc2c7", "/////gAAAAQAAAAaADkAeAAtAGMAbwBwAHkAcQAtAHQAYQBnAHMAAAAAFSwgMjAyMi0wMy0xMyAxMjo0OToxMgAAAC4AOQB4AC0AYwBvAHAAeQBxAC0AdQBzAGUAcgAtAGMAbwBwAHkALQB0AGkAbQBlAAAAAA0xNjQ3MTQ2OTUyODU4AAAACgA4AGgAdABtAGwAAAABLzxodG1sPgo8Ym9keT4KPCEtLVN0YXJ0RnJhZ21lbnQtLT48ZGl2IHN0eWxlPSJjb2xvcjogIzM1MzUzNTtiYWNrZ3JvdW5kLWNvbG9yOiAjZjhmOGY4O2ZvbnQtZmFtaWx5OiBDb25zb2xhcywgJ0NvdXJpZXIgTmV3JywgbW9ub3NwYWNlO2ZvbnQtd2VpZ2h0OiBub3JtYWw7Zm9udC1zaXplOiAxNHB4O2xpbmUtaGVpZ2h0OiAxOXB4O3doaXRlLXNwYWNlOiBwcmU7Ij48ZGl2PjxzcGFuIHN0eWxlPSJjb2xvcjogIzg0MzFjNTsiPm1pZ3JhdGU8L3NwYW4+PC9kaXY+PC9kaXY+PCEtLUVuZEZyYWdtZW50LS0+CjwvYm9keT4KPC9odG1sPgAAAAwAOABwAGwAYQBpAG4AAAAAB21pZ3JhdGU=");
	`
	db.Exec(Query)

	Query = `
	INSERT INTO "main"."configs" ("key", "value") VALUES ("version", "1.0.0");
	`
	db.Exec(Query)

	migrateVersion()
	var config Config
	db.First(&config, "key = ?", "version")
	assert.Equal(t, version, config.Value)
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get database handle: %s", err)
	}
	sqlDB.Close()
}

func TestMigrateVersion3(t *testing.T) {
	connectDatabase("file::memory:?cache=shared")

	Query := `
	CREATE TABLE "clipboard_items" ("id" integer,"created_at" datetime,"updated_at" datetime,"deleted_at" datetime,"clipboard_item_time" integer UNIQUE,"clipboard_item_text" text,"clipboard_item_hash" text UNIQUE,"clipboard_item_data" text,PRIMARY KEY ("id"))
	`
	db.Exec(Query)

	Query = `
	CREATE TABLE "configs" ("key" text,"value" text,PRIMARY KEY ("key"))
	`
	db.Exec(Query)

	Query = `
	INSERT INTO "main"."clipboard_items" ("id", "created_at", "updated_at", "deleted_at", "clipboard_item_time", "clipboard_item_text", "clipboard_item_hash", "clipboard_item_data") VALUES ("499", "2022-03-13 13:22:43.238644233+08:00", "2022-03-13 13:22:43.238644233+08:00", "", "1647146952858", "migrate", "2cb5fed12b27c377de172eb922161838b1343adf55dbd9db39aa50391f1fc2c7", "/////gAAAAQAAAAaADkAeAAtAGMAbwBwAHkAcQAtAHQAYQBnAHMAAAAAFSwgMjAyMi0wMy0xMyAxMjo0OToxMgAAAC4AOQB4AC0AYwBvAHAAeQBxAC0AdQBzAGUAcgAtAGMAbwBwAHkALQB0AGkAbQBlAAAAAA0xNjQ3MTQ2OTUyODU4AAAACgA4AGgAdABtAGwAAAABLzxodG1sPgo8Ym9keT4KPCEtLVN0YXJ0RnJhZ21lbnQtLT48ZGl2IHN0eWxlPSJjb2xvcjogIzM1MzUzNTtiYWNrZ3JvdW5kLWNvbG9yOiAjZjhmOGY4O2ZvbnQtZmFtaWx5OiBDb25zb2xhcywgJ0NvdXJpZXIgTmV3JywgbW9ub3NwYWNlO2ZvbnQtd2VpZ2h0OiBub3JtYWw7Zm9udC1zaXplOiAxNHB4O2xpbmUtaGVpZ2h0OiAxOXB4O3doaXRlLXNwYWNlOiBwcmU7Ij48ZGl2PjxzcGFuIHN0eWxlPSJjb2xvcjogIzg0MzFjNTsiPm1pZ3JhdGU8L3NwYW4+PC9kaXY+PC9kaXY+PCEtLUVuZEZyYWdtZW50LS0+CjwvYm9keT4KPC9odG1sPgAAAAwAOABwAGwAYQBpAG4AAAAAB21pZ3JhdGU=");
	`

	db.Exec(Query)

	Query = `
	INSERT INTO "main"."configs" ("key", "value") VALUES ("version", "2.0.0");
	`
	db.Exec(Query)

	Query = `
	CREATE VIRTUAL TABLE clipboard_items_fts USING fts5(
		clipboard_item_time, 
		clipboard_item_text, 
		content = clipboard_items, 
		content_rowid = clipboard_item_time
	);

	CREATE TRIGGER clipboard_items_ai AFTER INSERT ON clipboard_items BEGIN
		INSERT INTO clipboard_items_fts(rowid, clipboard_item_text) 
			VALUES (new.clipboard_item_time, new.clipboard_item_text);
	END;

	CREATE TRIGGER clipboard_items_ad AFTER DELETE ON clipboard_items BEGIN
		INSERT INTO clipboard_items_fts(clipboard_items_fts, rowid, clipboard_item_text) 
			VALUES("delete", old.clipboard_item_time, old.clipboard_item_text);
	END;

	CREATE TRIGGER clipboard_items_au AFTER UPDATE ON clipboard_items BEGIN
		INSERT INTO clipboard_items_fts(clipboard_items_fts, rowid, clipboard_item_text) 
			VALUES("delete", old.clipboard_item_time, old.clipboard_item_text);
		INSERT INTO clipboard_items_fts(rowid, clipboard_item_text) 
			VALUES (new.clipboard_item_time, new.clipboard_item_text);
	END;
	`
	db.Exec(Query)

	Query = `
	INSERT INTO clipboard_items_fts (rowid, clipboard_item_text)
	SELECT clipboard_items.clipboard_item_time, clipboard_items.clipboard_item_text 
	FROM clipboard_items;
	`
	db.Exec(Query)

	migrateVersion()
	var config Config
	db.First(&config, "key = ?", "version")
	assert.Equal(t, version, config.Value)
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get database handle: %s", err)
	}
	sqlDB.Close()
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
}

func TestGetMajorVersionError(t *testing.T) {
	v, err := getMajorVersion("a")
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
	v, err = getMajorVersion("184467440737095516150.0.0")
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

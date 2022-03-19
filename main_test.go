package main

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/tjarratt/babble"

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

func closeDatabase() {
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get database handle: %s", err)
	}
	sqlDB.Close()
	db = nil
}

func TestCloseDatabase(t *testing.T) {
	connectDatabase("file::memory:?cache=shared")
	closeDatabase()
}

func TestConnectDatabase(t *testing.T) {
	connectDatabase("file::memory:?cache=shared")
	assert.NotNil(t, db)

	closeDatabase()
}

func TestMigrateVersion(t *testing.T) {
	var config Config
	connectDatabase("file::memory:?cache=shared")
	migrateVersion()
	migrateVersion()
	db.First(&config, "key = ?", "version")
	assert.Equal(t, version, config.Value)

	closeDatabase()
}

func TestMigrateVersionInitializingDatabase(t *testing.T) {
	var config Config
	connectDatabase("file::memory:?cache=shared")
	migrateVersion()
	db.First(&config, "key = ?", "version")
	assert.Equal(t, version, config.Value)

	closeDatabase()
}

func createVersion0Database() {
	Query := `
	CREATE TABLE "clipboard_items" ("id" integer,"created_at" datetime,"updated_at" datetime,"deleted_at" datetime,"clipboard_item_time" integer UNIQUE,"clipboard_item_text" text,"clipboard_item_hash" text UNIQUE,"clipboard_item_data" text,PRIMARY KEY ("id"))
	`
	db.Exec(Query)

	Query = `
	INSERT INTO "main"."clipboard_items" ("id", "created_at", "updated_at", "deleted_at", "clipboard_item_time", "clipboard_item_text", "clipboard_item_hash", "clipboard_item_data") VALUES ("499", "2022-03-13 13:22:43.238644233+08:00", "2022-03-13 13:22:43.238644233+08:00", "", "1647146952858", "migrate", "2cb5fed12b27c377de172eb922161838b1343adf55dbd9db39aa50391f1fc2c7", "/////gAAAAQAAAAaADkAeAAtAGMAbwBwAHkAcQAtAHQAYQBnAHMAAAAAFSwgMjAyMi0wMy0xMyAxMjo0OToxMgAAAC4AOQB4AC0AYwBvAHAAeQBxAC0AdQBzAGUAcgAtAGMAbwBwAHkALQB0AGkAbQBlAAAAAA0xNjQ3MTQ2OTUyODU4AAAACgA4AGgAdABtAGwAAAABLzxodG1sPgo8Ym9keT4KPCEtLVN0YXJ0RnJhZ21lbnQtLT48ZGl2IHN0eWxlPSJjb2xvcjogIzM1MzUzNTtiYWNrZ3JvdW5kLWNvbG9yOiAjZjhmOGY4O2ZvbnQtZmFtaWx5OiBDb25zb2xhcywgJ0NvdXJpZXIgTmV3JywgbW9ub3NwYWNlO2ZvbnQtd2VpZ2h0OiBub3JtYWw7Zm9udC1zaXplOiAxNHB4O2xpbmUtaGVpZ2h0OiAxOXB4O3doaXRlLXNwYWNlOiBwcmU7Ij48ZGl2PjxzcGFuIHN0eWxlPSJjb2xvcjogIzg0MzFjNTsiPm1pZ3JhdGU8L3NwYW4+PC9kaXY+PC9kaXY+PCEtLUVuZEZyYWdtZW50LS0+CjwvYm9keT4KPC9odG1sPgAAAAwAOABwAGwAYQBpAG4AAAAAB21pZ3JhdGU=");
	`
	db.Exec(Query)
}

func TestMigrateVersion0To1(t *testing.T) {
	var config Config
	connectDatabase("file::memory:?cache=shared")
	createVersion0Database()
	migrateVersion()
	db.First(&config, "key = ?", "version")
	assert.Equal(t, version, config.Value)

	closeDatabase()
}

func createVersion1Database() {
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
}

func TestMigrateVersion1To2(t *testing.T) {
	var config Config
	connectDatabase("file::memory:?cache=shared")
	createVersion1Database()
	migrateVersion()
	db.First(&config, "key = ?", "version")
	assert.Equal(t, version, config.Value)

	closeDatabase()
}

func createVersion2Database() {
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
}

func TestMigrateVersion2To3(t *testing.T) {
	var config Config
	connectDatabase("file::memory:?cache=shared")
	createVersion2Database()
	migrateVersion()
	db.First(&config, "key = ?", "version")
	assert.Equal(t, version, config.Value)

	closeDatabase()
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

func preparationClipboardItem() ClipboardItem {
	babbler := babble.NewBabbler()
	babbler.Separator = " "
	ClipboardItemText := babbler.Babble()
	ClipboardItemTime := getUnixMillisTimestamp()
	ClipboardItemData := toBase64(ClipboardItemText)
	ClipboardItemHash := toSha256(ClipboardItemData)
	item := ClipboardItem{
		ClipboardItemText: ClipboardItemText,
		ClipboardItemTime: ClipboardItemTime,
		ClipboardItemData: ClipboardItemData,
		ClipboardItemHash: ClipboardItemHash,
	}
	return item
}

func ClipboardItemToGinH(s ClipboardItem) gin.H {
	var c gin.H
	b, _ := json.Marshal(&s)
	_ = json.Unmarshal(b, &c)
	return c
}

func TestInsertClipboardItem(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	connectDatabase("file::memory:?cache=shared")
	migrateVersion()
	item := ClipboardItemToGinH(preparationClipboardItem())
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
	closeDatabase()
}

func TestInsertClipboardItemBindJsonError(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	connectDatabase("file::memory:?cache=shared")
	migrateVersion()
	r := setupRouter()
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("POST", "/api/v1/ClipboardItem", strings.NewReader("{}"))
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	expected := gin.H{
		"status":  http.StatusBadRequest,
		"message": "Invalid JSON",
	}
	expected = reloadJSON(expected)
	got := loadJSON(w.Body.String())
	delete(got, "error")
	assert.Equal(t, expected, got)
	closeDatabase()
}

func TestInsertClipboardItemUniqueError(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	connectDatabase("file::memory:?cache=shared")
	migrateVersion()
	_item := preparationClipboardItem()
	item := ClipboardItemToGinH(_item)
	db.Create(&_item)
	r := setupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/ClipboardItem", strings.NewReader(dumpJSON(item)))
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusConflict, w.Code)
	expected := gin.H{
		"status":  http.StatusConflict,
		"message": "ClipboardItem already exists",
	}
	expected = reloadJSON(expected)
	got := loadJSON(w.Body.String())
	assert.Equal(t, expected, got)
	closeDatabase()
}

func TestInsertClipboardItemDatabaseError(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	connectDatabase("file::memory:?cache=shared")
	item := ClipboardItemToGinH(preparationClipboardItem())
	r := setupRouter()
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("POST", "/api/v1/ClipboardItem", strings.NewReader(dumpJSON(item)))
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	expected := gin.H{
		"status":  http.StatusInternalServerError,
		"message": "Error inserting ClipboardItem",
	}
	delete(expected, "error")
	expected = reloadJSON(expected)
	got := loadJSON(w.Body.String())
	delete(got, "error")
	assert.Equal(t, expected, got)
	closeDatabase()
}

func TestDeleteClipboardItem(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	connectDatabase("file::memory:?cache=shared")
	migrateVersion()
	_item := preparationClipboardItem()
	db.Create(&_item)
	r := setupRouter()
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/api/v1/ClipboardItem/%d", _item.ClipboardItemTime), nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	expected := gin.H{
		"status":            http.StatusOK,
		"message":           "ClipboardItem deleted successfully",
		"ClipboardItemTime": _item.ClipboardItemTime,
	}
	expected = reloadJSON(expected)
	got := loadJSON(w.Body.String())
	assert.Equal(t, expected, got)
	closeDatabase()
}

func TestDeleteClipboardItemParamsError(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	connectDatabase("file::memory:?cache=shared")
	migrateVersion()
	r := setupRouter()
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("DELETE", "/api/v1/ClipboardItem/a", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	expected := gin.H{
		"status":  http.StatusBadRequest,
		"message": "Invalid ID",
	}
	expected = reloadJSON(expected)
	got := loadJSON(w.Body.String())
	delete(got, "error")
	assert.Equal(t, expected, got)
	closeDatabase()
}

func TestDeleteClipboardItemNotFoundError(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	connectDatabase("file::memory:?cache=shared")
	migrateVersion()
	r := setupRouter()
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("DELETE", "/api/v1/ClipboardItem/0", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
	expected := gin.H{
		"status":  http.StatusNotFound,
		"message": "ClipboardItem not found",
	}
	expected = reloadJSON(expected)
	got := loadJSON(w.Body.String())
	assert.Equal(t, expected, got)
	closeDatabase()
}

func TestGetClipboardItems(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	connectDatabase("file::memory:?cache=shared")
	migrateVersion()
	r := setupRouter()
	w := httptest.NewRecorder()
	item := preparationClipboardItem()
	db.Create(&item)

	req, _ := http.NewRequest("GET", "/api/v1/ClipboardItem", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	items := []ClipboardItem{}
	items = append(items, item)
	requestedForm := gin.H{
		"startTimestamp": "",
		"endTimestamp":   "",
		"limit":          "",
		"search":         "",
	}
	expected := gin.H{
		"status":         http.StatusOK,
		"requested_form": requestedForm,
		"count":          1,
		"message":        "ClipboardItem found successfully",
		"ClipboardItem":  items,
	}
	expected = reloadJSON(expected)
	got := loadJSON(w.Body.String())
	delete(got, "function_start_time")
	delete(got, "function_end_time")
	assert.Equal(t, expected, got)
	closeDatabase()
}

func TestGetClipboardItemsStartTimestampQuery(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	connectDatabase("file::memory:?cache=shared")
	migrateVersion()
	r := setupRouter()
	w := httptest.NewRecorder()
	item := preparationClipboardItem()
	db.Create(&item)

	req, _ := http.NewRequest("GET", "/api/v1/ClipboardItem?startTimestamp=1", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	items := []ClipboardItem{}
	items = append(items, item)
	requestedForm := gin.H{
		"startTimestamp": "1",
		"endTimestamp":   "",
		"limit":          "",
		"search":         "",
	}
	expected := gin.H{
		"status":         http.StatusOK,
		"requested_form": requestedForm,
		"count":          1,
		"message":        "ClipboardItem found successfully",
		"ClipboardItem":  items,
	}
	expected = reloadJSON(expected)
	got := loadJSON(w.Body.String())
	delete(got, "function_start_time")
	delete(got, "function_end_time")
	assert.Equal(t, expected, got)
	closeDatabase()
}

func TestGetClipboardItemsStartTimestampQueryError(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	connectDatabase("file::memory:?cache=shared")
	migrateVersion()
	r := setupRouter()
	w := httptest.NewRecorder()
	item := preparationClipboardItem()
	db.Create(&item)

	req, _ := http.NewRequest("GET", "/api/v1/ClipboardItem?startTimestamp=a", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	expected := gin.H{
		"status":  http.StatusBadRequest,
		"message": "Invalid startTimestamp",
	}
	expected = reloadJSON(expected)
	got := loadJSON(w.Body.String())
	delete(got, "error")
	assert.Equal(t, expected, got)
	closeDatabase()
}

func TestGetClipboardItemsEndTimeStampQuery(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	connectDatabase("file::memory:?cache=shared")
	migrateVersion()
	r := setupRouter()
	w := httptest.NewRecorder()
	item := preparationClipboardItem()
	db.Create(&item)

	req, _ := http.NewRequest("GET", "/api/v1/ClipboardItem?endTimestamp=1844674407370955161", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	items := []ClipboardItem{}
	items = append(items, item)
	requestedForm := gin.H{
		"startTimestamp": "",
		"endTimestamp":   "1844674407370955161",
		"limit":          "",
		"search":         "",
	}
	expected := gin.H{
		"status":         http.StatusOK,
		"requested_form": requestedForm,
		"count":          1,
		"message":        "ClipboardItem found successfully",
		"ClipboardItem":  items,
	}
	expected = reloadJSON(expected)
	got := loadJSON(w.Body.String())
	delete(got, "function_start_time")
	delete(got, "function_end_time")
	assert.Equal(t, expected, got)
	closeDatabase()
}

func TestGetClipboardItemsEndTimeStampQueryError(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	connectDatabase("file::memory:?cache=shared")
	migrateVersion()
	r := setupRouter()
	w := httptest.NewRecorder()
	item := preparationClipboardItem()
	db.Create(&item)

	req, _ := http.NewRequest("GET", "/api/v1/ClipboardItem?endTimestamp=a", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	expected := gin.H{
		"status":  http.StatusBadRequest,
		"message": "Invalid endTimestamp",
	}
	expected = reloadJSON(expected)
	got := loadJSON(w.Body.String())
	delete(got, "error")
	assert.Equal(t, expected, got)
	closeDatabase()
}

func TestGetClipboardItemsLimitQuery(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	connectDatabase("file::memory:?cache=shared")
	migrateVersion()
	r := setupRouter()
	w := httptest.NewRecorder()
	item := preparationClipboardItem()
	db.Create(&item)
	item2 := preparationClipboardItem()
	item2.ClipboardItemTime = 1
	db.Create(&item2)

	req, _ := http.NewRequest("GET", "/api/v1/ClipboardItem?limit=1", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	items := []ClipboardItem{}
	items = append(items, item)
	requestedForm := gin.H{
		"startTimestamp": "",
		"endTimestamp":   "",
		"limit":          "1",
		"search":         "",
	}
	expected := gin.H{
		"status":         http.StatusOK,
		"requested_form": requestedForm,
		"count":          2,
		"message":        "ClipboardItem found successfully",
		"ClipboardItem":  items,
	}
	expected = reloadJSON(expected)
	got := loadJSON(w.Body.String())
	delete(got, "function_start_time")
	delete(got, "function_end_time")
	assert.Equal(t, expected, got)
	closeDatabase()
}

func TestGetClipboardItemsLimitQueryError(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	connectDatabase("file::memory:?cache=shared")
	migrateVersion()
	r := setupRouter()
	w := httptest.NewRecorder()
	item := preparationClipboardItem()
	db.Create(&item)

	req, _ := http.NewRequest("GET", "/api/v1/ClipboardItem?limit=a", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	expected := gin.H{
		"status":  http.StatusBadRequest,
		"message": "Invalid limit",
	}
	expected = reloadJSON(expected)
	got := loadJSON(w.Body.String())
	delete(got, "error")
	assert.Equal(t, expected, got)
	closeDatabase()
}

func TestGetClipboardItemsSearchQuery(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	connectDatabase("file::memory:?cache=shared")
	migrateVersion()
	r := setupRouter()
	w := httptest.NewRecorder()
	item := preparationClipboardItem()
	db.Create(&item)

	req, _ := http.NewRequest("GET", fmt.Sprintf("/api/v1/ClipboardItem?search=%s", item.ClipboardItemText), nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	items := []ClipboardItem{}
	items = append(items, item)
	requestedForm := gin.H{
		"startTimestamp": "",
		"endTimestamp":   "",
		"limit":          "",
		"search":         item.ClipboardItemText,
	}
	expected := gin.H{
		"status":         http.StatusOK,
		"requested_form": requestedForm,
		"count":          1,
		"message":        "ClipboardItem found successfully",
		"ClipboardItem":  items,
	}
	expected = reloadJSON(expected)
	got := loadJSON(w.Body.String())
	delete(got, "function_start_time")
	delete(got, "function_end_time")
	assert.Equal(t, expected, got)
	closeDatabase()
}

func TestGetClipboardItemsAllQuery(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	connectDatabase("file::memory:?cache=shared")
	migrateVersion()
	r := setupRouter()
	w := httptest.NewRecorder()
	item := preparationClipboardItem()
	db.Create(&item)

	req, _ := http.NewRequest("GET", fmt.Sprintf("/api/v1/ClipboardItem?startTimestamp=1&endTimestamp=1844674407370955161&limit=1&search=%s", item.ClipboardItemText), nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	items := []ClipboardItem{}
	items = append(items, item)
	requestedForm := gin.H{
		"startTimestamp": "1",
		"endTimestamp":   "1844674407370955161",
		"limit":          "1",
		"search":         item.ClipboardItemText,
	}
	expected := gin.H{
		"status":         http.StatusOK,
		"requested_form": requestedForm,
		"count":          1,
		"message":        "ClipboardItem found successfully",
		"ClipboardItem":  items,
	}
	expected = reloadJSON(expected)
	got := loadJSON(w.Body.String())
	delete(got, "function_start_time")
	delete(got, "function_end_time")
	assert.Equal(t, expected, got)
	closeDatabase()
}

func TestTakeClipboardItems(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	connectDatabase("file::memory:?cache=shared")
	migrateVersion()
	r := setupRouter()
	w := httptest.NewRecorder()
	item := preparationClipboardItem()
	db.Create(&item)

	req, _ := http.NewRequest("GET", fmt.Sprintf("/api/v1/ClipboardItem/%d", item.ClipboardItemTime), nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	expected := gin.H{
		"status":        http.StatusOK,
		"message":       "ClipboardItem taken successfully",
		"ClipboardItem": item,
	}
	expected = reloadJSON(expected)
	got := loadJSON(w.Body.String())
	assert.Equal(t, expected, got)
	closeDatabase()
}

func TestTakeClipboardItemsParamsError(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	connectDatabase("file::memory:?cache=shared")
	migrateVersion()
	r := setupRouter()
	w := httptest.NewRecorder()
	item := preparationClipboardItem()
	db.Create(&item)

	req, _ := http.NewRequest("GET", "/api/v1/ClipboardItem/a", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	expected := gin.H{
		"status":  http.StatusBadRequest,
		"message": "Invalid ID",
	}
	expected = reloadJSON(expected)
	got := loadJSON(w.Body.String())
	delete(got, "error")
	assert.Equal(t, expected, got)
}

func TestUpdateClipboardItem(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	r := setupRouter()
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("PUT", "/api/v1/ClipboardItem/1", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestUpdateClipboardItemParamsError(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	r := setupRouter()
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("PUT", "/api/v1/ClipboardItem/0", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	expected := gin.H{
		"status":  http.StatusBadRequest,
		"message": "invalid params",
	}
	expected = reloadJSON(expected)
	got := loadJSON(w.Body.String())
	assert.Equal(t, expected, got)
}

func TestUpdateClipboardItemBindJsonError(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	r := setupRouter()
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("PUT", "/api/v1/ClipboardItem/1", strings.NewReader("{}"))
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	expected := gin.H{
		"status":  http.StatusBadRequest,
		"message": "invalid params",
	}
	expected = reloadJSON(expected)
	got := loadJSON(w.Body.String())
	assert.Equal(t, expected, got)
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

package main

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
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

	assert.NotEmpty(t, b)
}

func toSha256(s string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(s)))
}

func TestToSha256(t *testing.T) {
	s := "The quick brown fox jumps over the lazy dog"
	b := toSha256(s)

	assert.NotEmpty(t, b)
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

func TestCreateVersion0Database(t *testing.T) {
	var config Config
	connectDatabase("file::memory:?cache=shared")

	createVersion0Database()

	db.First(&config, "key = ?", "version")
	assert.Equal(t, "", config.Value)

	closeDatabase()
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

func TestCreateVersion1Database(t *testing.T) {
	var config Config
	connectDatabase("file::memory:?cache=shared")

	createVersion1Database()

	db.First(&config, "key = ?", "version")
	assert.Equal(t, "1.0.0", config.Value)

	closeDatabase()
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

func TestCreateVersion2Database(t *testing.T) {
	var config Config
	connectDatabase("file::memory:?cache=shared")

	createVersion2Database()

	db.First(&config, "key = ?", "version")
	assert.Equal(t, "2.0.0", config.Value)

	closeDatabase()
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

func randString(l int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, l)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func TestRandString(t *testing.T) {
	assert.Len(t, randString(5), 5)
	a := randString(5)
	b := randString(5)
	log.Println("a:", a)
	log.Println("b:", b)
	assert.NotEqual(t, a, b)
}

func preparationClipboardItem() ClipboardItem {
	ClipboardItemText := randString(5)
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

func TestPreparationClipboardItem(t *testing.T) {
	preparationClipboardItem()
}

func clipboardItemToGinH(s ClipboardItem) gin.H {
	var c gin.H
	b, _ := json.Marshal(&s)
	_ = json.Unmarshal(b, &c)
	return c
}

func TestClipboardItemToGinH(t *testing.T) {
	clipboardItemToGinH(preparationClipboardItem())
}

func TestInsertClipboardItem(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	connectDatabase("file::memory:?cache=shared")
	migrateVersion()
	r := setupRouter()

	item := preparationClipboardItem()
	item.ClipboardItemText = `'; DELETE TABLE clipboard_items; --`
	item_req := clipboardItemToGinH(item)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/ClipboardItem", strings.NewReader(dumpJSON(item_req)))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	item_req["Index"] = 1
	expected := gin.H{
		"status":        http.StatusCreated,
		"message":       "ClipboardItem created successfully",
		"ClipboardItem": item_req,
	}
	expected = reloadJSON(expected)
	got := loadJSON(w.Body.String())
	assert.Equal(t, expected, got)

	item.Index = 1
	var item2 ClipboardItem
	db.Where("clipboard_item_time = ?", item.ClipboardItemTime).First(&item2)
	assert.Equal(t, item, item2)

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
	r := setupRouter()

	item := preparationClipboardItem()
	item_req := clipboardItemToGinH(item)
	db.Create(&item)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/ClipboardItem", strings.NewReader(dumpJSON(item_req)))
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
	r := setupRouter()

	item_req := clipboardItemToGinH(preparationClipboardItem())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/ClipboardItem", strings.NewReader(dumpJSON(item_req)))
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
	r := setupRouter()

	item := preparationClipboardItem()
	db.Create(&item)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/api/v1/ClipboardItem/%d", item.ClipboardItemTime), nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	expected := gin.H{
		"status":            http.StatusOK,
		"message":           "ClipboardItem deleted successfully",
		"ClipboardItemTime": item.ClipboardItemTime,
	}
	expected = reloadJSON(expected)
	got := loadJSON(w.Body.String())
	assert.Equal(t, expected, got)

	err := db.Where("clipboard_item_time = ?", item.ClipboardItemTime).First(&item).Error
	assert.Error(t, err)

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

	item := preparationClipboardItem()
	db.Create(&item)

	w := httptest.NewRecorder()
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

	item := preparationClipboardItem()
	db.Create(&item)

	w := httptest.NewRecorder()
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

	item := preparationClipboardItem()
	db.Create(&item)

	w := httptest.NewRecorder()
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

	item := preparationClipboardItem()
	db.Create(&item)

	w := httptest.NewRecorder()
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

	item := preparationClipboardItem()
	db.Create(&item)

	w := httptest.NewRecorder()
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

	item := preparationClipboardItem()
	db.Create(&item)
	item2 := preparationClipboardItem()
	item2.ClipboardItemTime = 1
	db.Create(&item2)

	w := httptest.NewRecorder()
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

	item := preparationClipboardItem()
	db.Create(&item)

	w := httptest.NewRecorder()
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

	item := preparationClipboardItem()
	db.Create(&item)

	w := httptest.NewRecorder()
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

	item := preparationClipboardItem()
	db.Create(&item)

	w := httptest.NewRecorder()
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

	item := preparationClipboardItem()
	db.Create(&item)

	w := httptest.NewRecorder()
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

	item := preparationClipboardItem()
	db.Create(&item)

	w := httptest.NewRecorder()
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

	closeDatabase()
}

func TestTakeClipboardItemsNotFoundError(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	connectDatabase("file::memory:?cache=shared")
	migrateVersion()
	r := setupRouter()

	item := preparationClipboardItem()
	db.Create(&item)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/ClipboardItem/1", nil)
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

func TestTakeClipboardItemsDatabaseError(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	connectDatabase("file::memory:?cache=shared")
	r := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/ClipboardItem/1", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	expected := gin.H{
		"status":  http.StatusInternalServerError,
		"message": "Error taking ClipboardItem",
	}
	expected = reloadJSON(expected)
	got := loadJSON(w.Body.String())
	delete(got, "error")
	assert.Equal(t, expected, got)

	closeDatabase()
}

func TestUpdateClipboardItem(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	connectDatabase("file::memory:?cache=shared")
	migrateVersion()
	r := setupRouter()

	item := preparationClipboardItem()
	db.Create(&item)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", fmt.Sprintf("/api/v1/ClipboardItem/%d", item.ClipboardItemTime), strings.NewReader(`{"clipboardItemText": "';DROP TABLE clipboard_items;"}`))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	item.ClipboardItemText = `';DROP TABLE clipboard_items;`
	expected := gin.H{
		"status":        http.StatusOK,
		"message":       "ClipboardItem updated successfully",
		"ClipboardItem": item,
	}
	expected = reloadJSON(expected)
	got := loadJSON(w.Body.String())
	assert.Equal(t, expected, got)
	var item2 ClipboardItem
	db.First(&item2)
	assert.Equal(t, item, item2)

	closeDatabase()
}

func TestUpdateClipboardItemParamsError(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	connectDatabase("file::memory:?cache=shared")
	migrateVersion()
	r := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/v1/ClipboardItem/a", strings.NewReader(`{"clipboardItemText": "test"}`))
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

func TestUpdateClipboardItemBindJsonError(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	connectDatabase("file::memory:?cache=shared")
	migrateVersion()
	r := setupRouter()

	item := preparationClipboardItem()
	db.Create(&item)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", fmt.Sprintf("/api/v1/ClipboardItem/%d", item.ClipboardItemTime), strings.NewReader(`a`))
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

func TestUpdateClipboardItemNotFoundError(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	connectDatabase("file::memory:?cache=shared")
	migrateVersion()
	r := setupRouter()
	item := preparationClipboardItem()
	db.Create(&item)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/v1/ClipboardItem/1", strings.NewReader(`{"clipboardItemText": "test"}`))
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

func TestUpdateClipboardItemDatabaseError(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	connectDatabase("file::memory:?cache=shared")
	r := setupRouter()

	item := preparationClipboardItem()
	db.Create(&item)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/v1/ClipboardItem/1", strings.NewReader(`{"clipboardItemText": "test"}`))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	expected := gin.H{
		"status":  http.StatusInternalServerError,
		"message": "Error updating ClipboardItem",
	}
	expected = reloadJSON(expected)
	got := loadJSON(w.Body.String())
	delete(got, "error")
	assert.Equal(t, expected, got)

	closeDatabase()
}

func TestGetClipboardItemCount(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	connectDatabase("file::memory:?cache=shared")
	migrateVersion()
	r := setupRouter()
	item := preparationClipboardItem()
	db.Create(&item)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/ClipboardItem/count", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	expected := gin.H{
		"status":  http.StatusOK,
		"message": "1 items in clipboard",
		"count":   1,
	}
	expected = reloadJSON(expected)
	got := loadJSON(w.Body.String())
	assert.Equal(t, expected, got)

	closeDatabase()
}

func TestGetMajorVersion(t *testing.T) {
	v, err := getMajorVersion("1.2.3")
	assert.NoError(t, err)
	assert.Equal(t, uint64(1), v)

	v, err = getMajorVersion("0.0.0")
	assert.NoError(t, err)
	assert.Equal(t, uint64(0), v)

	v, err = getMajorVersion("65535.0.0")
	assert.NoError(t, err)
	assert.Equal(t, uint64(65535), v)
}

func TestGetMajorVersionError(t *testing.T) {
	v, err := getMajorVersion("a")
	assert.Error(t, err)
	assert.Equal(t, uint64(0), v)

	v, err = getMajorVersion("1.1.1.1")
	assert.Error(t, err)
	assert.Equal(t, uint64(0), v)

	v, err = getMajorVersion("-1.0.0")
	assert.Error(t, err)
	assert.Equal(t, uint64(0), v)

	v, err = getMajorVersion("184467440737095516150.0.0")
	assert.Error(t, err)
	assert.Equal(t, uint64(0), v)
}

func TestGetUnixMillisTimestamp(t *testing.T) {
	ts := getUnixMillisTimestamp()
	assert.True(t, ts > 0)
}

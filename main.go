package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "modernc.org/sqlite"
)

var db *sql.DB
var err error

func main() {
	openDatabase()
	defer db.Close()
	go webServer()
	defer log.Println("再见👋")
	db.Exec(`CREATE VIRTUAL TABLE email USING fts5(sender, title, body)`)
	select {}
}

func initializeDatabase() {
	statement, err := db.Prepare("CREATE TABLE IF NOT EXISTS config (key TEXT PRIMARY KEY, value TEXT)")

	if err != nil {
		log.Fatal(err)
	}

	_, err = statement.Exec()

	if err != nil {
		log.Fatal(err)
	}
	json_string, err := json.Marshal(map[string]interface{}{"a": "a"})
	if err != nil {
		log.Fatal(err)
	}
	query := fmt.Sprintf(`INSERT OR IGNORE INTO config (key, value) VALUES ('test', json('%s'))`, json_string)
	statement, err = db.Prepare(query)
	if err != nil {
		log.Fatal(err)
	}
	_, err = statement.Exec()
	if err != nil {
		log.Fatal(err)
	}
}

func openDatabase() {
	db, err = sql.Open("sqlite", "./database.db")

	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()

	if err != nil {
		log.Fatal(err)
	}

	initializeDatabase()
}

func webServer() {
	// gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	err = r.Run() // listen and serve on
	if err != nil {
		log.Fatal(err)
	}
}

package cmd

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/used255/clipboard_archive/v5/database"
	"github.com/used255/clipboard_archive/v5/route"
	"github.com/used255/clipboard_archive/v5/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var err error

func Start() {
	databasePathFlagPtr := flag.String("database", "clipboard_archive.sqlite3", "database path")
	bindFlagPtr := flag.String("bind", ":8080", "bind address")
	versionFlagPtr := flag.Bool("v", false, "show version")
	debugFlagPtr := flag.Bool("debug", false, "")

	flag.Parse()

	gin.SetMode(gin.ReleaseMode)
	database.OrmConfig = &gorm.Config{}
	utils.DebugLog = log.New(io.Discard, "", 0)

	if *versionFlagPtr {
		fmt.Println(database.Version)
		os.Exit(0)
	}

	if *debugFlagPtr {
		gin.SetMode(gin.DebugMode)
		database.OrmConfig = &gorm.Config{Logger: logger.Default.LogMode(logger.Info)}
		utils.DebugLog = log.Default()
		utils.DebugLog.SetFlags(log.LstdFlags | log.Lshortfile)
	}

	log.Println("Welcome 🐱‍🏍")

	err = database.Open(*databasePathFlagPtr)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		err = route.SetupRouter().Run(*bindFlagPtr)
		if err != nil {
			log.Fatal(err)
		}
	}()
	awaitSignalAndExit()
}

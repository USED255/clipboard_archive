package cmd

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/used255/clipboard_archive/v3/database"
	"github.com/used255/clipboard_archive/v3/route"
)

var err error

func Start() {
	bindFlagPtr := flag.String("bind", ":8080", "bind address")
	versionFlagPtr := flag.Bool("v", false, "show version")
	disableGinModeFlagPtr := flag.Bool("disable-gin-debug-mode", false, "gin.ReleaseMode")

	flag.Parse()

	if *versionFlagPtr {
		fmt.Println(database.Version)
		os.Exit(0)
	}

	if *disableGinModeFlagPtr {
		gin.SetMode(gin.ReleaseMode)
	}

	log.Println("Welcome 🐱‍🏍")
	database.Open("clipboard_archive.db")
	go func() {
		err = route.SetupRouter().Run(*bindFlagPtr)
		if err != nil {
			log.Fatal(err)
		}
	}()
	awaitSignalAndExit()
}

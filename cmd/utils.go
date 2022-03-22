package cmd

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

func awaitSignalAndExit() {
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT)
	<-s
	log.Println("Bey 🐱‍👤")
	os.Exit(0)
}

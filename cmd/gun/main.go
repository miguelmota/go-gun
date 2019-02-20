package main

import (
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/miguelmota/go-gun/server"
)

func main() {
	var debug bool
	if os.Getenv("DEBUG") != "" {
		log.SetReportCaller(true)
		debug = true
	}

	srv := server.NewServer(&server.Config{
		Port:  8080,
		Debug: debug,
	})

	log.Fatal(srv.Start())
}

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

	port := uint(8080)
	srv := server.NewServer(&server.Config{
		Port:  &port,
		Debug: debug,
	})

	log.Fatal(srv.Start())
}

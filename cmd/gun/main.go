package main

import (
	"log"

	"github.com/miguelmota/go-gun/server"
)

func main() {
	srv := server.NewServer(&server.Config{
		Port: 8080,
	})

	log.Fatal(srv.Start())
}

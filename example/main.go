package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.Handle("/", http.FileServer(http.Dir("./example")))
	fmt.Println("listening on port 3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}

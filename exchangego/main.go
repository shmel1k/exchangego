package main

import (
	"log"
	"net/http"

	"github.com/shmel1k/exchangego/config"
	"github.com/shmel1k/exchangego/exchange/auth"
)

func main() {
	http.HandleFunc("/auth", auth.Authorize)

	port := ":" + config.HTTPServer().Port
	log.Printf("Starting listening http server on port %q", port)
	http.ListenAndServe(port, nil)
}

package main

import (
	"log"
	"net/http"

	"github.com/shmel1k/exchangego/config"
	"github.com/shmel1k/exchangego/exchange/auth"
	"github.com/shmel1k/exchangego/exchange/register"
)

func main() {
	http.HandleFunc("/auth", auth.AuthorizeHandler)
	http.HandleFunc("/register", register.RegisterHandler)

	port := ":" + config.HTTPServer().Port
	log.Printf("Starting listening http server on port %q", port)
	http.ListenAndServe(port, nil)
}

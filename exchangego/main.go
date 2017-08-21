package main

import (
	"log"
	"net/http"
	"time"

	"github.com/shmel1k/exchangego/config"
	"github.com/shmel1k/exchangego/exchange/auth"
	"github.com/shmel1k/exchangego/exchange/register"
	"github.com/shmel1k/exchangego/broadcast"
)

var broadCaster *server.EasyCast

func init() {
	log.Println("Set config websocket")
	broadCaster = server.NewEasyCast(1*time.Second, 5)
}

func connectWebSocketHandler(w http.ResponseWriter, r *http.Request) {
	// TODO CheckCookie
	ok := broadCaster.Subscribe(w, r)
	if !ok {
		log.Fatal("cannot subscribe")
	}
}

func main() {
	http.HandleFunc("/auth", auth.AuthorizeHandler)
	http.HandleFunc("/register", register.RegisterHandler)

	http.HandleFunc("/ws", connectWebSocketHandler)

	/* TODO nginx */
	fs := http.FileServer(http.Dir("./exchangego/static"))
	http.Handle("/", fs)

	port := ":" + config.HTTPServer().Port
	log.Printf("Starting listening http server on port %q", port)
	http.ListenAndServe(port, nil)
}

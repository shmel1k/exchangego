package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/shmel1k/exchangego/broadcast"
	"github.com/shmel1k/exchangego/config"
	"github.com/shmel1k/exchangego/currency"
	"github.com/shmel1k/exchangego/exchange/auth"
	"github.com/shmel1k/exchangego/exchange/register"
)

var broadCaster *server.EasyCast

type Currencies struct {
	History *[]int `json:"history"`
}

type Error struct {
	Err string `json:"error"`
}

func init() {
	log.Println("Set websocket")
	currency.InitCurrency()
	broadCaster = server.NewEasyCast(currency.UpdateCurrency, 1*time.Second, 5)
}

func connectWebSocketHandler(w http.ResponseWriter, r *http.Request) {
	// TODO CheckCookie
	ok := broadCaster.Subscribe(w, r)
	if !ok {
		log.Fatal("cannot subscribe")
	}
}

func simpleParam(r *http.Request, key string) (string, bool) {
	value, ok := r.Form[key]
	if !ok || len(value) != 1 {
		return "", false
	}

	return value[0], true
}

func getLastCurrency(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	w.Header().Set("Content-Type", "application/json")

	sizeStr, ok := simpleParam(r, "size")
	if !ok {
		json.NewEncoder(w).Encode(Error{"bad params"})
		return
	}

	size, err := strconv.Atoi(sizeStr)
	if err != nil {
		json.NewEncoder(w).Encode(Error{"bad params"})
		return
	}

	if size != 10 {
		json.NewEncoder(w).Encode(Error{"bad size"})
		return
	}

	historyArray := currency.GetHistory(size)
	json.NewEncoder(w).Encode(Currencies{historyArray})
}

func main() {
	http.HandleFunc("/auth", auth.AuthorizeHandler)
	http.HandleFunc("/register", register.RegisterHandler)

	http.HandleFunc("/ws", connectWebSocketHandler)
	http.HandleFunc("/get", getLastCurrency)

	/* TODO nginx */
	fs := http.FileServer(http.Dir("./exchangego/static"))
	http.Handle("/", fs)

	port := ":" + config.HTTPServer().Port
	log.Printf("Starting listening http server on port %q", port)
	http.ListenAndServe(port, nil)
}

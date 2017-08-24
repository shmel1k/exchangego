package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"errors"

	"github.com/shmel1k/exchangego/broadcast"
	"github.com/shmel1k/exchangego/config"
	"github.com/shmel1k/exchangego/context/errs"
	"github.com/shmel1k/exchangego/currency"
	"github.com/shmel1k/exchangego/exchange/auth"
	"github.com/shmel1k/exchangego/exchange/register"
	"github.com/shmel1k/exchangego/exchange/session/context"
	"github.com/shmel1k/exchangego/game"
)

var broadCaster *server.EasyCast

type Currencies struct {
	History []int `json:"history"`
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
	ctx, err := context.InitFromHTTP(w, r)
	if err != nil {
		errs.WriteError(w, err)
	}
	ok := broadCaster.Subscribe(ctx)
	if !ok {
		errs.WriteError(w, errs.Error{
			Status: http.StatusServiceUnavailable,
			Err:    "failed to subscribe",
		})
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
	ctx, err := context.InitFromHTTP(w, r)
	if err != nil {
		errs.WriteError(w, err)
		return
	}

	ctx.HTTPRequest().ParseForm()
	ctx.HTTPResponseWriter().Header().Set("Content-Type",
		"application/json")

	sizeStr, ok := simpleParam(ctx.HTTPRequest(), "size")
	if !ok {
		json.NewEncoder(w).Encode(Error{"bad params"})
		return
	}

	size, err := strconv.Atoi(sizeStr)
	if err != nil {
		ctx.WriteError(err)
		return
	}

	if size != 10 {
		ctx.WriteError(errors.New("need 10"))
		return
	}

	historyArray := currency.GetHistory(size)
	json.NewEncoder(ctx.HTTPResponseWriter()).Encode(Currencies{historyArray})
}

func main() {
	http.HandleFunc("/api/auth", auth.AuthorizeHandler)
	http.HandleFunc("/api/register", register.RegisterHandler)

	http.HandleFunc("/ws", connectWebSocketHandler)
	http.HandleFunc("/get", getLastCurrency)

	/* TODO nginx */
	fs := http.FileServer(http.Dir("./exchangego/static"))
	http.Handle("/", fs)

	port := ":" + config.HTTPServer().Port
	log.Printf("Starting listening http server on port %q", port)

	errs := make(chan error, 2)
	go func() {
		errs <- game.RunScheduler()
	}()

	go func() {
		errs <- http.ListenAndServe(port, nil)
	}()

	select {
	case t := <-errs:
		log.Fatal(t)
	}
}

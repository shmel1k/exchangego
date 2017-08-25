package main

import (
	"log"
	"net/http"
	"time"

	"github.com/shmel1k/exchangego/broadcast"
	"github.com/shmel1k/exchangego/config"
	"github.com/shmel1k/exchangego/context/errs"
	"github.com/shmel1k/exchangego/currency"
	"github.com/shmel1k/exchangego/exchange/auth"
	"github.com/shmel1k/exchangego/exchange/register"
	"github.com/shmel1k/exchangego/exchange/session/context"
	"github.com/shmel1k/exchangego/game"
	"github.com/shmel1k/exchangego/exchange/money"
)

var broadCaster *server.EasyCast

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

func main() {
	http.HandleFunc("/api/auth", auth.AuthorizeHandler)
	http.HandleFunc("/api/register", register.RegisterHandler)
	http.HandleFunc("/get", money.GetLastCurrency)

	http.HandleFunc("/ws", connectWebSocketHandler)

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

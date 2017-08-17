package auth

import (
	"context"
	"net/http"
	"time"

	cntxt "github.com/shmel1k/exchangego/context"
	"github.com/shmel1k/exchangego/context/errs"
	"github.com/shmel1k/exchangego/exchange"
)

const (
	maxQueryDuration = time.Second
)

func init() {

}

func Authorize(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), maxQueryDuration)
	defer cancel()

	var err error
	ctx, err = cntxt.InitFromHTTP(ctx, r)
	if err != nil {
		errs.WriteError(w, err)
		return
	}
	exchange.WriteOK(w, struct {
		Test string `json:"test"`
	}{
		Test: "ok",
	})
}

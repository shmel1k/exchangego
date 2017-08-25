package money

import (
	"net/http"
	"encoding/json"
	"strconv"
	"errors"
	"github.com/shmel1k/exchangego/context/errs"
	"github.com/shmel1k/exchangego/currency"
	"github.com/shmel1k/exchangego/exchange/session/context"
	"github.com/shmel1k/exchangego/exchange"
)

type Currencies struct {
	History []int `json:"history"`
}

func GetLastCurrency(w http.ResponseWriter, r *http.Request) {
	ctx, err := context.InitFromHTTP(w, r)
	if err != nil {
		errs.WriteError(w, err)
		return
	}

	ctx.HTTPResponseWriter().Header().Set("Content-Type",
		"application/json")

	var sizeStr string
	if !exchange.SimpleParam(ctx, "size", &sizeStr) {
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

package exchange

import (
	"encoding/json"
	"net/http"

	"github.com/shmel1k/exchangego/base/errs"
	"errors"
	"github.com/shmel1k/exchangego/exchange/session/context"
)

var (
	BadParamsError = errors.New("bad params")
	BadMethodError = errors.New("bad method")
)
type Response struct {
	Status int         `json:"status"`
	Body   interface{} `json:"body"`
}

type Error struct {
	Err string `json:"error"`
}

func SimpleParam(ctx *context.ExContext, key string, str *string) bool {
	value, ok := ctx.HTTPRequest().Form[key]
	if !ok || len(value) != 1 {
		ctx.WriteError(BadParamsError)
		return false
	}
	*str = value[0]
	return true
}

func IsOnlyMethod(ctx *context.ExContext, ar ...string) bool {
	success := false
	for _, method := range ar {
		if method == ctx.HTTPRequest().Method {
			success = true
			break
		}
	}
	if !success {
		ctx.WriteError(BadMethodError)
		return false
	}
	return true
}

func WriteOK(w http.ResponseWriter, data interface{}) {
	// FIXME(shmel1k): add easyjson or something like that
	r := Response{
		Status: http.StatusOK,
		Body:   data,
	}
	dt, err := json.Marshal(r)
	if err != nil {
		errs.WriteInternal(w)
		return
	}
	w.Write(dt)
}

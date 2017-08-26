package exgame

import (
	"net/http"
	"github.com/shmel1k/exchangego/exchange/session/context"
	"github.com/shmel1k/exchangego/exchange"
	"github.com/shmel1k/exchangego/database"
	"strconv"
	"github.com/shmel1k/exchangego/game"
	"github.com/shmel1k/exchangego/currency"
)

func getAuthContext(w http.ResponseWriter, r *http.Request) (*context.ExContext, error) {
	ctx, err := context.InitFromHTTP(w, r)
	if err != nil {
		return nil, err
	}

	err = ctx.InitUserFromCookie()
	if err != nil {
		return ctx, err
	}

	return ctx, err
}

func WelcomePage(w http.ResponseWriter, r *http.Request) {
	ctx, err := getAuthContext(w, r)
	if err != nil {
		if ctx != nil {
			http.Redirect(ctx.HTTPResponseWriter(), ctx.HTTPRequest(),
				"/auth", http.StatusMovedPermanently)
		}
		return
	}

	defer func() {
		ctx.Exit(recover())
	}()

	if !exchange.IsOnlyMethod(ctx, http.MethodGet) {
		return
	}

	exchange.ReturnTemplate(ctx, exchange.GameTmpl)
}

func StartGame(w http.ResponseWriter, r *http.Request) {
	ctx, err := getAuthContext(w, r)
	if err != nil {
		if ctx != nil {
			http.Redirect(ctx.HTTPResponseWriter(), ctx.HTTPRequest(),
				"/auth", http.StatusMovedPermanently)
		}
		return
	}

	defer func() {
		ctx.Exit(recover())
	}()

	if !exchange.IsOnlyMethod(ctx, http.MethodGet) {
		return
	}

	var typeMove, seconds string
	if !exchange.SimpleParam(ctx, "type", &typeMove) {
		return
	}
	if !exchange.SimpleParam(ctx, "seconds", &seconds) {
		return
	}

	/* TODO check noraml scope */
	move, ok := game.CastMoveType(typeMove)
	if !ok {
		ctx.WriteError(exchange.BadParamsError)
		return
	}

	duration, err := strconv.Atoi(seconds)
	if err != nil {
		ctx.WriteError(exchange.BadParamsError)
		return
	}

	id, err := database.AddUserTransaction(ctx.User(), string(move), duration)
	game.AddPlayer(ctx.User(), id, int64(duration), move, currency.GetCurrency())

	exchange.WriteOK(ctx.HTTPResponseWriter(), nil)
}
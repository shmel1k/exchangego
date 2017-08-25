package game

import (
	"net/http"
	"github.com/shmel1k/exchangego/exchange/session/context"
	"github.com/shmel1k/exchangego/context/errs"
	"github.com/shmel1k/exchangego/exchange"
	"github.com/shmel1k/exchangego/database"
	"strconv"
	"github.com/shmel1k/exchangego/game"
)

func startGame(w http.ResponseWriter, r *http.Request) {
	ctx, err := context.InitFromHTTP(w, r)
	if err != nil {
		errs.WriteError(w, err)
		return
	}
	defer func() {
		ctx.Exit(recover())
	}()

	err = ctx.InitUserFromCookie()
	if err != nil {
		errs.WriteError(w, err)
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

	transactionId, err := database.AddUserTransaction(ctx, move, duration)

	// start go
}
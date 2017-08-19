package register

import (
	"net/http"

	"github.com/shmel1k/exchangego/context/errs"
	"github.com/shmel1k/exchangego/database"
	"github.com/shmel1k/exchangego/exchange"
	"github.com/shmel1k/exchangego/exchange/session/context"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	ctx, err := context.InitFromHTTP(w, r)
	if err != nil {
		errs.WriteError(w, err)
		return
	}
	defer ctx.Exit(recover())

	user := r.URL.Query().Get("Login")
	password := r.URL.Query().Get("Password")

	resp, err := Register(ctx, RegisterRequest{
		Login:    user,
		Password: password,
	})

	switch {
	case err != nil:
		ctx.WriteError(err)
		return
	}
	exchange.WriteOK(ctx.HTTPResponseWriter(), resp)
}

type RegisterRequest struct {
	Password string
	Login    string
}

type RegisterResponse struct {
}

func Register(ctx *context.ExContext, param RegisterRequest) (RegisterResponse, error) {
	var err error
	if err = ctx.InitUser(); err != nil {
		return RegisterResponse{}, err
	}
	_, err = database.AddUser(ctx, param.Login, param.Password)
	if err != nil {
		if err == database.ErrUserExists {
			return RegisterResponse{}, errs.Error{
				Status: http.StatusForbidden,
				Err:    "user exists",
			}
		}
		return RegisterResponse{}, err
	}

	return RegisterResponse{}, nil
}

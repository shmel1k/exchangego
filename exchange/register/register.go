package register

import (
	"net/http"

	"github.com/shmel1k/exchangego/base/errs"
	"github.com/shmel1k/exchangego/database"
	"github.com/shmel1k/exchangego/exchange"
	"github.com/shmel1k/exchangego/exchange/session/context"
	"fmt"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	ctx, err := context.InitFromHTTP(w, r)
	if err != nil {
		errs.WriteError(w, err)
		return
	}
	defer ctx.Exit(recover())

	if !exchange.IsOnlyMethod(ctx, http.MethodGet, http.MethodPost) {
		return
	}

	if ctx.HTTPRequest().Method == http.MethodGet {
		exchange.ReturnTemplate(ctx, exchange.RegTmpl)
		return
	}

	var user, password string
	if !exchange.SimpleParam(ctx, "Login", &user) {
		return
	}

	if !exchange.SimpleParam(ctx, "Password", &password) {
		return
	}

	fmt.Println("user", user)
	fmt.Println(password)
	_, err = Register(ctx, RegisterRequest{
		Login:    user,
		Password: password,
	})

	switch {
	case err != nil:
		ctx.WriteError(err)
		return
	}

	http.Redirect(ctx.HTTPResponseWriter(), ctx.HTTPRequest(),
		"/auth", http.StatusMovedPermanently)
	// exchange.WriteOK(ctx.HTTPResponseWriter(), resp)
}

type RegisterRequest struct {
	Password string
	Login    string
}

type RegisterResponse struct {
}

func Register(ctx *context.ExContext, param RegisterRequest) (RegisterResponse, error) {
	var err error
	if err = ctx.InitUser(param.Login); err != nil && err != errs.ErrUserNotExists {
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

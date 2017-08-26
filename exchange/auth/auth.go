package auth

import (
	"net/http"

	"github.com/shmel1k/exchangego/base/errs"
	"github.com/shmel1k/exchangego/exchange"
	"github.com/shmel1k/exchangego/exchange/auth/cookie"
	"github.com/shmel1k/exchangego/exchange/session/context"
)

type AuthorizeRequest struct {
	Login    string
	Password string
}

type AuthorizeResponse struct {
}

func AuthorizeHandler(w http.ResponseWriter, r *http.Request) {
	ctx, err := context.InitFromHTTP(w, r)
	if err != nil {
		errs.WriteError(w, err)
		return
	}
	defer func() {
		ctx.Exit(recover())
	}()

	if !exchange.IsOnlyMethod(ctx, http.MethodGet, http.MethodPost) {
		return
	}

	if (r.Method == http.MethodGet) {
		exchange.ReturnTemplate(ctx, exchange.AuthTmpl)
		return
	}

	var user, password string
	if !exchange.SimpleParam(ctx, "Login", &user) {
		return
	}

	if !exchange.SimpleParam(ctx, "Password", &password) {
		return
	}

	_, err = Authorize(ctx, AuthorizeRequest{
		Login:    user,
		Password: password,
	})
	switch {
	case err != nil:
		ctx.WriteError(err)
		return
	}

	http.Redirect(ctx.HTTPResponseWriter(), ctx.HTTPRequest(), "/", http.StatusMovedPermanently)
	// exchange.WriteOK(ctx.HTTPResponseWriter(), resp)
}

func Authorize(ctx *context.ExContext, req AuthorizeRequest) (AuthorizeResponse, error) {
	if err := ctx.InitUser(req.Login); err != nil {
		return AuthorizeResponse{}, err
	}
	if req.Login == "" {
		return AuthorizeResponse{}, errs.Error{
			Status: http.StatusForbidden,
			Err:    "invalid login",
		}
	}
	if ctx.User().Password != req.Password {
		return AuthorizeResponse{}, errs.Error{
			Status: http.StatusForbidden,
			Err:    "forbidden",
		}
	}
	// FIXME(shmel1k): add adequate http headers here.
	cookieVal, err := cookie.GenerateCookie(ctx.User().Name)
	if err != nil {
		return AuthorizeResponse{}, err
	}

	http.SetCookie(ctx.HTTPResponseWriter(), &http.Cookie{
		Domain:   "",
		HttpOnly: true,
		Name:     cookie.CookieName,
		Value:    cookieVal,
	})

	return AuthorizeResponse{}, nil
}

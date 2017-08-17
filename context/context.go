package context

import (
	"context"
	"net/http"

	"github.com/shmel1k/exchangego/context/errs"
	"github.com/shmel1k/exchangego/database"
)

type Context struct {
	context.Context
	User database.User
}

func InitFromHTTP(ctx context.Context, r *http.Request) (*Context, error) {
	user := r.URL.Query().Get("Login")
	password := r.URL.Query().Get("Password")
	if password == "" {
		return nil, errs.Error{
			Status: http.StatusForbidden,
			Err:    "forbidden",
		}
	}

	u, err := database.FetchUser(ctx, user, password)
	_ = u
	return nil, err
}

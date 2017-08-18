package register

import (
	"log"
	"net/http"

	"github.com/shmel1k/exchangego/context"
	"github.com/shmel1k/exchangego/context/errs"
	"github.com/shmel1k/exchangego/database"
	"github.com/shmel1k/exchangego/exchange"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	ctx, err := context.InitFromHTTP(w, r)
	if err != nil {
		log.Println(err)
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

	// FIXME(a.petrukhin): add ctxlog!
	switch {
	case err != nil:
		log.Printf("[%s]: error %s", ctx.User().Name, err)
		errs.WriteError(ctx.HTTPResponseWriter(), err)
		return
	}
	exchange.WriteOK(ctx, resp)
}

type RegisterRequest struct {
	Password string
	Login    string
}

type RegisterResponse struct {
}

func Register(ctx *context.Context, param RegisterRequest) (RegisterResponse, error) {
	var err error
	u := ctx.User()
	if u.Name != "" {
		// XXX(a.petrukhin): add context logging.
		err = errs.Error{
			Status: http.StatusForbidden,
			Err:    "forbidden",
		}
		return RegisterResponse{}, err
	}

	_, err = database.AddUser(param.Login, param.Password)
	if err != nil {
		if err == database.ErrUserExists {
			return RegisterResponse{}, errs.Error{
				Status: http.StatusForbidden,
				Err:    "forbidden",
			}
		}
		return RegisterResponse{}, err
	}

	return RegisterResponse{}, nil
}

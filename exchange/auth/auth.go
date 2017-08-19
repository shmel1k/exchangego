package auth

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/shmel1k/exchangego/context/errs"
	"github.com/shmel1k/exchangego/exchange"
	"github.com/shmel1k/exchangego/exchange/session/context"
)

// XXX(shmel1k): move to config:
var key []byte = []byte(`HonestOption1234`)

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

	login := r.URL.Query().Get("Login")
	pass := r.URL.Query().Get("Password")
	resp, err := Authorize(ctx, AuthorizeRequest{
		Login:    login,
		Password: pass,
	})
	switch {
	case err != nil:
		ctx.WriteError(err)
		return
	}
	exchange.WriteOK(ctx.HTTPResponseWriter(), resp)
}

func Authorize(ctx *context.ExContext, req AuthorizeRequest) (AuthorizeResponse, error) {
	if err := ctx.InitUser(); err != nil {
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
	cookie, err := GenerateCookie(ctx)
	if err != nil {
		return AuthorizeResponse{}, err
	}

	http.SetCookie(ctx.HTTPResponseWriter(), &http.Cookie{
		Domain:   "",
		HttpOnly: true,
		Name:     "exchange",
		Value:    cookie,
	})

	return AuthorizeResponse{}, nil
}

func GenerateCookie(ctx *context.ExContext) (string, error) {
	buf := make([]byte, 0, 32)
	buf = append(buf, ctx.User().Name...)
	buf = append(buf, ':')
	buf = strconv.AppendInt(buf, time.Now().Unix(), 10)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to generate cookie[1]: %s", err)
	}
	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("failed to generate cookie[2]: %s", err)
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to generate cookie[3]: %s", err)
	}
	ciphertext := aesgcm.Seal(nil, nonce, buf, nil)
	return hex.EncodeToString(ciphertext), nil
}

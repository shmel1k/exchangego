package auth

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/shmel1k/exchangego/context"
	"github.com/shmel1k/exchangego/context/errs"
	"github.com/shmel1k/exchangego/exchange"
)

const (
	maxQueryDuration = time.Second
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
	resp, err := Authorize(ctx, AuthorizeRequest{
		Login:    r.URL.Query().Get("Login"),
		Password: r.URL.Query().Get("Password"),
	})

	switch {
	case err != nil:
		log.Printf("[%s]: error %s", ctx.User().Name, err)
		errs.WriteError(ctx.HTTPResponseWriter(), err)
		return
	}
	exchange.WriteOK(ctx, resp)
}

func Authorize(ctx *context.Context, req AuthorizeRequest) (AuthorizeResponse, error) {
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

func GenerateCookie(ctx *context.Context) (string, error) {
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

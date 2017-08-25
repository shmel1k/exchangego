package context

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"sync"

	"net"

	"log"

	cntxt "github.com/shmel1k/exchangego/context"
	"github.com/shmel1k/exchangego/context/contextlog"
	"github.com/shmel1k/exchangego/context/errs"
	"github.com/shmel1k/exchangego/database"
	"github.com/shmel1k/exchangego/exchange"
	"github.com/shmel1k/exchangego/exchange/auth/cookie"
)

type ExContext struct {
	cntxt.Context
	context context.Context

	scope RequestScope

	responseStatus  int
	responseMessage string
	prefix          string
}

type RequestScope struct {
	mu       sync.Mutex
	deferred []func()

	cn net.Conn

	request   *http.Request
	requestID string
	writer    http.ResponseWriter
	user      exchange.User
}

func withCancel(w http.ResponseWriter, r *http.Request) (*ExContext, context.CancelFunc) {
	ctx1, cancel := context.WithCancel(context.Background())
	ctx := &ExContext{
		context: ctx1,
		scope: RequestScope{
			request: r,
			writer:  w,
		},
	}
	return ctx, cancel
}

func InitFromHTTP(w http.ResponseWriter, r *http.Request) (*ExContext, error) {
	r.ParseForm()

	ctx, cancel := withCancel(w, r)
	ctx.Defer(cancel)
	ctx.Defer(ctx.writeAccessLog)
	defer func() {
		ctx.setLogPrefix()
	}()

	return ctx, nil
}

func (ctx *ExContext) InitUser() error {
	login := ctx.scope.request.URL.Query().Get("Login")
	if login == "" {
		return errs.ErrUserNotExists
	}
	err := ctx.fetchUser(login)
	if err != nil {
		return err
	}

	ctx.setLogPrefix()
	return nil
}

func (ctx *ExContext) InitUserFromCookie() error {
	c, err := ctx.HTTPRequest().Cookie(cookie.CookieName)
	if err != nil {
		return errs.Error{
			Status: http.StatusForbidden,
			Err:    "bad cookie",
		}
	}
	user, err := cookie.CheckCookie(c)
	if err != nil {
		if err == http.ErrNoCookie {
			return errs.Error{
				Status: http.StatusForbidden,
				Err:    "bad cookie",
			}
		}
		return err
	}

	err = ctx.fetchUser(user)
	if err != nil {
		return err
	}
	ctx.setLogPrefix()
	return nil
}

func (ctx *ExContext) PutCn(conn net.Conn) {
	ctx.scope.cn = conn
}

func (ctx *ExContext) Cn() net.Conn {
	if ctx.scope.cn == nil {
		log.Fatal("no connection")
	}
	return ctx.scope.cn
}

func (ctx *ExContext) fetchUser(user string) error {
	u, err := database.FetchUser(ctx, user)
	if err != nil {
		return err
	}
	ctx.scope.user = u
	if u.Name == "" {
		return errs.ErrUserNotExists
	}
	return nil
}

func (ctx *ExContext) Defer(f func()) {
	ctx.scope.mu.Lock()
	defer ctx.scope.mu.Unlock()

	ctx.scope.deferred = append(ctx.scope.deferred, f)
}

func (ctx *ExContext) Done() <-chan struct{} {
	return ctx.context.Done()
}

func (ctx *ExContext) Exit(panc interface{}) {
	// FIXME(shmel1k): add panic handling
	def := ctx.scope.deferred
	ctx.scope.mu.Lock()
	ctx.scope.deferred = nil
	ctx.scope.mu.Unlock()

	for _, v := range def {
		v()
	}
}

func (ctx *ExContext) writeAccessLog() {
	if ctx.responseStatus == 0 {
		ctx.responseStatus = http.StatusOK
	}
	// [prefix]: url=%q, method=%s, status=%d
	contextlog.Printf(ctx, "url=%q, method=%q, status=%d, msg=%s", ctx.HTTPRequest().URL.Path,
		ctx.HTTPRequest().Method, ctx.responseStatus, ctx.responseMessage)
}

func (ctx *ExContext) WriteError(err error) {
	st := errs.WriteError(ctx.HTTPResponseWriter(), err)
	ctx.responseMessage = fmt.Sprintf("%s", err)
	ctx.responseStatus = st
}

func (ctx *ExContext) HTTPResponseWriter() http.ResponseWriter {
	return ctx.scope.writer
}

func (ctx *ExContext) HTTPRequest() *http.Request {
	return ctx.scope.request
}

func (ctx *ExContext) User() exchange.User {
	return ctx.scope.user
}

func (ctx *ExContext) LogPrefix() string {
	return ctx.prefix
}

func (ctx *ExContext) RequestID() string {
	return ctx.scope.requestID
}

func (ctx *ExContext) setLogPrefix() {
	rid := ctx.HTTPRequest().Header.Get("Request-Id")
	if rid == "" {
		rid = newRequestID()
	}
	ctx.scope.requestID = rid
	ctx.prefix = fmt.Sprintf("[request_id=%s login=%s] ", rid, ctx.User().Name)
}

func newRequestID() string {
	return fmt.Sprintf("%08x%08x", rand.Uint32(), rand.Uint32())
}

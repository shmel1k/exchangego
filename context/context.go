package context

import (
	"context"
	"net/http"
	"sync"

	"github.com/shmel1k/exchangego/context/errs"
	"github.com/shmel1k/exchangego/database"
)

type Context struct {
	context context.Context
	user    database.User

	scope RequestScope
}

type RequestScope struct {
	mu       sync.Mutex
	deferred []func()

	request *http.Request
	writer  http.ResponseWriter
}

func withCancel(w http.ResponseWriter, r *http.Request) (*Context, context.CancelFunc) {
	ctx1, cancel := context.WithCancel(context.Background())
	ctx := &Context{
		context: ctx1,
		scope: RequestScope{
			request: r,
			writer:  w,
		},
	}
	return ctx, cancel
}

func InitFromHTTP(w http.ResponseWriter, r *http.Request) (*Context, error) {
	ctx, cancel := withCancel(w, r)
	ctx.Defer(cancel)

	user := r.URL.Query().Get("Login")
	password := r.URL.Query().Get("Password")
	if password == "" {
		return nil, errs.Error{
			Status: http.StatusForbidden,
			Err:    "forbidden",
		}
	}

	u, err := database.FetchUser(user, password)
	if err != nil {
		return nil, err
	}
	ctx.user = u
	return ctx, nil
}

func (ctx *Context) Defer(f func()) {
	ctx.scope.mu.Lock()
	defer ctx.scope.mu.Unlock()

	ctx.scope.deferred = append(ctx.scope.deferred, f)
}

func (ctx *Context) Done() <-chan struct{} {
	return ctx.context.Done()
}

func (ctx *Context) Exit(panc interface{}) {
	def := ctx.scope.deferred
	ctx.scope.mu.Lock()
	ctx.scope.deferred = nil
	ctx.scope.mu.Unlock()

	for _, v := range def {
		v()
	}
}

func (c *Context) HTTPResponseWriter() http.ResponseWriter {
	return c.scope.writer
}

func (c *Context) HTTPRequest() *http.Request {
	return c.scope.request
}

func (c *Context) User() database.User {
	return c.user
}

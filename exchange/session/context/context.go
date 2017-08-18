package context

import (
	"context"
	"net/http"
	"sync"

	cntxt "github.com/shmel1k/exchangego/context"
	"github.com/shmel1k/exchangego/database"
	"github.com/shmel1k/exchangego/exchange"
)

type ExContext struct {
	cntxt.Context
	context context.Context

	scope RequestScope

	prefix string
}

type RequestScope struct {
	mu       sync.Mutex
	deferred []func()

	request *http.Request
	writer  http.ResponseWriter
	user    exchange.User
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
	ctx, cancel := withCancel(w, r)
	ctx.Defer(cancel)

	return ctx, nil
}

func (ctx *ExContext) fetchUser(user string) error {
	u, err := database.FetchUser(ctx, user)
	if err != nil {
		return err
	}
	ctx.scope.user = u
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
	def := ctx.scope.deferred
	ctx.scope.mu.Lock()
	ctx.scope.deferred = nil
	ctx.scope.mu.Unlock()

	for _, v := range def {
		v()
	}
}

func (c *ExContext) HTTPResponseWriter() http.ResponseWriter {
	return c.scope.writer
}

func (c *ExContext) HTTPRequest() *http.Request {
	return c.scope.request
}

func (c *ExContext) User() exchange.User {
	return c.scope.user
}

func (c *ExContext) LogPrefix() string {
	return c.prefix
}

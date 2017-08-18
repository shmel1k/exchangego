package context

import (
	"context"
	"net/http"
)

type Context interface {
	context.Context
	Defer(f func())
	HTTPResponseWriter() http.ResponseWriter
	HTTPRequest() *http.Request
	LogPrefix() string
}

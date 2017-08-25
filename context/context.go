package context

import (
	"context"
	"net"
	"net/http"
)

type Context interface {
	context.Context
	Defer(f func())
	PutCn(cn net.Conn)
	Cn() net.Conn
	HTTPResponseWriter() http.ResponseWriter
	HTTPRequest() *http.Request
	LogPrefix() string
	RequestID() string
}

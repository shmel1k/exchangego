package context

import (
	"context"
	"net/http"
	"net"
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

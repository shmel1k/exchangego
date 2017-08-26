package base

import (
	"net/http"
	"net"
	"time"
)

type User struct {
	ID               uint32
	Name             string
	Password         string
	RegistrationDate time.Time

	Money int64
}

type Context interface {
	Defer(f func())
	PutCn(cn net.Conn)
	Cn() net.Conn
	HTTPResponseWriter() http.ResponseWriter
	HTTPRequest() *http.Request
	LogPrefix() string
	RequestID() string
}

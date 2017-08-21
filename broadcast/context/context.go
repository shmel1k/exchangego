package context

import (
	"context"
	"errors"
	"net"
)

var (
	cannotInitWebSocketError = errors.New("cannot init websocket")
)

type WsContext struct {
	ctx context.Context
	Cn  net.Conn

	context.CancelFunc

	Data interface{}
}

func newWsContext(cn net.Conn) *WsContext {
	ctx, cancel := context.WithCancel(context.Background())
	return &WsContext{
		ctx:        ctx,
		Cn:         cn,
		CancelFunc: cancel,
	}
}

func InitWebSocketContext(cn net.Conn) (*WsContext, error) {
	return newWsContext(cn), nil
}

func (ctx *WsContext) AttachData(foo interface{}) {
	ctx.Data = foo
}

func (ctx *WsContext) Exit() {
	ctx.CancelFunc()
}

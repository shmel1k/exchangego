package server

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/shmel1k/exchangego/exchange/session/context"
)

type PassPoolStruct struct {
	ctx      *context.ExContext
	mainCast *EasyCast

	msg string
}

type message struct {
	Message string `json:"value"`
	Time    int64  `json:"time"`
}

type EasyCast struct {
	ConnectionMap *CnMap
	pool          *Pool
}

func NewEasyCast(generator func() string, castDelay time.Duration, poolSize int) *EasyCast {
	easyCast := new(EasyCast)
	easyCast.pool = NewPool(poolSize)

	easyCast.ConnectionMap = NewConnectionStorage()

	go func(cast *EasyCast, delay time.Duration) {
		for {
			message := generator()

			lockMap := cast.ConnectionMap.GetAndLock()
			for ctx, _ := range lockMap {
				cast.pool.ThrowTask(shareAllUsers, &PassPoolStruct{
					ctx:      ctx,
					mainCast: cast,
					msg:      message,
				})
			}
			cast.ConnectionMap.UnLock()

			time.Sleep(delay)
		}
	}(easyCast, castDelay)

	return easyCast
}

func shareAllUsers(msg_ interface{}) {
	passPool, _ := msg_.(*PassPoolStruct)
	ctx := passPool.ctx

	now := time.Now().Unix()

	resp, _ := json.Marshal(message{passPool.msg, now})
	err := wsutil.WriteServerMessage(ctx.Cn(), ws.OpText, resp)
	if err != nil {
		/* close connection */
		log.Println("Close connection")

		ctx.Exit(recover())
		passPool.mainCast.ConnectionMap.TryRemove(ctx)
	}
}

func tryOpenConnection(w http.ResponseWriter, r *http.Request) (net.Conn, bool) {
	conn, _, _, err := ws.UpgradeHTTP(r, w, nil)
	if err != nil {
		return nil, false
	}

	return conn, true
}

func (ec *EasyCast) Subscribe(ctx *context.ExContext) bool {
	cn, ok := tryOpenConnection(ctx.HTTPResponseWriter(), ctx.HTTPRequest())
	if !ok {
		return false
	}
	ctx.PutCn(cn)

	ec.ConnectionMap.Put(ctx)
	return true
}

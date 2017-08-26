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
	"fmt"
)

type PassPoolStruct struct {
	userName string

	cn net.Conn
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
			for name, cn := range lockMap {
				cast.pool.ThrowTask(shareAllUsers, &PassPoolStruct{
					cn:      	cn,
					userName:	name,
					mainCast: 	cast,
					msg:      	message,
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
	cn := passPool.cn

	now := time.Now().Unix()

	resp, _ := json.Marshal(message{passPool.msg, now})
	err := wsutil.WriteServerMessage(cn, ws.OpText, resp)

	if err != nil {
		/* close connection */
		log.Println("Close connection")
		passPool.mainCast.ConnectionMap.TryRemove(passPool.userName)
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

	fmt.Println("Add to ", ctx.User().Name)
	ec.ConnectionMap.Put(ctx.User().Name, cn)
	return true
}

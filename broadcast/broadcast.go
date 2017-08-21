package server

import (
	"encoding/json"
	"log"
	"math/rand"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/shmel1k/exchangego/broadcast/context"
)

type PassPoolStruct struct {
	ctx *context.WsContext
	msg string
}

type message struct {
	Message string `json:"value"`
	Time    int64  `json:"time"`
}

type EasyCast struct {
	ConnectionMap *CnMap

	pool *Pool
}

func getEasyCastCtx(wsContext *context.WsContext) *EasyCast {
	if wsContext.Data == nil {
		log.Fatal("cannot get data")
	}

	return wsContext.Data.(*EasyCast)
}

func NewEasyCast(castDelay time.Duration, poolSize int) *EasyCast {
	easyCast := new(EasyCast)
	easyCast.pool = NewPool(poolSize)

	easyCast.ConnectionMap = NewConnectionStorage()

	go func(cast *EasyCast, delay time.Duration) {
		var currency = 50
		for {
			/* base func */
			currency += rand.Intn(10) - 5
			currencyString := strconv.Itoa(currency)

			lockMap := cast.ConnectionMap.GetAndLock()
			for ctx, _ := range lockMap {
				cast.pool.ThrowTask(shareAllUsers, &PassPoolStruct{
					ctx: ctx,
					msg: currencyString,
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
	err := wsutil.WriteServerMessage(ctx.Cn, ws.OpText, resp)
	if err != nil {
		/* close connection */
		log.Println("Close connection")
		getEasyCastCtx(ctx).ConnectionMap.TryRemove(ctx)
	}
}

func tryOpenConnection(w http.ResponseWriter, r *http.Request) (net.Conn, bool) {
	conn, _, _, err := ws.UpgradeHTTP(r, w, nil)
	if err != nil {
		return nil, false
	}

	return conn, true
}

func (ec *EasyCast) Subscribe(w http.ResponseWriter, r *http.Request) bool {
	cn, ok := tryOpenConnection(w, r)
	if !ok {
		return false
	}

	ctx, err := context.InitWebSocketContext(cn)
	if err != nil {
		log.Println(err)
		return false
	}

	ctx.AttachData(ec)

	ec.ConnectionMap.Put(ctx)
	return true
}

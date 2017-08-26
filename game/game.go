package game

import (
	"errors"
	"sync"
	"time"
	"github.com/shmel1k/exchangego/base"
	"github.com/shmel1k/exchangego/currency"
	"github.com/shmel1k/exchangego/database"
	"fmt"
	"github.com/shmel1k/exchangego/broadcast"
	"encoding/json"
	"github.com/gobwas/ws/wsutil"
	"github.com/gobwas/ws"
)

var localCast *server.EasyCast
func InitGame(cast *server.EasyCast) {
	localCast = cast
}

type updateInfo struct {
	Status	string 		`json:"status"`
	Money	int64	 	`json:"money"`
}

type MoveType string
const (
	UpMoveType 	 	MoveType 	=	"up"
	DownMoveType 	MoveType 	= 	"down"
	UnknownMoveType MoveType	=	""
)

func CastMoveType(param string) (MoveType, bool) {
	if param == "up" {
		return UpMoveType, true
	} else if param == "down" {
		return DownMoveType, true
	}

	return UnknownMoveType, false
}

type TransactionResult int
const (
	InWaitResult TransactionResult = 0
	FinishResult TransactionResult = 1
)

type game struct {
	transactionID int64
	duration      int64
	end           int64

	move MoveType // False -- down, True -- up
	startmoney int
}

const (
	updateTime = 1 * time.Second
)

var (
	ErrUserExists = errors.New("failed to add user to exgame: user exists")
)

var players Players

type Players struct {
	players map[base.User]game

	mu sync.Mutex
}

func (p *Players) Add(user base.User, transactionID int64, duration int64, move MoveType, startmoney int) error {
	if p.players != nil {
		if _, ok := p.players[user]; ok {
			return ErrUserExists
		}
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	if p.players == nil {
		p.players = make(map[base.User]game)
	}

	p.players[user] = game{
		transactionID: transactionID,
		duration:      duration,
		end:           time.Now().Unix() + duration,
		move:          move,

		startmoney: startmoney,
	}

	return nil
}

func (p *Players) Get(user base.User) game {
	if p.players == nil {
		return game{}
	}
	return p.players[user]
}

func (p *Players) Delete(user base.User) {
	p.mu.Lock()
	delete(p.players, user)
	p.mu.Unlock()
}

func AddPlayer(user base.User, transactionID int64, duration int64, move MoveType, startmoney int) error {
	return players.Add(user, transactionID, duration, move, startmoney)
}

func RunScheduler() error {
	for {
		err := schedule()
		if err != nil {
			return err
		}
	}
}

func schedule() error {
	playersToUpdate := make([]base.User, 0, len(players.players))
	for {
		curr := currency.GetCurrency()

		t := time.Now().Unix()
		for k, v := range players.players {
			if v.end <= t {
				playersToUpdate = append(playersToUpdate, k)
			}
		}
		var err error
		for _, v := range playersToUpdate {
			fmt.Println("test", v)

			g := players.Get(v)
			mon := v.Money

			var status string
			if g.startmoney >= curr {
				status = "win"
				mon = mon + 100
			} else {
				status = "lost"
				mon = mon - 100
			}

			err = database.UpdateMoney(v.ID, mon)
			if err != nil {
				return err
			}
			players.Delete(v)

			fmt.Println("try to user", v.Name)
			mapCn := localCast.ConnectionMap.GetAndLock()
			fmt.Println(mapCn)

			connection := mapCn[v.Name]
			resp, _ := json.Marshal(updateInfo{
				Status: status,
				Money: mon,
			})

			err := wsutil.WriteServerMessage(connection, ws.OpText, resp)
			if err != nil {
				fmt.Println("error!", err)
			}

			localCast.ConnectionMap.UnLock()
		}
		playersToUpdate = playersToUpdate[:0]

		time.Sleep(updateTime)
	}
}

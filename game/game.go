package game

import (
	"errors"
	"sync"
	"time"

	"github.com/shmel1k/exchangego/base"
	"github.com/shmel1k/exchangego/database"
)

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
	duration int64
	end      int64

	move MoveType // False -- down, True -- up
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

func (p *Players) Add(user base.User, duration int64, move MoveType) error {
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
		duration: duration,
		end:      time.Now().Unix() + duration,
		move:     move,
	}

	return nil
}

func (p *Players) Delete(user base.User) {
	p.mu.Lock()
	delete(p.players, user)
	p.mu.Unlock()
}

func AddPlayer(user base.User, duration int64, move MoveType) error {
	return players.Add(user, duration, move)
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
		t := time.Now().Unix()
		for k, v := range players.players {
			if v.end <= t {
				playersToUpdate = append(playersToUpdate, k)
			}
		}
		var err error
		for _, v := range playersToUpdate {
			err = database.UpdateMoney(v.ID, v.Money/2)
			if err != nil {
				return err
			}
		}
		playersToUpdate = playersToUpdate[:0]

		time.Sleep(updateTime)
	}
}

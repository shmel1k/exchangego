package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/shmel1k/exchangego/config"
	"github.com/shmel1k/exchangego/context"
	"github.com/shmel1k/exchangego/context/contextlog"
	"github.com/shmel1k/exchangego/exchange"
)

const (
	dbName = "exchange"

	defaultMoney = 100
)

var ErrUserExists = errors.New("User already exists")

var db *sql.DB
var once sync.Once

func initClient() {
	once.Do(func() {
		cfg := config.Database()
		conf := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", cfg.User, cfg.Password, cfg.Address, cfg.Port, dbName)

		var err error
		db, err = sql.Open("mysql", conf)
		if err != nil {
			log.Fatalf("failed to init mysql database: %s", err)
		}

		db.SetMaxOpenConns(cfg.MaxOpenConns)

		err = db.Ping()
		if err != nil {
			log.Fatalf("failed to init mysql database: failed to ping: %s", err)
		}
	})
}

func FetchUser(ctx context.Context, user string) (exchange.User, error) {
	initClient()
	// XXX(shmel1k): fix registration_date in scanning
	q := fmt.Sprintf("SELECT id, name, password FROM users WHERE name = ?")
	resp, err := db.Query(q, user)
	if err != nil {
		return exchange.User{}, fmt.Errorf("failed to perform query %q: %s", q, err)
	}

	var u exchange.User
	for resp.Next() {
		err = resp.Scan(&u.ID, &u.Name, &u.Password)
		if err != nil {
			return exchange.User{}, fmt.Errorf("failed to scan response from mysql: %s", err)
		}
	}

	q = fmt.Sprintf("SELECT score FROM money WHERE user_id = ?")
	resp, err = db.Query(q, u.ID)
	if err != nil {
		return exchange.User{}, fmt.Errorf("failed to fetch money from mysql: %s", err)
	}

	for resp.Next() {
		err = resp.Scan(&u.Money)
		if err != nil {
			return exchange.User{}, fmt.Errorf("failed to scan money-response from mysql: %s", err)
		}
	}

	return u, nil
}

func AddUser(ctx context.Context, user, password string) (exchange.User, error) {
	initClient()
	// FIXME(shmel1k): add transaction
	q := fmt.Sprint("SELECT COUNT(*) FROM users WHERE name = ?")
	resp, err := db.Query(q, user)
	if err != nil {
		return exchange.User{}, fmt.Errorf("failed to perform query %q: %s", q, err)
	}

	var cnt uint32
	for resp.Next() {
		err = resp.Scan(&cnt)
		if err != nil {
			return exchange.User{}, fmt.Errorf("failed to scan response from mysql: %s", err)
		}
		if cnt != 0 {
			return exchange.User{}, ErrUserExists
		}
	}

	q = fmt.Sprint("INSERT INTO users(name, password, registration_date) VALUES(?, ?, ?)")
	t := time.Now()

	_, err = db.Query(q, user, password, t)
	if err != nil {
		return exchange.User{}, fmt.Errorf("failed to perform query %q: %s", q, err)
	}
	contextlog.Printf(ctx, "User %s successfully added", user)

	us := exchange.User{
		ID:               0,
		Name:             user,
		Password:         password,
		RegistrationDate: t,
	}

	res, err := db.Query("SELECT id FROM users WHERE name = ?", user)
	if err != nil {
		contextlog.Printf(ctx, "failed to get user_id for user %q: %s", user, err)
		return us, nil
	}

	var userID uint32
	for res.Next() {
		err = res.Scan(&userID)
		if err != nil {
			contextlog.Printf(ctx, "failed to scan user_id for user %q: %s", user, err)
		}
	}

	us.ID = userID

	_, err = db.Query("INSERT INTO money VALUES(?, ?)", userID, defaultMoney)
	if err != nil {
		contextlog.Printf(ctx, "failed to insert money for user %q: %s", user, err)
	}

	// FIXME(shmel1k): add UserID
	return us, nil
}

func UpdateMoney(userid uint32, money int64) error {
	// NOTE: this function in used in scheduler.
	q := fmt.Sprint("UPDATE money SET money = ? WHERE user_id = ?")
	_, err := db.Query(q, money, userid)
	if err != nil {
		return fmt.Errorf("failed to update money for user %d: %s", userid, err)
	}
	return nil
}

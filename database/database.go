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
	"github.com/shmel1k/exchangego/exchange"
)

const (
	dbName = "exchange"
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

	return u, nil
}

func AddUser(ctx context.Context, user, password string) (exchange.User, error) {
	initClient()

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

	// FIXME(shmel1k): add UserID
	return exchange.User{
		ID:               0,
		Name:             user,
		Password:         password,
		RegistrationDate: t,
	}, nil
}

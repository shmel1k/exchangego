package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/shmel1k/exchangego/config"
)

const (
	dbName = "exchange"
)

var db *sql.DB

func init() {
	cfg := config.Database()
	conf := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", cfg.User, cfg.Password, cfg.Address, cfg.Port, dbName)

	var err error
	db, err = sql.Open("mysql", conf)
	if err != nil {
		log.Fatalf("failed to init mysql database: %s", err)
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)
}

type User struct {
	Name             string
	Password         string
	RegistrationDate time.Time
}

func FetchUser(ctx context.Context, user, password string) (User, error) {
	// FIXME(shmel1k): add context here.
	return User{}, nil
}

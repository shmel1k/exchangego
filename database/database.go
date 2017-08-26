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
	"github.com/shmel1k/exchangego/base"
	"github.com/shmel1k/exchangego/base/contextlog"
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

func FetchUser(ctx base.Context, user string) (base.User, error) {
	initClient()
	// XXX(shmel1k): fix registration_date in scanning
	q := fmt.Sprintf("SELECT id, name, password FROM users WHERE name = ?")
	resp, err := db.Query(q, user)
	if err != nil {
		return base.User{}, fmt.Errorf("failed to perform query %q: %s", q, err)
	}

	var u base.User
	for resp.Next() {
		err = resp.Scan(&u.ID, &u.Name, &u.Password)
		if err != nil {
			return base.User{}, fmt.Errorf("failed to scan response from mysql: %s", err)
		}
	}

	q = fmt.Sprintf("SELECT score FROM money WHERE user_id = ?")
	resp, err = db.Query(q, u.ID)
	if err != nil {
		return base.User{}, fmt.Errorf("failed to fetch money from mysql: %s", err)
	}

	for resp.Next() {
		err = resp.Scan(&u.Money)
		if err != nil {
			return base.User{}, fmt.Errorf("failed to scan money-response from mysql: %s", err)
		}
	}

	return u, nil
}

func AddUser(ctx base.Context, user, password string) (base.User, error) {
	initClient()
	// FIXME(shmel1k): add transaction
	q := fmt.Sprint("SELECT COUNT(*) FROM users WHERE name = ?")
	resp, err := db.Query(q, user)
	if err != nil {
		return base.User{}, fmt.Errorf("failed to perform query %q: %s", q, err)
	}

	var cnt uint32
	for resp.Next() {
		err = resp.Scan(&cnt)
		if err != nil {
			return base.User{}, fmt.Errorf("failed to scan response from mysql: %s", err)
		}
		if cnt != 0 {
			return base.User{}, ErrUserExists
		}
	}

	q = fmt.Sprint("INSERT INTO users(name, password, registration_date) VALUES(?, ?, ?)")
	t := time.Now()

	_, err = db.Query(q, user, password, t)
	if err != nil {
		return base.User{}, fmt.Errorf("failed to perform query %q: %s", q, err)
	}
	contextlog.Printf(ctx, "User %s successfully added", user)

	us := base.User{
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

	return us, nil
}

func AddUserTransaction(user base.User, moveType string, duration int) (int64, error) {
	initClient()

	sql_ := "INSERT INTO transactions (user_id, type, ts, duration_s, result) VALUES (?, ?, ?, ?, ?);"
	resp, err := db.Exec(sql_, user.ID, moveType, time.Now().Unix(), duration, 0 /* Start Wait */)
	if err != nil {
		log.Printf("failed to insert transaction %q: %s", user, err)
		return 0, err
	}

	id, err := resp.LastInsertId()

	/* TODO money */

	return id, err
}

func ChangeStatusTransaction(transactionId int, status int) error {
	initClient()

	sql_ := "UPDATE transactions SET result = ? WHERE id = ?"
	_, err := db.Query(sql_, status, transactionId)
	if err != nil {
		fmt.Errorf("failed to update transaction id %d: %s", transactionId, err)
		return err
	}

	return nil
}

func UpdateMoney(userid uint32, money int64) error {
	// NOTE: this function in used in scheduler.
	q := fmt.Sprint("UPDATE money SET score = ? WHERE user_id = ?")
	_, err := db.Query(q, money, userid)
	if err != nil {
		return fmt.Errorf("failed to update money for user %d: %s", userid, err)
	}
	return nil
}

package kzbotdb

import (
	"context"
	"fmt"
	"os"
	"queue/config"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/mentalisit/logger"

	_ "github.com/lib/pq"
)

type Db struct {
	db     *sqlx.DB
	log    *logger.Logger
	client *sqlx.DB
}

func NewDb(log *logger.Logger, cfg *config.ConfigBot) *Db {
	db, err := NewClient(context.Background(), log, 5, cfg)
	if err != nil {
		log.Fatal(err.Error())
	}
	d := &Db{
		db:     db,
		log:    log,
		client: db,
	}
	return d
}

func NewClient(ctx context.Context, log *logger.Logger, maxAttempts int, conf *config.ConfigBot) (*sqlx.DB, error) {
	var db *sqlx.DB
	var err error
	dns := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		conf.Postgress.Username, conf.Postgress.Password, conf.Postgress.Host, conf.Postgress.Name)

	err = DoWithTries(func() error {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		db, err = sqlx.ConnectContext(ctx, "postgres", dns)
		if err != nil {
			dns = fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
				conf.Postgress.Username, conf.Postgress.Password, "192.168.100.131:5435", conf.Postgress.Name)

			db, err = sqlx.ConnectContext(ctx, "postgres", dns)
			if err != nil {
				log.ErrorErr(err)
				os.Exit(1)
			}
		}
		return nil
	}, maxAttempts, 5*time.Second)
	if err != nil {
		log.Fatal(err.Error())
	}

	return db, nil
}
func DoWithTries(fn func() error, attemtps int, delay time.Duration) (err error) {
	for attemtps > 0 {
		if err = fn(); err != nil {
			time.Sleep(delay)
			attemtps--

			continue
		}
		return nil
	}
	return
}

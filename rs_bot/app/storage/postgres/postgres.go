package postgres

import (
	"context"
	"fmt"
	"os"
	"rs/config"
	"rs/pkg/utils"
	"time"

	_ "github.com/lib/pq"

	"github.com/jmoiron/sqlx"
	"github.com/mentalisit/logger"
)

type Db struct {
	db  *sqlx.DB
	log *logger.Logger
}

func NewDb(log *logger.Logger, cfg *config.ConfigBot) *Db {
	db, err := NewClient(log, 5, cfg)
	if err != nil {
		log.Fatal(err.Error())
	}
	d := &Db{
		db:  db,
		log: log,
	}

	return d
}

func (d *Db) Shutdown() {
	d.db.Close()
}

func NewClient(log *logger.Logger, maxAttempts int, conf *config.ConfigBot) (*sqlx.DB, error) {
	var db *sqlx.DB
	var err error
	dns := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		conf.Postgress.Username, conf.Postgress.Password, conf.Postgress.Host, conf.Postgress.Name)

	err = utils.DoWithTries(func() error {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
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

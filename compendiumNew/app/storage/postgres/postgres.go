package postgres

import (
	"compendium/config"
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/mentalisit/logger"

	_ "github.com/lib/pq"
)

type Db struct {
	db  *sqlx.DB
	log *logger.Logger
}

func NewDb(log *logger.Logger, cfg *config.ConfigBot) *Db {
	dns := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		cfg.Postgress.Username, cfg.Postgress.Password, cfg.Postgress.Host, cfg.Postgress.Name)
	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFunc()

	db, err := sqlx.ConnectContext(ctx, "postgres", dns)
	if err != nil {
		log.ErrorErr(err)
		os.Exit(1)
	}
	database := &Db{
		db:  db,
		log: log,
	}

	return database
}

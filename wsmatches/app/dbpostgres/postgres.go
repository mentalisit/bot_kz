package dbpostgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/mentalisit/logger"
	"os"
)

type Db struct {
	log  *logger.Logger
	pool *pgxpool.Pool
}

func NewDb(log *logger.Logger) *Db {
	dns := fmt.Sprintf("postgres://postgres:root@%s/postgres", "postgres:5432")
	ctx := context.Background()

	pool, err := pgxpool.Connect(ctx, dns)
	if err != nil {
		dns = fmt.Sprintf("postgres://postgres:root@%s/postgres", "192.168.100.131:5435")
		pool, err = pgxpool.Connect(ctx, dns)
		if err != nil {
			log.ErrorErr(err)
			os.Exit(1)
		}
	}
	d := &Db{
		log:  log,
		pool: pool,
	}

	return d
}

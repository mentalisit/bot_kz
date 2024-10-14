package dbpostgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/mentalisit/logger"
	"os"
	"time"
)

type Db struct {
	log  *logger.Logger
	pool *pgxpool.Pool
}

func NewDb(log *logger.Logger) *Db {
	dns := fmt.Sprintf("postgres://postgres:root@%s/postgres", "postgres:5432")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

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
func (d *Db) GetContext() (ctx context.Context, cancel context.CancelFunc) {
	return context.WithTimeout(context.Background(), 10*time.Second)
}

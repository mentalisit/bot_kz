package postgresLocal

import (
	"context"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/mentalisit/logger"
	"kz_bot/config"
	"kz_bot/pkg/utils"
	"os"
	"time"
)

type Client interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Begin(ctx context.Context) (pgx.Tx, error)
}

func NewClient(ctx context.Context, log *logger.Logger, maxAttempts int, conf *config.ConfigBot) (*pgxpool.Pool, error) {
	var pool *pgxpool.Pool
	var err error
	dns := fmt.Sprintf("postgres://%s:%s@%s/%s",
		conf.Postgress.Username, conf.Postgress.Password, conf.Postgress.Host, conf.Postgress.Name)

	err = utils.DoWithTries(func() error {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		pool, err = pgxpool.Connect(ctx, dns)
		if err != nil {
			dns = fmt.Sprintf("postgres://%s:%s@%s/%s",
				conf.Postgress.Username, conf.Postgress.Password, "192.168.100.155:5435", conf.Postgress.Name)

			pool, err = pgxpool.Connect(ctx, dns)
			if err != nil {
				log.ErrorErr(err)
				os.Exit(1)
			}

			//return err
		}
		return nil
	}, maxAttempts, 5*time.Second)
	if err != nil {
		log.Fatal(err.Error())
	}

	return pool, nil
}

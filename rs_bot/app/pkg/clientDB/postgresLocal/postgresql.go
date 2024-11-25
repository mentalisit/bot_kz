package postgresLocal

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mentalisit/logger"
	"os"
	"rs/config"
	"rs/pkg/utils"
	"time"
)

type Client interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Begin(ctx context.Context) (pgx.Tx, error)
}

func NewClient(log *logger.Logger, maxAttempts int, conf *config.ConfigBot) (*pgxpool.Pool, error) {
	var pool *pgxpool.Pool
	var err error
	dns := fmt.Sprintf("postgres://%s:%s@%s/%s",
		conf.Postgress.Username, conf.Postgress.Password, conf.Postgress.Host, conf.Postgress.Name)

	err = utils.DoWithTries(func() error {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		pool, err = pgxpool.New(ctx, dns)
		if err != nil {
			dns = fmt.Sprintf("postgres://%s:%s@%s/%s",
				conf.Postgress.Username, conf.Postgress.Password, "192.168.100.131:5435", conf.Postgress.Name)

			pool, err = pgxpool.New(ctx, dns)
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

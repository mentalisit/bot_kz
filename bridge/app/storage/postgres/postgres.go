package postgres

import (
	"bridge/config"
	"context"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/mentalisit/logger"
	"os"
	"time"
)

type Db struct {
	db  Client
	log *logger.Logger
}
type Client interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Begin(ctx context.Context) (pgx.Tx, error)
}

func NewDb(log *logger.Logger, cfg *config.ConfigBot) *Db {
	dns := fmt.Sprintf("postgres://%s:%s@%s/%s",
		cfg.Postgress.Username, cfg.Postgress.Password, cfg.Postgress.Host, cfg.Postgress.Name)

	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFunc()
	pool, err := pgxpool.Connect(ctx, dns)
	if err != nil {
		log.ErrorErr(err)
		os.Exit(1)
		//return err
	}
	if err != nil {
		log.Fatal(err.Error())
	}
	db := &Db{
		db:  pool,
		log: log,
	}
	go db.createTable()
	return db
}
func (d *Db) createTable() {
	ctx, cancel := d.GetContext()
	defer cancel()
	_, err := d.db.Exec(ctx,
		`CREATE TABLE IF NOT EXISTS kzbot.bridge_config (
        id SERIAL PRIMARY KEY,
		name_relay TEXT,
		host_relay TEXT,
		role TEXT[],
		channel_ds JSONB,
		channel_tg JSONB,
		forbidden_prefixes TEXT[]
        );
    `)
	if err != nil {
		d.log.ErrorErr(err)
	}
}
func (d *Db) GetContext() (ctx context.Context, cancel context.CancelFunc) {
	return context.WithTimeout(context.Background(), 10*time.Second)
}

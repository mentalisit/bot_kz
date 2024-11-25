package postgres

import (
	"context"
	"discord/config"
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

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	pool, err := pgxpool.Connect(ctx, dns)
	if err != nil {
		log.ErrorErr(err)
		time.Sleep(5 * time.Second)
		os.Exit(1)
		//return err
	}
	db := &Db{
		db:  pool,
		log: log,
	}

	go db.createTable()

	return db
}
func (d *Db) getContext() (ctx context.Context, cancel context.CancelFunc) {
	return context.WithTimeout(context.Background(), 5*time.Second)
}
func (d *Db) createTable() {
	ctx, cancel := d.getContext()
	defer cancel()
	d.db.Exec(ctx, "CREATE SCHEMA IF NOT EXISTS kzbot")
	// Создание таблиц
	_, err := d.db.Exec(ctx,
		`CREATE TABLE IF NOT EXISTS kzbot.event (
            id             BIGSERIAL PRIMARY KEY,
            dateStart      TEXT,
            dateStop       TEXT,
            message        TEXT
        );
    `)
	if err != nil {
		d.log.ErrorErr(err)
		return
	}
}

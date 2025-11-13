package postgres

import (
	"context"
	"fmt"
	"os"
	"telegram/config"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mentalisit/logger"
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
	pool, err := pgxpool.New(ctx, dns)
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
	db.createTable()
	return db
}
func (d *Db) getContext() (ctx context.Context, cancel context.CancelFunc) {
	return context.WithTimeout(context.Background(), 10*time.Second)
}
func (d *Db) createTable() {
	ctx, cancel := d.getContext()
	defer cancel()
	// Создание таблиц
	_, err := d.db.Exec(ctx,
		`CREATE TABLE IF NOT EXISTS rs_bot.chat_members (
		chat_id BIGINT PRIMARY KEY,
		chat_name text NOT NULL DEFAULT '',
		data JSONB NOT NULL DEFAULT '{}',
		roles JSONB NOT NULL DEFAULT '{}',
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`)
	if err != nil {
		d.log.ErrorErr(err)
		return
	}
}

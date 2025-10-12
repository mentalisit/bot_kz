package postgres

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"
	"whatsapp/config"

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
		slog.Error(err.Error())
		//log.ErrorErr(err)
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
		slog.Error(err.Error())
		//d.log.ErrorErr(err)
		return
	}

	_, err = d.db.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS rs_bot.battlestop
	(
		id     bigserial        primary key,
		corporation text NOT NULL DEFAULT '',
		name text NOT NULL DEFAULT '',
		level    integer NOT NULL DEFAULT 0,
		count   integer NOT NULL DEFAULT 0
	);`)
	if err != nil {
		slog.Error(err.Error())
		//d.log.ErrorErr(err)
		return
	}

	_, err = d.db.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS rs_bot.scoreboard
	(
		id     bigserial        primary key,
		Name text NOT NULL DEFAULT '',
		WebhookChannel    text NOT NULL DEFAULT '',
		ScoreChannel   text NOT NULL DEFAULT '',
		LastMessage text NOT NULL DEFAULT ''
	);`)
	if err != nil {
		slog.Error(err.Error())
		//d.log.ErrorErr(err)
		return
	}

	_, err = d.db.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS rs_bot.webhooks
	(
		id     bigserial        primary key,
		TsUnix bigint NOT NULL DEFAULT 0,
		corp    text NOT NULL DEFAULT '',
		message   jsonb
	);`)
	if err != nil {
		slog.Error(err.Error())
		//d.log.ErrorErr(err)
		return
	}

	//_, err = d.db.Exec(ctx, `
	//	CREATE TABLE rs_bot.name_aliases (
	//	alias TEXT PRIMARY KEY,
	//	canonical_name TEXT NOT NULL
	//);
	//	CREATE INDEX idx_alias ON rs_bot.name_aliases(alias);`)
	//if err != nil {
	//	slog.Error(err.Error())
	//	//d.log.ErrorErr(err)
	//	return
	//}
}

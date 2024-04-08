package postgres

import (
	"compendium/config"
	"context"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/mentalisit/logger"
	"os"
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

	pool, err := pgxpool.Connect(context.Background(), dns)
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
	d.db.Exec(context.Background(), "CREATE SCHEMA IF NOT EXISTS compendium")
	// Создание таблицы compendium.corpmember
	_, err := d.db.Exec(context.Background(), `CREATE TABLE IF NOT EXISTS compendium.corpmember (
        guildid      TEXT,
        name         TEXT,
        userid       TEXT,
        clientuserid TEXT,
        avatar       TEXT,
        tech         JSONB,
        avatarurl    TEXT,
        timezona     TEXT,
        zonaoffset   NUMERIC,
        afkfor       TEXT
    )`)
	if err != nil {
		fmt.Println("Ошибка при создании таблицы compendium.corpmember:", err)
		return
	}

	// Создание таблицы compendium.guild
	_, err = d.db.Exec(context.Background(), `CREATE TABLE IF NOT EXISTS compendium.guild (
        token TEXT,
        url   TEXT,
        id    TEXT,
        name  TEXT,
        icon  TEXT
    )`)
	if err != nil {
		fmt.Println("Ошибка при создании таблицы compendium.guild:", err)
		return
	}

	// Создание таблицы compendium.identity
	_, err = d.db.Exec(context.Background(), `CREATE TABLE IF NOT EXISTS compendium.identity (
        id TEXT PRIMARY KEY
    )`)
	if err != nil {
		fmt.Println("Ошибка при создании таблицы compendium.identity:", err)
		return
	}

	// Создание таблицы compendium.user
	_, err = d.db.Exec(context.Background(), `CREATE TABLE IF NOT EXISTS compendium."user" (
        token         TEXT,
        id            TEXT,
        username      TEXT,
        discriminator TEXT,
        avatar        TEXT,
        avatarurl     TEXT
    )`)
	if err != nil {
		fmt.Println("Ошибка при создании таблицы compendium.user:", err)
		return
	}
}

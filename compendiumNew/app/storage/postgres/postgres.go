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
	d.db.Exec(context.Background(), "CREATE SCHEMA IF NOT EXISTS hs_compendium")

	// Создание таблицы users
	_, err := d.db.Exec(context.Background(), `CREATE TABLE IF NOT EXISTS hs_compendium.users (
        id            bigserial primary key,
        userid        TEXT,
        username      TEXT,
        discriminator TEXT,
        avatar        TEXT,
        avatarurl     TEXT,
        alts text[]
    )`)
	if err != nil {
		fmt.Println("Ошибка при создании таблицы users:", err)
		return
	}

	// Создание таблицы guilds
	_, err = d.db.Exec(context.Background(), `CREATE TABLE IF NOT EXISTS hs_compendium.guilds (
	   id bigserial primary key,
	   url   TEXT,
	   guildid    TEXT,
	   name  TEXT,
	   icon  TEXT
	)`)
	if err != nil {
		fmt.Println("Ошибка при создании таблицы guilds:", err)
		return
	}

	// Создание таблицы list_users
	_, err = d.db.Exec(context.Background(), `CREATE TABLE IF NOT EXISTS hs_compendium.list_users (
	   id bigserial primary key,
	   token   TEXT,
	   userid    TEXT,
	   guildid  TEXT
	)`)
	if err != nil {
		fmt.Println("Ошибка при создании таблицы guilds:", err)
		return
	}

	// Создание таблицы corpmember
	_, err = d.db.Exec(context.Background(), `CREATE TABLE IF NOT EXISTS hs_compendium.corpmember (
	id           bigserial primary key,
	username     TEXT,
	userid       TEXT,
	guildid 	 TEXT,
	avatar       TEXT,
	avatarurl    TEXT,
	timezona     TEXT,
	zonaoffset   NUMERIC,
	afkfor       TEXT
	)`)
	if err != nil {
		fmt.Println("Ошибка при создании таблицы corpmember:", err)
		return
	}

	_, err = d.db.Exec(context.Background(), `CREATE TABLE IF NOT EXISTS hs_compendium.tech (
    id bigserial primary key,
    username text,
    userid text,
    guildid text,
    tech jsonb
    )`)
	if err != nil {
		fmt.Println("Ошибка при создании таблицы tech:", err)
		return
	}

	// Создание таблицы userroles
	_, err = d.db.Exec(context.Background(), `CREATE TABLE IF NOT EXISTS hs_compendium.userroles (
	   id           bigserial primary key,
	   guildid      TEXT,
	   role         TEXT,
	   username     TEXT,
	   userid       TEXT
	)`)
	if err != nil {
		fmt.Println("Ошибка при создании таблицы userroles:", err)
		return
	}

	// Создание таблицы guildroles
	_, err = d.db.Exec(context.Background(), `CREATE TABLE IF NOT EXISTS hs_compendium.guildroles (
	   id           bigserial primary key,
	   guildid      TEXT,
	   role         TEXT
	)`)
	if err != nil {
		fmt.Println("Ошибка при создании таблицы guildroles:", err)
		return
	}

	// Создание таблицы wskill
	_, err = d.db.Exec(context.Background(), `CREATE TABLE IF NOT EXISTS hs_compendium.wskill (
	id           bigserial primary key,
	guildid 	 TEXT,
	chatid 	     TEXT,
	username     TEXT,
	mention      TEXT,
	shipname     TEXT,
	timestampend BIGSERIAL
	)`)
	if err != nil {
		fmt.Println("Ошибка при создании таблицы wskill:", err)
		return
	}
}

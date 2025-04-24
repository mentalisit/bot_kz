package multi

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mentalisit/logger"
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

func NewDb(log *logger.Logger, cl *pgxpool.Pool) *Db {
	db := &Db{
		db:  cl,
		log: log,
	}
	go db.createTable()
	return db
}
func (d *Db) getContext() (ctx context.Context, cancel context.CancelFunc) {
	return context.WithTimeout(context.Background(), 10*time.Second)
}

func (d *Db) createTable() {
	ctx, cancel := d.getContext()
	defer cancel()
	d.db.Exec(ctx, "CREATE SCHEMA IF NOT EXISTS compendium")

	_, err := d.db.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS compendium.multi_accounts
	(
		uuid  uuid primary key DEFAULT gen_random_uuid(),
		nickname text NOT NULL DEFAULT '',
		telegram_id text unique,
		telegram_username  TEXT DEFAULT '',
		discord_id text unique,
		discord_username  TEXT DEFAULT '',
		whatsapp_id  text unique,
		whatsapp_username  TEXT DEFAULT '',
		avatarUrl     TEXT NOT NULL DEFAULT '',
        alts 	TEXT[] NOT NULL DEFAULT ARRAY[]::TEXT[],
		created_at timestamp default now()
	);`)
	if err != nil {
		d.log.ErrorErr(err)
		return
	}

	_, err = d.db.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS compendium.accounts_link_codes
	(
		code  text primary key,
		uuid uuid references compendium.multi_accounts(uuid) on delete cascade,
		expires_at timestamp not null,
		created_at timestamp default now()
	);`)
	if err != nil {
		d.log.ErrorErr(err)
		return
	}

	_, err = d.db.Exec(ctx, `CREATE TABLE IF NOT EXISTS compendium.technologies (
    uid uuid references compendium.multi_accounts(uuid) on delete cascade,
    username text,
    tech jsonb
    )`)
	if err != nil {
		d.log.ErrorErr(err)
		return
	}

	// Создание таблицы guilds
	_, err = d.db.Exec(ctx, `CREATE TABLE IF NOT EXISTS compendium.guilds (
	   gid  uuid primary key DEFAULT gen_random_uuid(),
	   GuildName   TEXT NOT NULL DEFAULT '',
	   Channels  TEXT[] NOT NULL DEFAULT ARRAY[]::TEXT[],
	   AvatarUrl   TEXT NOT NULL DEFAULT ''
	)`)
	if err != nil {
		fmt.Println("Ошибка при создании таблицы guilds:", err)
		return
	}

	// Создание таблицы corpmember
	_, err = d.db.Exec(ctx, `CREATE TABLE IF NOT EXISTS compendium.corpMember (
	   uid uuid REFERENCES compendium.multi_accounts(uuid) ON DELETE CASCADE,
	   guildIds UUID[] NOT NULL DEFAULT ARRAY[]::UUID[],
	   timeZona TEXT NOT NULL DEFAULT '',
	   zonaOffset INTEGER NOT NULL DEFAULT 0,
	   afkFor TEXT NOT NULL DEFAULT ''
)`)
	if err != nil {
		d.log.ErrorErr(err)
		return
	}

	//// Создание таблицы list_users
	//_, err = d.db.Exec(ctx, `CREATE TABLE IF NOT EXISTS compendium.list_users (
	//   uid uuid references compendium.multi_accounts(uuid) on delete cascade,
	//   guildId 	 TEXT,
	//   token   TEXT primary key
	//
	//)`)
	//if err != nil {
	//	d.log.ErrorErr(err)
	//	return
	//}
	//
	//// Создание таблицы userroles
	//_, err = d.db.Exec(ctx, `CREATE TABLE IF NOT EXISTS compendium.userRoles (
	//   id           bigserial primary key,
	//   guildId      TEXT,
	//   role         TEXT,
	//   username     TEXT,
	//   uid uuid references compendium.multi_accounts(uuid) on delete cascade
	//)`)
	//if err != nil {
	//	d.log.ErrorErr(err)
	//	return
	//}
	//
	//// Создание таблицы guildroles
	//_, err = d.db.Exec(ctx, `CREATE TABLE IF NOT EXISTS compendium.guildRoles (
	//   id           bigserial primary key,
	//   guildId      TEXT,
	//   role         TEXT
	//)`)
	//if err != nil {
	//	d.log.ErrorErr(err)
	//	return
	//}
	//
	//// Создание таблицы wskill
	//_, err = d.db.Exec(ctx, `CREATE TABLE IF NOT EXISTS compendium.wsKill (
	//id           bigserial primary key,
	//guildId 	 TEXT,
	//chatId 	     TEXT,
	//username     TEXT,
	//mention      TEXT,
	//shipName     TEXT,
	//timestampEnd BIGSERIAL
	//)`)
	//if err != nil {
	//	d.log.ErrorErr(err)
	//	return
	//}

}

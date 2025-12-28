package postgresv2

import (
	"compendium/config"
	"fmt"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/mentalisit/logger"
)

type Db struct {
	db  *sqlx.DB
	log *logger.Logger
}

func NewDb(log *logger.Logger, cfg *config.ConfigBot) *Db {
	dns := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		cfg.Postgress.Username, cfg.Postgress.Password, cfg.Postgress.Host, cfg.Postgress.Name)

	db, err := sqlx.Open("postgres", dns)
	if err != nil {
		log.ErrorErr(err)
		os.Exit(1)
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(1 * time.Minute)

	d := &Db{
		db:  db,
		log: log,
	}

	go d.createTable()

	return d
}

func (d *Db) createTable() {
	var query []string
	query = append(query, "CREATE SCHEMA IF NOT EXISTS my_compendium")

	query = append(query, `CREATE TABLE IF NOT EXISTS my_compendium.multi_accounts (
			uuid  uuid primary key DEFAULT gen_random_uuid(),
			nickname text NOT NULL DEFAULT '',
			telegram_id text NOT NULL DEFAULT '',
			telegram_username  TEXT DEFAULT '',
			discord_id TEXT DEFAULT '',
			discord_username  TEXT DEFAULT '',
			whatsapp_id  TEXT DEFAULT '',
			whatsapp_username  TEXT DEFAULT '',
			avatarUrl     TEXT NOT NULL DEFAULT '',
			alts 	TEXT[] NOT NULL DEFAULT ARRAY[]::TEXT[],
			created_at timestamp default now()
		)`)
	query = append(query, `CREATE TABLE IF NOT EXISTS my_compendium.technologies (
		uid uuid references my_compendium.multi_accounts(uuid) on delete cascade,
		username text,
		tech jsonb
		)`)
	query = append(query, `CREATE TABLE IF NOT EXISTS my_compendium.guilds (
		gid  uuid primary key DEFAULT gen_random_uuid(),
		GuildName   TEXT NOT NULL DEFAULT '',
		Channels  JSONB NOT NULL DEFAULT '{}',
		AvatarUrl   TEXT NOT NULL DEFAULT ''
		)`)
	query = append(query, `CREATE TABLE IF NOT EXISTS my_compendium.corpMember (
		uid uuid REFERENCES my_compendium.multi_accounts(uuid) ON DELETE CASCADE,
		guildIds UUID[] NOT NULL DEFAULT ARRAY[]::UUID[],
		timeZona TEXT NOT NULL DEFAULT '',
		zonaOffset INTEGER NOT NULL DEFAULT 0,
		afkFor TEXT NOT NULL DEFAULT ''
	)`)
	query = append(query, `CREATE TABLE IF NOT EXISTS my_compendium.codes (
    	id           bigserial primary key,
		code    	 TEXT,
		time 	     bigint,
		identity     jsonb
	)`)
	query = append(query, `CREATE TABLE IF NOT EXISTS my_compendium.wskill (
    	id           bigserial primary key,
		guildid 	 TEXT,
		chatid 	     TEXT,
		username     TEXT,
		mention      TEXT,
		shipname     TEXT,
		timestampend BIGSERIAL,
		language     TEXT                                                
	)`)

	for _, s := range query {
		if _, err := d.db.Exec(s); err != nil {
			d.log.ErrorErr(fmt.Errorf("failed to create table: %w", err))
		}
	}
}

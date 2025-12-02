package postgresv2

import (
	"compendium_s/config"
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
	// Create schema
	if _, err := d.db.Exec("CREATE SCHEMA IF NOT EXISTS my_compendium"); err != nil {
		d.log.ErrorErr(fmt.Errorf("failed to create schema: %w", err))
		return
	}

	var corpMemberTable string
	// Create multi_accounts table
	multiAccountsTable := `
		CREATE TABLE IF NOT EXISTS my_compendium.multi_accounts
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
		)`
	if _, err := d.db.Exec(multiAccountsTable); err != nil {
		d.log.ErrorErr(fmt.Errorf("failed to create multi_accounts table: %w", err))
		return
	}

	// Create accounts_link_codes table
	linkCodesTable := `
		CREATE TABLE IF NOT EXISTS my_compendium.accounts_link_codes
		(
			code  text primary key,
			uuid uuid references my_compendium.multi_accounts(uuid) on delete cascade,
			expires_at timestamp not null,
			created_at timestamp default now()
		)`
	if _, err := d.db.Exec(linkCodesTable); err != nil {
		d.log.ErrorErr(fmt.Errorf("failed to create accounts_link_codes table: %w", err))
		return
	}

	// Create technologies table
	technologiesTable := `CREATE TABLE IF NOT EXISTS my_compendium.technologies (
		uid uuid references my_compendium.multi_accounts(uuid) on delete cascade,
		username text,
		tech jsonb
		)`
	if _, err := d.db.Exec(technologiesTable); err != nil {
		d.log.ErrorErr(fmt.Errorf("failed to create technologies table: %w", err))
		return
	}

	// Create guilds table
	guildsTable := `CREATE TABLE IF NOT EXISTS my_compendium.guilds (
		gid  uuid primary key DEFAULT gen_random_uuid(),
		GuildName   TEXT NOT NULL DEFAULT '',
		Channels  JSONB NOT NULL DEFAULT '{}',
		AvatarUrl   TEXT NOT NULL DEFAULT ''
		)`
	if _, err := d.db.Exec(guildsTable); err != nil {
		d.log.ErrorErr(fmt.Errorf("failed to create guilds table: %w", err))
		return
	}

	// Create corpMember table
	corpMemberTable = `CREATE TABLE IF NOT EXISTS my_compendium.corpMember (
		uid uuid REFERENCES my_compendium.multi_accounts(uuid) ON DELETE CASCADE,
		guildIds UUID[] NOT NULL DEFAULT ARRAY[]::UUID[],
		timeZona TEXT NOT NULL DEFAULT '',
		zonaOffset INTEGER NOT NULL DEFAULT 0,
		afkFor TEXT NOT NULL DEFAULT ''
	)`
	if _, err := d.db.Exec(corpMemberTable); err != nil {
		d.log.ErrorErr(fmt.Errorf("failed to create corpMember table: %w", err))
		return
	}
}

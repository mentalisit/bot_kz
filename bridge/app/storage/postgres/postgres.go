package postgres

import (
	"bridge/config"
	"database/sql"
	"fmt"
	"os"

	"github.com/mentalisit/logger"
)

type Db struct {
	db  *sql.DB
	log *logger.Logger
}

func NewDb(log *logger.Logger, cfg *config.ConfigBot) *Db {
	dns := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		cfg.Postgress.Username, cfg.Postgress.Password, cfg.Postgress.Host, cfg.Postgress.Name)

	// Открытие соединения
	conn, err := sql.Open("postgres", dns)
	if err != nil {
		log.ErrorErr(err)
		os.Exit(1)
	}

	// Проверка подключения
	if err = conn.Ping(); err != nil {
		log.ErrorErr(err)
		os.Exit(1)
	}

	db := &Db{
		db:  conn,
		log: log,
	}

	go db.createTable()
	return db
}

func (d *Db) createTable() {
	_, err := d.db.Exec(
		`CREATE TABLE IF NOT EXISTS rs_bot.bridge_config (
        id SERIAL PRIMARY KEY,
		name_relay TEXT,
		host_relay TEXT,
		role TEXT[],
		channel JSONB,
		forbidden_prefixes TEXT[]
        );
    `)
	if err != nil {
		d.log.ErrorErr(err)
	}

	_, err = d.db.Exec(
		`CREATE TABLE IF NOT EXISTS rs_bot.message_maps (
    	id SERIAL PRIMARY KEY,
    	message_ids JSONB NOT NULL,
    	created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
    	);
		CREATE INDEX IF NOT EXISTS idx_message_ids_gin ON rs_bot.message_maps USING GIN (message_ids);
	`)
	if err != nil {
		d.log.ErrorErr(err)
	}
}

package postgres

import (
	"bridge/config"
	"database/sql"
	"fmt"
	"github.com/mentalisit/logger"
	"os"
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

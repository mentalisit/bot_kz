package postgres

import (
	"context"
	"discord/config"
	"fmt"
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/mentalisit/logger"
	"github.com/mentalisit/restapi/models"

	_ "github.com/lib/pq"
)

type Db struct {
	db  *sqlx.DB
	log *logger.Logger
	sync.RWMutex
	RsBotConfig  map[string]models.CorporationConfigV2
	BridgeConfig map[string]models.Bridge2Config
	KzBotConfig  map[string]models.CorporationConfig
	pool         *sqlx.DB
	dns          string
}

func NewDb(log *logger.Logger, cfg *config.ConfigBot) *Db {
	dns := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		cfg.Postgress.Username, cfg.Postgress.Password, cfg.Postgress.Host, cfg.Postgress.Name)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	db, err := sqlx.ConnectContext(ctx, "postgres", dns)
	if err != nil {
		slog.Error(err.Error())
		time.Sleep(5 * time.Second)
		os.Exit(1)
	}
	database := &Db{
		db:           db,
		log:          log,
		RsBotConfig:  make(map[string]models.CorporationConfigV2),
		BridgeConfig: make(map[string]models.Bridge2Config),
		KzBotConfig:  make(map[string]models.CorporationConfig),
		pool:         db,
	}
	database.dns = dns

	go database.createTable()

	database.loadConfig()
	go database.StartConfigWatcher(context.Background())

	return database
}

func (d *Db) createTable() {
	fmt.Printf("lastEventId %d\n", d.readLastEvent())

	// Создание таблиц
	_, err := d.db.Exec(
		`CREATE TABLE IF NOT EXISTS rs_bot2.event_schedule (
            id             BIGSERIAL PRIMARY KEY,
            dateStart      TEXT,
            dateStop       TEXT,
            season        integer
        );
    `)
	if err != nil {
		slog.Error(err.Error())
		//d.log.ErrorErr(err)
		return
	}

	_, err = d.db.Exec(`
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

}

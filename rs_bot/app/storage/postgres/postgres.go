package postgres

import (
	"context"
	"rs/config"
	"rs/pkg/clientDB/postgresLocal"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mentalisit/logger"
)

type Db struct {
	db     postgresLocal.Client
	log    *logger.Logger
	client *pgxpool.Pool
	DB     postgresLocal.Client
}

func NewDb(log *logger.Logger, cfg *config.ConfigBot) *Db {
	db, err := postgresLocal.NewClient(log, 5, cfg)
	if err != nil {
		log.Fatal(err.Error())
	}
	d := &Db{
		db:     db,
		log:    log,
		client: db,
		DB:     db,
	}
	go d.createTable()
	return d
}
func (d *Db) getContext() (ctx context.Context, cancel context.CancelFunc) {
	return context.WithTimeout(context.Background(), 10*time.Second)
}
func (d *Db) createTable() {
	ctx, cancel := d.getContext()
	defer cancel()
	d.db.Exec(ctx, "CREATE SCHEMA IF NOT EXISTS kzbot")
	// Создание таблиц
	_, err := d.db.Exec(ctx,
		`CREATE TABLE IF NOT EXISTS kzbot.config (
            id             BIGSERIAL PRIMARY KEY,
            corpname       TEXT,
            dschannel      TEXT,
            tgchannel      TEXT,
            mesiddshelp    TEXT,
            mesidtghelp    TEXT,
            country        TEXT,
            delmescomplite BIGINT,
            guildid        TEXT,
            forvard        boolean
        );
    `)
	if err != nil {
		d.log.ErrorErr(err)
		return
	}

	_, err = d.db.Exec(ctx,
		`CREATE TABLE IF NOT EXISTS kzbot.sborkz(
		id          bigserial        primary key,
		corpname    text,
		name        text,
		mention     text,
		tip         text,
		dsmesid     text,    
		tgmesid     bigint,
		wamesid     text,    
		time        text,
		date        text,    
		lvlkz       text,
		numkzn      bigint,    
		numberkz    bigint,
		numberevent bigint,
		eventpoints bigint,    
		active      bigint,
		timedown    bigint,
		userid      text);`)
	if err != nil {
		d.log.ErrorErr(err)
		return
	}

	_, err = d.db.Exec(ctx,
		`CREATE TABLE IF NOT EXISTS kzbot.numkz(
    	id       bigserial	primary key,
        lvlkz    text,
        number   bigint,
        corpname text
	);`)
	if err != nil {
		d.log.ErrorErr(err)
		return
	}

	_, err = d.db.Exec(ctx,
		`CREATE TABLE IF NOT EXISTS kzbot.rsevent(
    id          bigserial        primary key,
	corpname    text,    
	numevent    bigint,
	activeevent bigint,    
	number      bigint
	);`)
	if err != nil {
		d.log.ErrorErr(err)
		return
	}

	_, err = d.db.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS kzbot.rsevent(
		id          bigserial        primary key,
		corpname    text,    numevent    bigint,
		activeevent bigint,    number      bigint
	);`)
	if err != nil {
		d.log.ErrorErr(err)
		return
	}

	_, err = d.db.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS kzbot.event
	(
		id     bigserial        primary key,
		datestart   text,
		datestop  text,
		message text
	);`)
	if err != nil {
		d.log.ErrorErr(err)
		return
	}

}
func (d *Db) Shutdown() {
	d.client.Close()
}

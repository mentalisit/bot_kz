package postgres

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mentalisit/logger"
	"rs/config"
	"rs/pkg/clientDB/postgresLocal"
	"time"
)

type Db struct {
	db     postgresLocal.Client
	log    *logger.Logger
	client *pgxpool.Pool
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
	}
	go d.createTable()
	return d
}
func (d *Db) GetContext() (ctx context.Context, cancel context.CancelFunc) {
	return context.WithTimeout(context.Background(), 10*time.Second)
}
func (d *Db) createTable() {
	ctx, cancel := d.GetContext()
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
	_, err = d.db.Exec(ctx,
		`CREATE TABLE IF NOT EXISTS kzbot.subscribe(
    id        bigserial	primary key,
    name      text,
    nameid    text,
    lvlkz     text,
    tip       bigint,
    chatid    text,
    timestart text,
    timeend   text
	);`)
	if err != nil {
		d.log.ErrorErr(err)
		return
	}
	_, err = d.db.Exec(ctx,
		`CREATE TABLE IF NOT EXISTS kzbot.timer(
    id       bigserial primary key,
    dsmesid  text,
    dschatid text,
    tgmesid  text,
    tgchatid text,
    timed    bigint
	);`)
	if err != nil {
		d.log.ErrorErr(err)
		return
	}
	_, err = d.db.Exec(ctx,
		`CREATE TABLE IF NOT EXISTS kzbot.users(
    id      bigserial primary key,
    tip     text,
	name    text,
	em1     text,
	em2     text,
	em3     text,
	em4     text,
	module1 text,
	module2 text,
	module3 text,
	weapon  text
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

	_, _ = d.db.Exec(ctx, "CREATE SCHEMA IF NOT EXISTS rs_bot")

	_, err = d.db.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS rs_bot.user_account
	(
		internal_user_id     bigserial        primary key,
		general_name text NOT NULL DEFAULT '',
		user_id_tg   text NOT NULL DEFAULT '',
		user_id_ds  text NOT NULL DEFAULT '',
		user_id_game text[] DEFAULT '{}',
		user_name_active text NOT NULL DEFAULT '',
		user_accounts text[] DEFAULT '{}'
	);`)
	if err != nil {
		d.log.ErrorErr(err)
		return
	}

	_, err = d.db.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS rs_bot.eventnumber
	(
		id     bigserial        primary key,
		number   integer NOT NULL DEFAULT 0,
		event  integer NOT NULL DEFAULT 0,
		status bool NOT NULL DEFAULT false
	);`)
	if err != nil {
		d.log.ErrorErr(err)
		return
	}
}
func (d *Db) Shutdown() {
	d.client.Close()
}

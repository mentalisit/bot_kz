package multi

import (
	"compendium_s/models"
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mentalisit/logger"
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

}

func (d *Db) DeleteOldClient(uid uuid.UUID) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	del := "delete from compendium.corpmember where uid = $1"
	_, _ = d.db.Exec(ctx, del, uid)
	del = "delete from compendium.multi_accounts where uuid = $1"
	_, _ = d.db.Exec(ctx, del, uid)
	del = "delete from compendium.technologies where uid = $1"
	_, _ = d.db.Exec(ctx, del, uid)
}
func (d *Db) SearchOldData(i models.Identity) (m models.Moving) {
	account, _ := d.FindMultiAccountUUID(i.MAccount.UUID)
	if account != nil {
		m.MAcc = *account
	}
	corpMember, _ := d.CorpMemberByUId(i.MAccount.UUID)
	if corpMember != nil {
		m.CorpMember = *corpMember
	}
	m.Tech = d.TechnologiesGetMember(i.MAccount.UUID)
	return m
}

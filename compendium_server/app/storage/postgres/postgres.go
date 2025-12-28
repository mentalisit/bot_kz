package postgres

import (
	"compendium_s/config"
	"compendium_s/models"
	"compendium_s/storage/postgres/multi"
	"context"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mentalisit/logger"
)

type Db struct {
	db    Client
	log   *logger.Logger
	Multi *multi.Db
}
type Client interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Begin(ctx context.Context) (pgx.Tx, error)
}

func NewDb(log *logger.Logger, cfg *config.ConfigBot) *Db {
	dns := fmt.Sprintf("postgres://%s:%s@%s/%s",
		cfg.Postgress.Username, cfg.Postgress.Password, cfg.Postgress.Host, cfg.Postgress.Name)
	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFunc()
	pool, err := pgxpool.New(ctx, dns)
	if err != nil {
		log.ErrorErr(err)
		os.Exit(1)
		//return err
	}

	db := &Db{
		db:    pool,
		log:   log,
		Multi: multi.NewDb(log, pool),
	}
	go db.createTable()
	return db
}
func (d *Db) createTable() {
	ctx, cancel := d.getContext()
	defer cancel()
	d.db.Exec(ctx, "CREATE SCHEMA IF NOT EXISTS hs_compendium")

	// Создание таблицы users
	_, err := d.db.Exec(ctx, `CREATE TABLE IF NOT EXISTS hs_compendium.users (
        id            bigserial primary key,
        userid        TEXT,
        username      TEXT,
        discriminator TEXT,
        avatar        TEXT,
        avatarurl     TEXT,
        alts text[],
        gamename      TEXT
    )`)
	if err != nil {
		fmt.Println("Ошибка при создании таблицы users:", err)
		return
	}

	//// Создание таблицы guilds
	//_, err = d.db.Exec(ctx, `CREATE TABLE IF NOT EXISTS hs_compendium.guilds (
	//   id bigserial primary key,
	//   url   TEXT,
	//   guildid    TEXT,
	//   name  TEXT,
	//   icon  TEXT
	//)`)
	//if err != nil {
	//	fmt.Println("Ошибка при создании таблицы guilds:", err)
	//	return
	//}

	// Создание таблицы list_users
	_, err = d.db.Exec(ctx, `CREATE TABLE IF NOT EXISTS hs_compendium.list_users (
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
	_, err = d.db.Exec(ctx, `CREATE TABLE IF NOT EXISTS hs_compendium.corpmember (
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

	_, err = d.db.Exec(ctx, `CREATE TABLE IF NOT EXISTS hs_compendium.tech (
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
	_, err = d.db.Exec(ctx, `CREATE TABLE IF NOT EXISTS hs_compendium.userroles (
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
	_, err = d.db.Exec(ctx, `CREATE TABLE IF NOT EXISTS hs_compendium.guildroles (
	   id           bigserial primary key,
	   guildid      TEXT,
	   role         TEXT
	)`)
	if err != nil {
		fmt.Println("Ошибка при создании таблицы guildroles:", err)
		return
	}

	// Создание таблицы wskill
	_, err = d.db.Exec(ctx, `CREATE TABLE IF NOT EXISTS hs_compendium.wskill (
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

	// Создание таблицы code
	_, err = d.db.Exec(ctx, `CREATE TABLE IF NOT EXISTS hs_compendium.codes (
	id           bigserial primary key,
	code    	 TEXT,
	identity     jsonb,
	timestamp 	 bigint
                                               
	)`)
	if err != nil {
		fmt.Println("Ошибка при создании таблицы codes:", err)
		return
	}
}

func (d *Db) getContext() (ctx context.Context, cancel context.CancelFunc) {
	return context.WithTimeout(context.Background(), 10*time.Second)
}

func (d *Db) DeleteOldClient(userid string) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	del := "delete from hs_compendium.corpmember where userid = $1"
	_, _ = d.db.Exec(ctx, del, userid)
	del = "delete from hs_compendium.list_users where userid = $1"
	_, _ = d.db.Exec(ctx, del, userid)
	del = "delete from hs_compendium.tech where userid = $1"
	_, _ = d.db.Exec(ctx, del, userid)
	del = "delete from hs_compendium.userroles where userid = $1"
	_, _ = d.db.Exec(ctx, del, userid)
	del = "delete from hs_compendium.users where userid = $1"
	_, _ = d.db.Exec(ctx, del, userid)
}

func (d *Db) SearchOldData(i models.Identity) (m models.Moving) {
	m.MAcc = *i.MAccount
	m.CorpMember.Uid = i.MAccount.UUID
	var cmm []models.CorpMember
	var guildMap map[string]struct{}
	guildMap = make(map[string]struct{})
	if m.MAcc.DiscordID != "" {
		cm, _ := d.CorpMemberRead(m.MAcc.DiscordID)
		cmm = append(cmm, cm...)
	}
	if m.MAcc.TelegramID != "" {
		cm, _ := d.CorpMemberRead(m.MAcc.TelegramID)
		cmm = append(cmm, cm...)
	}
	if m.MAcc.WhatsappID != "" {
		cm, _ := d.CorpMemberRead(m.MAcc.WhatsappID)
		cmm = append(cmm, cm...)
	}
	var techMap map[string]models.TechLevels
	techMap = make(map[string]models.TechLevels)
	for _, member := range cmm {
		if techMap[member.Name] == nil {
			techMap[member.Name] = member.Tech
		} else {
			for module, data := range member.Tech {
				if techMap[member.Name][module].Level == 0 || techMap[member.Name][module].Ts < data.Ts {
					techMap[member.Name][module] = data
				}
			}
		}
		if m.CorpMember.AfkFor == "" {
			m.CorpMember.AfkFor = member.AfkFor
		}
		if m.CorpMember.TimeZona == "" {
			m.CorpMember.TimeZona = member.TimeZone
		}
		if m.CorpMember.ZonaOffset == 0 {
			m.CorpMember.ZonaOffset = member.ZoneOffset
		}
		guildMap[member.GuildId] = struct{}{}
	}
	for name, tech := range techMap {
		m.Tech = append(m.Tech, models.Technology{
			Tech: tech,
			Name: name,
		})
	}

	for g, _ := range guildMap {
		guid, err := uuid.Parse(g)
		if err == nil {
			m.CorpMember.GuildIds = append(m.CorpMember.GuildIds, guid)
		}
	}

	return m
}

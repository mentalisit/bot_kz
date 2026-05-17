package postgres

import (
	"compendium_s/config"
	"compendium_s/models"
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/mentalisit/logger"

	_ "github.com/lib/pq"
)

type Db struct {
	//db  Client
	db  *sqlx.DB
	log *logger.Logger
}
type Client interface {
	Exec(ctx context.Context, sql string, arguments ...any) (sql.Result, error)
	Query(ctx context.Context, sql string, args ...any) (*sql.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) *sql.Row
	Begin(ctx context.Context) (*sql.Tx, error)
}

func NewDb(log *logger.Logger, cfg *config.ConfigBot) *Db {
	dns := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		cfg.Postgress.Username, cfg.Postgress.Password, cfg.Postgress.Host, cfg.Postgress.Name)
	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFunc()
	db, err := sqlx.ConnectContext(ctx, "postgres", dns)
	if err != nil {
		log.ErrorErr(err)
		os.Exit(1)
	}

	database := &Db{
		db:  db,
		log: log,
	}
	return database
}

func (d *Db) DeleteOldClient(userid string) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	del := "delete from hs_compendium.corpmember where userid = $1"
	_, _ = d.db.ExecContext(ctx, del, userid)
	del = "delete from hs_compendium.list_users where userid = $1"
	_, _ = d.db.ExecContext(ctx, del, userid)
	del = "delete from hs_compendium.tech where userid = $1"
	_, _ = d.db.ExecContext(ctx, del, userid)
	del = "delete from hs_compendium.users where userid = $1"
	_, _ = d.db.ExecContext(ctx, del, userid)
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

package postgres

import (
	"compendium_s/config"
	"compendium_s/models"
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
	db  Client
	log *logger.Logger
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
		db:  pool,
		log: log,
	}
	return db
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

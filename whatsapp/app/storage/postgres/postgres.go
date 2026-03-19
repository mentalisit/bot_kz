package postgres

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"sync"
	"time"
	"whatsapp/config"
	"whatsapp/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mentalisit/logger"
)

type Db struct {
	db  Client
	log *logger.Logger
	sync.RWMutex
	RsBotConfig  map[string]models.CorporationConfigV2
	BridgeConfig map[string]models.Bridge2Config
	KzBotConfig  map[string]models.CorporationConfig
	pool         *pgxpool.Pool
}
type Client interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Begin(ctx context.Context) (pgx.Tx, error)
}

func NewDb(log *logger.Logger, cfg *config.ConfigBot) *Db {
	dns := fmt.Sprintf("postgres://%s:%s@%s/%s",
		cfg.Postgress.Username, cfg.Postgress.Password, cfg.Postgress.Host, cfg.Postgress.Name)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	pool, err := pgxpool.New(ctx, dns)
	if err != nil {
		slog.Error(err.Error())
		//log.ErrorErr(err)
		time.Sleep(5 * time.Second)
		os.Exit(1)
		//return err
	}
	db := &Db{
		db:           pool,
		log:          log,
		RsBotConfig:  make(map[string]models.CorporationConfigV2),
		BridgeConfig: make(map[string]models.Bridge2Config),
		KzBotConfig:  make(map[string]models.CorporationConfig),
		pool:         pool,
	}
	db.loadConfig()
	go db.StartConfigWatcher(context.Background())

	return db
}
func (d *Db) getContext() (ctx context.Context, cancel context.CancelFunc) {
	return context.WithTimeout(context.Background(), 5*time.Second)
}

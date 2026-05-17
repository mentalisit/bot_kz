package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"sync"
	"time"
	"whatsapp/config"
	"whatsapp/models"

	"github.com/jmoiron/sqlx"
	"github.com/mentalisit/logger"

	_ "github.com/lib/pq"
)

type Db struct {
	db  Client
	log *logger.Logger
	sync.RWMutex
	RsBotConfig  map[string]models.CorporationConfigV2
	BridgeConfig map[string]models.Bridge2Config
	KzBotConfig  map[string]models.CorporationConfig
	pool         *sqlx.DB
}
type Client interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (sql.Result, error)
	Query(ctx context.Context, sql string, args ...interface{}) (*sql.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) *sql.Row
	Begin(ctx context.Context) (*sql.Tx, error)
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
	database.loadConfig()
	go database.StartConfigWatcher(context.Background())

	return database
}
func (d *Db) getContext() (ctx context.Context, cancel context.CancelFunc) {
	return context.WithTimeout(context.Background(), 5*time.Second)
}

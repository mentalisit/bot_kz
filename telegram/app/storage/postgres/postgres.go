package postgres

import (
	"context"
	"fmt"
	"os"
	"sync"
	"telegram/config"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mentalisit/logger"
	"github.com/mentalisit/restapi/models"
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
		log.ErrorErr(err)
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
	db.CreateTables(context.Background())

	db.loadConfig()
	go db.StartConfigWatcher(context.Background())

	return db
}

func (d *Db) getContext() (ctx context.Context, cancel context.CancelFunc) {
	return context.WithTimeout(context.Background(), 10*time.Second)
}

// CreateTables создает необходимые таблицы в схеме telegram
func (d *Db) CreateTables(ctx context.Context) error {
	// Сначала создаем схему если она не существует
	queries := []string{
		`CREATE SCHEMA IF NOT EXISTS telegram`,

		// Таблица чатов
		`CREATE TABLE IF NOT EXISTS telegram.chats (
			chat_id BIGINT PRIMARY KEY,
			chat_name VARCHAR(255) NOT NULL,
			created_at TIMESTAMP DEFAULT NOW(),
			updated_at TIMESTAMP DEFAULT NOW()
		)`,

		// Таблица участников чатов
		`CREATE TABLE IF NOT EXISTS telegram.chat_members (
			chat_id BIGINT NOT NULL,
			user_id BIGINT NOT NULL,
			first_name VARCHAR(255),
			last_name VARCHAR(255),
			user_name VARCHAR(255),
			is_admin BOOLEAN DEFAULT FALSE,
			last_updated TIMESTAMP DEFAULT NOW(),
			PRIMARY KEY (chat_id, user_id)
		)`,

		// Таблица ролей
		`CREATE TABLE IF NOT EXISTS telegram.roles (
			id BIGSERIAL PRIMARY KEY,
			chat_id BIGINT NOT NULL,
			name VARCHAR(100) NOT NULL,
			created_by BIGINT NOT NULL,
			created_at TIMESTAMP DEFAULT NOW(),
			UNIQUE(chat_id, name)
		)`,

		// Таблица связи пользователей и ролей
		`CREATE TABLE IF NOT EXISTS telegram.user_roles (
			id BIGSERIAL PRIMARY KEY,
			user_id BIGINT NOT NULL,
			role_id BIGINT NOT NULL,
			chat_id BIGINT NOT NULL,
			assigned_at TIMESTAMP DEFAULT NOW(),
			UNIQUE(user_id, role_id),
			FOREIGN KEY (role_id) REFERENCES telegram.roles(id) ON DELETE CASCADE
		)`,

		// Таблица прав доступа
		`CREATE TABLE IF NOT EXISTS telegram.chat_permissions (
			chat_id BIGINT NOT NULL,
			user_id BIGINT NOT NULL,
			is_admin BOOLEAN DEFAULT FALSE,
			PRIMARY KEY (chat_id, user_id)
		)`,

		// Индексы для оптимизации
		`CREATE INDEX IF NOT EXISTS idx_chat_members_chat_id ON telegram.chat_members(chat_id)`,
		`CREATE INDEX IF NOT EXISTS idx_chat_members_user_id ON telegram.chat_members(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_roles_chat_id ON telegram.roles(chat_id)`,
		`CREATE INDEX IF NOT EXISTS idx_user_roles_user_chat ON telegram.user_roles(user_id, chat_id)`,
		`CREATE INDEX IF NOT EXISTS idx_user_roles_role_id ON telegram.user_roles(role_id)`,
		`CREATE INDEX IF NOT EXISTS idx_chat_permissions_chat_user ON telegram.chat_permissions(chat_id, user_id)`,
	}

	for _, query := range queries {
		_, err := d.db.Exec(ctx, query)
		if err != nil {
			return fmt.Errorf("failed to create table: %w, query: %s", err, query)
		}
	}
	return nil
}

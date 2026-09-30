package postgres

import (
	"context"
	"fmt"
	"sync"

	"github.com/jmoiron/sqlx"
	"github.com/mentalisit/conf/config"
	"github.com/mentalisit/conf/logger"
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
	dns          string
}

func NewDb(log *logger.Logger, db *sqlx.DB) *Db {

	database := &Db{
		db:           db,
		log:          log,
		RsBotConfig:  make(map[string]models.CorporationConfigV2),
		BridgeConfig: make(map[string]models.Bridge2Config),
		KzBotConfig:  make(map[string]models.CorporationConfig),
		dns:          config.Instance.GetDNS(),
	}
	database.CreateTables()

	database.loadConfig()
	go database.StartConfigWatcher(context.Background())

	return database
}

// CreateTables создает необходимые таблицы в схеме telegram
func (d *Db) CreateTables() error {
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
			last_updated TIMESTAMP DEFAULT NOW(),
			PRIMARY KEY (chat_id, user_id)
		)`,

		// Таблица участников чатов
		`CREATE TABLE IF NOT EXISTS telegram.members (
			user_id BIGINT NOT NULL,
			first_name VARCHAR(255),
			last_name VARCHAR(255),
			user_name VARCHAR(255),
			last_updated TIMESTAMP DEFAULT NOW(),
			PRIMARY KEY (user_id)
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

		// Таблица прав доступа
		`CREATE TABLE IF NOT EXISTS telegram.topic_cache (
			chat_id BIGINT NOT NULL,
			thread_id INTEGER NOT NULL,
			topic_name TEXT NOT NULL,
			updated_at TIMESTAMP NOT NULL DEFAULT now(),
		PRIMARY KEY ("chat_id", "thread_id")
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
		_, err := d.db.Exec(query)
		if err != nil {
			return fmt.Errorf("failed to create table: %w, query: %s", err, query)
		}
	}
	return nil
}

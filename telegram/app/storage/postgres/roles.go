package postgres

import (
	"context"
	"fmt"
	"telegram/models"
	"time"
)

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

// GetChatRoles возвращает роли определенного чата
func (d *Db) GetChatsRoles(ctx context.Context, chatID int64) ([]models.Role, error) {
	query := `SELECT id, chat_id, name FROM telegram.roles WHERE chat_id = $1`

	rows, err := d.db.Query(ctx, query, chatID)
	if err != nil {
		return nil, fmt.Errorf("failed to query chats roles: %w", err)
	}
	defer rows.Close()

	var roles []models.Role
	for rows.Next() {
		var r models.Role
		if err := rows.Scan(&r.ID, &r.ChatID, &r.Name); err != nil {
			return nil, fmt.Errorf("failed to scan roles: %w", err)
		}
		roles = append(roles, r)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return roles, nil
}

// createAllRole создает системную роль "all"
func (d *Db) createAllRole(ctx context.Context, chatID int64) (int64, error) {
	query := `
        INSERT INTO telegram.roles (chat_id, name, created_by, created_at)
        VALUES ($1, $2, $3, $4)
        RETURNING id
    `

	var roleID int64
	err := d.db.QueryRow(ctx, query,
		chatID, "all", 0, time.Now(), // created_by = 0 для системных ролей
	).Scan(&roleID)

	if err != nil {
		return 0, fmt.Errorf("failed to create 'all' role: %w", err)
	}

	fmt.Printf("Created 'all' role for chat %d with ID %d", chatID, roleID)
	return roleID, nil
}

// CreateRole создает новую роль
func (d *Db) CreateRole(ctx context.Context, role *models.Role) error {
	query := `
		INSERT INTO telegram.roles (chat_id, name, created_by, created_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	err := d.db.QueryRow(ctx, query,
		role.ChatID, role.Name, role.CreatedBy, time.Now(),
	).Scan(&role.ID)

	if err != nil {
		return fmt.Errorf("failed to create role: %w", err)
	}

	return nil
}

// SetUserRole назначает или снимает роль с пользователя
func (d *Db) SetUserRole(ctx context.Context, userID, roleID, chatID int64, assign bool) error {
	if assign {
		return d.JoinRole(ctx, userID, roleID, chatID)
	} else {
		return d.LeaveRole(ctx, userID, roleID)
	}
}

// DeleteRole удаляет роль
func (d *Db) DeleteRole(ctx context.Context, roleID, chatID int64) error {
	// Сначала проверяем, не является ли роль системной "all"
	var roleName string
	checkQuery := `SELECT name FROM telegram.roles WHERE id = $1 AND chat_id = $2`
	err := d.db.QueryRow(ctx, checkQuery, roleID, chatID).Scan(&roleName)
	if err != nil {
		return fmt.Errorf("failed to find role: %w", err)
	}

	if roleName == "all" {
		return fmt.Errorf("cannot delete system role 'all'")
	}

	query := `DELETE FROM telegram.roles WHERE id = $1 AND chat_id = $2`
	result, err := d.db.Exec(ctx, query, roleID, chatID)
	if err != nil {
		return fmt.Errorf("failed to delete role: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("role not found or access denied")
	}

	return nil
}

// UpdateRoleName обновляет название роли
func (d *Db) UpdateRoleName(ctx context.Context, roleID, chatID int64, newName string) error {
	var currentName string
	checkQuery := `SELECT name FROM telegram.roles WHERE id = $1 AND chat_id = $2`
	if err := d.db.QueryRow(ctx, checkQuery, roleID, chatID).Scan(&currentName); err != nil {
		return fmt.Errorf("failed to find role: %w", err)
	}

	if currentName == "all" {
		return fmt.Errorf("cannot rename system role 'all'")
	}

	query := `UPDATE telegram.roles SET name = $1 WHERE id = $2 AND chat_id = $3`
	result, err := d.db.Exec(ctx, query, newName, roleID, chatID)
	if err != nil {
		return fmt.Errorf("failed to update role name: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("role not found or access denied")
	}

	return nil
}

// GetRoleName возвращает название роли по ID
func (d *Db) GetRoleName(ctx context.Context, roleID int64, roleName *string) error {
	query := `SELECT name FROM telegram.roles WHERE id = $1`
	err := d.db.QueryRow(ctx, query, roleID).Scan(roleName)
	if err != nil {
		return fmt.Errorf("failed to get role name: %w", err)
	}
	return nil
}

// GetRoleName возвращает  ID роли
func (d *Db) GetRoleByName(ctx context.Context, roleName string, ChatID int64) (roleId int64, err error) {
	query := `SELECT id FROM telegram.roles WHERE name = $1 AND chat_id = $2`
	err = d.db.QueryRow(ctx, query, roleName, ChatID).Scan(&roleId)
	if err != nil {
		return 0, fmt.Errorf("failed to get role name: %w", err)
	}
	return roleId, nil
}

// SetChatAdmin назначает пользователя администратором чата
func (d *Db) SetChatAdmin(ctx context.Context, chatID, userID int64, isAdmin bool) error {
	query := `
        INSERT INTO telegram.chat_permissions (chat_id, user_id, is_admin)
        VALUES ($1, $2, $3)
        ON CONFLICT (chat_id, user_id) 
        DO UPDATE SET is_admin = EXCLUDED.is_admin
    `

	_, err := d.db.Exec(ctx, query, chatID, userID, isAdmin)
	if err != nil {
		return fmt.Errorf("failed to set chat admin: %w", err)
	}

	return nil
}

// RemoveChatAdmin удаляет права администратора у пользователя
func (d *Db) RemoveChatAdmin(ctx context.Context, chatID, userID int64) error {
	query := `DELETE FROM telegram.chat_permissions WHERE chat_id = $1 AND user_id = $2`

	result, err := d.db.Exec(ctx, query, chatID, userID)
	if err != nil {
		return fmt.Errorf("failed to remove chat admin: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("admin permission not found")
	}

	return nil
}

// RoleExists проверяет, существует ли роль с указанным именем в текущем чате и возвращает её ID
func (d *Db) RoleExists(ctx context.Context, chatID int64, roleName string) (int64, error) {
	query := `SELECT id FROM telegram.roles WHERE chat_id = $1 AND name = $2`

	var roleID int64
	err := d.db.QueryRow(ctx, query, chatID, roleName).Scan(&roleID)
	if err != nil {
		return 0, fmt.Errorf("failed to check role existence: %w", err)
	}

	return roleID, nil
}

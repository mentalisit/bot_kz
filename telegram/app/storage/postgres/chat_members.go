package postgres

import (
	"context"
	"fmt"
	"telegram/models"

	"github.com/jackc/pgx/v5"
)

// getAllRoleMemberCount возвращает количество участников в роли "all" (все пользователи чата)
func (d *Db) getAllRoleMemberCount(ctx context.Context, chatID int64) (int, error) {
	query := `
        SELECT COUNT(*) 
        FROM telegram.chat_members 
        WHERE chat_id = $1
    `

	var count int
	err := d.db.QueryRow(ctx, query, chatID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get member count: %w", err)
	}

	return count, nil
}

// GetChatUsers возвращает пользователей чата с их ролями
func (d *Db) GetChatUsers(ctx context.Context, chatID int64) ([]models.User, error) {
	query := `
        SELECT 
            cm.user_id,
            cm.first_name,
            cm.last_name,
            cm.user_name,
            cm.is_admin
        FROM telegram.chat_members cm
        WHERE cm.chat_id = $1
        ORDER BY cm.is_admin DESC, cm.first_name, cm.last_name
    `

	rows, err := d.db.Query(ctx, query, chatID)
	if err != nil {
		return nil, fmt.Errorf("failed to query chat users: %w", err)
	}
	defer rows.Close()

	var users []models.User
	userMap := make(map[int64]*models.User)

	for rows.Next() {
		var user models.User
		if err := rows.Scan(
			&user.ID,
			&user.FirstName,
			&user.LastName,
			&user.UserName,
			&user.IsAdmin,
		); err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		user.Roles = make(map[int64]bool)
		users = append(users, user)
		userMap[user.ID] = &users[len(users)-1]
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	// Находим ID роли "all" для этого чата
	var allRoleID int64
	roleQuery := `SELECT id FROM telegram.roles WHERE chat_id = $1 AND name = 'all'`
	err = d.db.QueryRow(ctx, roleQuery, chatID).Scan(&allRoleID)
	if err != nil && err != pgx.ErrNoRows {
		return nil, fmt.Errorf("failed to get 'all' role ID: %w", err)
	}

	// Если роль "all" существует, добавляем всех пользователей в нее
	if allRoleID != 0 {
		for i := range users {
			users[i].Roles[allRoleID] = true
		}
	}

	// Теперь загружаем остальные роли для каждого пользователя
	if len(users) > 0 {
		rolesQuery := `
            SELECT user_id, role_id 
            FROM telegram.user_roles 
            WHERE chat_id = $1
        `

		roleRows, err := d.db.Query(ctx, rolesQuery, chatID)
		if err != nil {
			return nil, fmt.Errorf("failed to query user roles: %w", err)
		}
		defer roleRows.Close()

		for roleRows.Next() {
			var userID, roleID int64
			if err := roleRows.Scan(&userID, &roleID); err != nil {
				return nil, fmt.Errorf("failed to scan user role: %w", err)
			}
			if user, exists := userMap[userID]; exists {
				user.Roles[roleID] = true
			}
		}

		if err := roleRows.Err(); err != nil {
			return nil, fmt.Errorf("error iterating role rows: %w", err)
		}
	}

	return users, nil
}

// IsChatAdmin проверяет, является ли пользователь администратором чата
func (d *Db) IsChatAdmin(ctx context.Context, chatID, userID int64) (bool, error) {
	// Сначала проверяем в таблице chat_permissions
	query := `SELECT is_admin FROM telegram.chat_permissions WHERE chat_id = $1 AND user_id = $2`
	var isAdmin bool
	err := d.db.QueryRow(ctx, query, chatID, userID).Scan(&isAdmin)

	if err == nil {
		return isAdmin, nil
	}

	if err != pgx.ErrNoRows {
		return false, fmt.Errorf("failed to check admin status: %w", err)
	}

	// Если записи нет, проверяем в таблице chat_members
	query = `SELECT is_admin FROM telegram.chat_members WHERE chat_id = $1 AND user_id = $2`
	err = d.db.QueryRow(ctx, query, chatID, userID).Scan(&isAdmin)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil // Пользователь не найден в чате
		}
		return false, fmt.Errorf("failed to get user from chat_members: %w", err)
	}

	return isAdmin, nil
}

// UpdateUserCache обновляет кэш пользователей чата
func (d *Db) UpdateUserCache(ctx context.Context, chatID int64, users map[int64]models.User) error {
	// Начнем транзакцию
	tx, err := d.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Очищаем старые данные для этого чата
	_, err = tx.Exec(ctx, "DELETE FROM telegram.chat_members WHERE chat_id = $1", chatID)
	if err != nil {
		return fmt.Errorf("failed to clear chat members: %w", err)
	}

	// Вставляем новых пользователей
	for userID, user := range users {
		_, err := tx.Exec(ctx, `
            INSERT INTO telegram.chat_members (chat_id, user_id, first_name, last_name, user_name, is_admin)
            VALUES ($1, $2, $3, $4, $5, $6)
        `, chatID, userID, user.FirstName, user.LastName, user.UserName, user.IsAdmin)

		if err != nil {
			return fmt.Errorf("failed to insert chat member %d: %w", userID, err)
		}
	}

	// Коммитим транзакцию
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// RemoveUserFromChat удаляет пользователя из чата
func (d *Db) RemoveUserFromChat(ctx context.Context, chatID, userID int64) error {
	query := `DELETE FROM telegram.chat_members WHERE chat_id = $1 AND user_id = $2`

	_, err := d.db.Exec(ctx, query, chatID, userID)
	if err != nil {
		return fmt.Errorf("failed to remove user from chat: %w", err)
	}

	query = `DELETE FROM telegram.chat_permissions WHERE chat_id = $1 AND user_id = $2`
	_, _ = d.db.Exec(ctx, query, chatID, userID)

	query = `DELETE FROM telegram.user_roles WHERE chat_id = $1 AND user_id = $2`
	_, _ = d.db.Exec(ctx, query, chatID, userID)

	return nil
}

// UpdateUserInfo обновляет информацию о пользователе
func (d *Db) UpdateUserInfo(ctx context.Context, userID int64, firstName, lastName, userName string) error {
	query := `
        UPDATE telegram.chat_members 
        SET first_name = $1, last_name = $2, user_name = $3, last_updated = NOW()
        WHERE user_id = $4
    `

	result, err := d.db.Exec(ctx, query, firstName, lastName, userName, userID)
	if err != nil {
		return fmt.Errorf("failed to update user info: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// AddUserToChat добавляет пользователя в чат
func (d *Db) AddUserToChat(ctx context.Context, chatID int64, user models.User) error {
	query := `
        INSERT INTO telegram.chat_members (chat_id, user_id, first_name, last_name, user_name, is_admin)
        VALUES ($1, $2, $3, $4, $5, $6)
        ON CONFLICT (chat_id, user_id) 
        DO UPDATE SET
            first_name = EXCLUDED.first_name,
            last_name = EXCLUDED.last_name,
            user_name = EXCLUDED.user_name,
            is_admin = EXCLUDED.is_admin,
            last_updated = NOW()
    `

	_, err := d.db.Exec(ctx, query, chatID, user.ID, user.FirstName, user.LastName, user.UserName, user.IsAdmin)
	if err != nil {
		return fmt.Errorf("failed to add user to chat: %w", err)
	}

	return nil
}

// GetChatAdmins возвращает список администраторов чата
func (d *Db) GetChatAdmins(ctx context.Context, chatID int64) ([]models.User, error) {
	query := `
        SELECT cm.user_id, cm.first_name, cm.last_name, cm.user_name, cm.is_admin
        FROM telegram.chat_members cm
        LEFT JOIN telegram.chat_permissions cp ON cm.chat_id = cp.chat_id AND cm.user_id = cp.user_id
        WHERE cm.chat_id = $1 AND (cm.is_admin = true OR cp.is_admin = true)
        ORDER BY cm.first_name, cm.last_name
    `

	rows, err := d.db.Query(ctx, query, chatID)
	if err != nil {
		return nil, fmt.Errorf("failed to query chat admins: %w", err)
	}
	defer rows.Close()

	var admins []models.User
	for rows.Next() {
		var admin models.User
		if err := rows.Scan(
			&admin.ID,
			&admin.FirstName,
			&admin.LastName,
			&admin.UserName,
			&admin.IsAdmin,
		); err != nil {
			return nil, fmt.Errorf("failed to scan admin: %w", err)
		}
		admins = append(admins, admin)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return admins, nil
}

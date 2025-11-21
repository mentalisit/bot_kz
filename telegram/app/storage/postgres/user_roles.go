package postgres

import (
	"context"
	"fmt"
	"telegram/models"
	"time"
)

// GetChatRoles возвращает роли в чате с информацией о подписке пользователя
func (d *Db) GetChatRoles(ctx context.Context, chatID, userID int64) ([]models.Role, error) {
	// Сначала получаем роли из базы данных
	query := `
        SELECT r.id, r.chat_id, r.name, r.created_by, r.created_at,
               COUNT(ur.user_id) as member_count,
               EXISTS(SELECT 1 FROM telegram.user_roles WHERE user_id = $2 AND role_id = r.id) as is_member
        FROM telegram.roles r
        LEFT JOIN telegram.user_roles ur ON r.id = ur.role_id
        WHERE r.chat_id = $1
        GROUP BY r.id
        ORDER BY r.created_at
    `

	rows, err := d.db.Query(ctx, query, chatID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query chat roles: %w", err)
	}
	defer rows.Close()

	var roles []models.Role
	var hasAllRole bool
	var allRoleID int64

	for rows.Next() {
		var role models.Role
		if err := rows.Scan(
			&role.ID, &role.ChatID, &role.Name, &role.CreatedBy, &role.CreatedAt,
			&role.MemberCount, &role.IsMember,
		); err != nil {
			return nil, fmt.Errorf("failed to scan role: %w", err)
		}
		roles = append(roles, role)

		// Проверяем есть ли роль "all"
		if role.Name == "all" {
			hasAllRole = true
			allRoleID = role.ID
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	// Если роли "all" нет - создаем ее
	if !hasAllRole {
		allRoleID, err = d.createAllRole(ctx, chatID)
		if err != nil {
			fmt.Printf("Warning: failed to create 'all' role: %v", err)
		} else {
			// Добавляем виртуальную роль "all" в результат
			allRole := models.Role{
				ID:          allRoleID,
				ChatID:      chatID,
				Name:        "all",
				CreatedBy:   0, // Системная роль
				CreatedAt:   time.Now(),
				MemberCount: 0,    // Будет рассчитано ниже
				IsMember:    true, // Все пользователи автоматически в этой роли
			}
			roles = append([]models.Role{allRole}, roles...)
		}
	}

	// Обновляем количество участников для роли "all"
	for i, role := range roles {
		if role.Name == "all" {
			memberCount, err := d.getAllRoleMemberCount(ctx, chatID)
			if err != nil {
				fmt.Printf("Warning: failed to get member count for 'all' role: %v", err)
				memberCount = 0
			}
			roles[i].MemberCount = memberCount
			roles[i].IsMember = true // Все пользователи всегда в роли "all"
			break
		}
	}

	return roles, nil
}

// JoinRole добавляет пользователя в роль
func (d *Db) JoinRole(ctx context.Context, userID, roleID, chatID int64) error {
	// Проверяем, не является ли роль "all"
	var roleName string
	checkQuery := `SELECT name FROM telegram.roles WHERE id = $1`
	err := d.db.QueryRow(ctx, checkQuery, roleID).Scan(&roleName)
	if err != nil {
		return fmt.Errorf("failed to find role: %w", err)
	}

	if roleName == "all" {
		return fmt.Errorf("cannot manually join system role 'all' - all users are automatically members")
	}

	query := `
        INSERT INTO telegram.user_roles (user_id, role_id, chat_id)
        VALUES ($1, $2, $3)
        ON CONFLICT (user_id, role_id) DO NOTHING
    `

	_, err = d.db.Exec(ctx, query, userID, roleID, chatID)
	if err != nil {
		return fmt.Errorf("failed to join role: %w", err)
	}

	return nil
}

// LeaveRole удаляет пользователя из роли
func (d *Db) LeaveRole(ctx context.Context, userID, roleID int64) error {
	// Проверяем, не является ли роль "all"
	var roleName string
	checkQuery := `SELECT name FROM telegram.roles WHERE id = $1`
	err := d.db.QueryRow(ctx, checkQuery, roleID).Scan(&roleName)
	if err != nil {
		return fmt.Errorf("failed to find role: %w", err)
	}

	if roleName == "all" {
		return fmt.Errorf("cannot leave system role 'all'")
	}

	query := `DELETE FROM telegram.user_roles WHERE user_id = $1 AND role_id = $2`
	_, err = d.db.Exec(ctx, query, userID, roleID)
	if err != nil {
		return fmt.Errorf("failed to leave role: %w", err)
	}

	return nil
}

// IsUserSubscribedToRole проверяет, подписан ли пользователь на указанную роль
func (d *Db) IsUserSubscribedToRole(ctx context.Context, userID, roleID int64) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM telegram.user_roles WHERE user_id = $1 AND role_id = $2)`

	var isSubscribed bool
	err := d.db.QueryRow(ctx, query, userID, roleID).Scan(&isSubscribed)
	if err != nil {
		return false, fmt.Errorf("failed to check user role subscription: %w", err)
	}

	return isSubscribed, nil
}

package postgresV2

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"rs/models"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// ==================== Chats ====================

// GetUserChats возвращает список чатов пользователя
func (d *Db) GetUserChats(ctx context.Context, userID int64) ([]models.Chat, error) {
	query := `
        SELECT DISTINCT c.chat_id, c.chat_name 
        FROM telegram.chats c
        INNER JOIN telegram.chat_members cm ON c.chat_id = cm.chat_id
        WHERE cm.user_id = $1
        ORDER BY c.chat_name
    `

	rows, err := d.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query user chats: %w", err)
	}
	defer rows.Close()

	var chats []models.Chat
	for rows.Next() {
		var chat models.Chat
		if err := rows.Scan(&chat.ChatID, &chat.ChatName); err != nil {
			return nil, fmt.Errorf("failed to scan chat: %w", err)
		}
		chats = append(chats, chat)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return chats, nil
}

// UpdateChatTitle обновляет название чата
func (d *Db) UpdateChatTitle(ctx context.Context, chatID int64, chatTitle string) error {
	query := `
        UPDATE telegram.chats 
        SET chat_name = $1, updated_at = NOW() 
        WHERE chat_id = $2 AND chat_name != $1
    `

	_, err := d.db.Exec(query, chatTitle, chatID)
	if err != nil {
		return fmt.Errorf("failed to update chat title: %w", err)
	}

	return nil
}

// ==================== Roles ====================

// GetChatRoles возвращает роли в чате с информацией о подписке пользователя
func (d *Db) GetChatRoles(ctx context.Context, chatID, userID int64) ([]models.Role, error) {
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

	rows, err := d.db.Query(query, chatID, userID)
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

		if role.Name == "all" {
			hasAllRole = true
			allRoleID = role.ID
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	// Если роли "all" нет — создаем
	if !hasAllRole {
		allRoleID, err = d.createAllRole(chatID)
		if err != nil {
			fmt.Printf("Warning: failed to create 'all' role: %v", err)
		} else {
			allRole := models.Role{
				ID:          allRoleID,
				ChatID:      chatID,
				Name:        "all",
				CreatedBy:   0,
				CreatedAt:   time.Now(),
				MemberCount: 0,
				IsMember:    true,
			}
			roles = append([]models.Role{allRole}, roles...)
		}
	}

	// Обновляем количество участников для роли "all"
	for i, role := range roles {
		if role.Name == "all" {
			memberCount, err := d.getAllRoleMemberCount(chatID)
			if err != nil {
				fmt.Printf("Warning: failed to get member count for 'all' role: %v", err)
				memberCount = 0
			}
			roles[i].MemberCount = memberCount
			roles[i].IsMember = true
			break
		}
	}

	return roles, nil
}

// createAllRole создает системную роль "all"
func (d *Db) createAllRole(chatID int64) (int64, error) {
	query := `
        INSERT INTO telegram.roles (chat_id, name, created_by, created_at)
        VALUES ($1, $2, $3, $4)
        RETURNING id
    `

	var roleID int64
	err := d.db.QueryRow(query, chatID, "all", 0, time.Now()).Scan(&roleID)
	if err != nil {
		return 0, fmt.Errorf("failed to create 'all' role: %w", err)
	}

	return roleID, nil
}

// getAllRoleMemberCount возвращает количество участников в чате (для роли "all")
func (d *Db) getAllRoleMemberCount(chatID int64) (int, error) {
	query := `SELECT COUNT(*) FROM telegram.chat_members WHERE chat_id = $1`

	var count int
	err := d.db.QueryRow(query, chatID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get member count: %w", err)
	}

	return count, nil
}

// CreateRole создает новую роль
func (d *Db) CreateRole(ctx context.Context, role *models.Role) error {
	query := `
		INSERT INTO telegram.roles (chat_id, name, created_by, created_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	err := d.db.QueryRow(query,
		role.ChatID, role.Name, role.CreatedBy, time.Now(),
	).Scan(&role.ID)

	if err != nil {
		return fmt.Errorf("failed to create role: %w", err)
	}

	return nil
}

// UpdateRoleName обновляет название роли
func (d *Db) UpdateRoleName(ctx context.Context, roleID, chatID int64, newName string) error {
	var currentName string
	checkQuery := `SELECT name FROM telegram.roles WHERE id = $1 AND chat_id = $2`
	if err := d.db.QueryRow(checkQuery, roleID, chatID).Scan(&currentName); err != nil {
		return fmt.Errorf("failed to find role: %w", err)
	}

	if currentName == "all" {
		return fmt.Errorf("cannot rename system role 'all'")
	}

	query := `UPDATE telegram.roles SET name = $1 WHERE id = $2 AND chat_id = $3`
	result, err := d.db.Exec(query, newName, roleID, chatID)
	if err != nil {
		return fmt.Errorf("failed to update role name: %w", err)
	}

	r, _ := result.RowsAffected()
	if r == 0 {
		return fmt.Errorf("role not found or access denied")
	}

	return nil
}

// DeleteRole удаляет роль
func (d *Db) DeleteRole(ctx context.Context, roleID, chatID int64) error {
	var roleName string
	checkQuery := `SELECT name FROM telegram.roles WHERE id = $1 AND chat_id = $2`
	err := d.db.QueryRow(checkQuery, roleID, chatID).Scan(&roleName)
	if err != nil {
		return fmt.Errorf("failed to find role: %w", err)
	}

	if roleName == "all" {
		return fmt.Errorf("cannot delete system role 'all'")
	}

	query := `DELETE FROM telegram.roles WHERE id = $1 AND chat_id = $2`
	result, err := d.db.Exec(query, roleID, chatID)
	if err != nil {
		return fmt.Errorf("failed to delete role: %w", err)
	}

	r, _ := result.RowsAffected()
	if r == 0 {
		return fmt.Errorf("role not found or access denied")
	}

	return nil
}

// GetRoleName возвращает название роли по ID
func (d *Db) GetRoleName(ctx context.Context, roleID int64, roleName *string) error {
	query := `SELECT name FROM telegram.roles WHERE id = $1`
	err := d.db.QueryRow(query, roleID).Scan(roleName)
	if err != nil {
		return fmt.Errorf("failed to get role name: %w", err)
	}
	return nil
}

// ==================== User Roles ====================

// JoinRole добавляет пользователя в роль
func (d *Db) JoinRole(userID, roleID, chatID int64) error {
	var roleName string
	checkQuery := `SELECT name FROM telegram.roles WHERE id = $1`
	err := d.db.QueryRow(checkQuery, roleID).Scan(&roleName)
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

	_, err = d.db.Exec(query, userID, roleID, chatID)
	if err != nil {
		return fmt.Errorf("failed to join role: %w", err)
	}

	return nil
}

// LeaveRole удаляет пользователя из роли
func (d *Db) LeaveRole(userID, roleID int64) error {
	var roleName string
	checkQuery := `SELECT name FROM telegram.roles WHERE id = $1`
	err := d.db.QueryRow(checkQuery, roleID).Scan(&roleName)
	if err != nil {
		return fmt.Errorf("failed to find role: %w", err)
	}

	if roleName == "all" {
		return fmt.Errorf("cannot leave system role 'all'")
	}

	query := `DELETE FROM telegram.user_roles WHERE user_id = $1 AND role_id = $2`
	_, err = d.db.Exec(query, userID, roleID)
	if err != nil {
		return fmt.Errorf("failed to leave role: %w", err)
	}

	return nil
}

// IsUserSubscribedToRole проверяет, подписан ли пользователь на указанную роль
func (d *Db) IsUserSubscribedToRole(userID, roleID int64) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM telegram.user_roles WHERE user_id = $1 AND role_id = $2)`

	err := d.db.QueryRow(query, userID, roleID).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

// SetUserRole назначает или снимает роль с пользователя
func (d *Db) SetUserRole(ctx context.Context, userID, roleID, chatID int64, assign bool) error {
	if assign {
		return d.JoinRole(userID, roleID, chatID)
	}
	return d.LeaveRole(userID, roleID)
}

// ==================== Chat Members ====================

// GetChatUsers возвращает пользователей чата с их ролями
func (d *Db) GetChatUsers(ctx context.Context, chatID int64) ([]models.ChatUser, error) {
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

	rows, err := d.db.Query(query, chatID)
	if err != nil {
		return nil, fmt.Errorf("failed to query chat users: %w", err)
	}
	defer rows.Close()

	var users []models.ChatUser
	userMap := make(map[int64]*models.ChatUser)

	for rows.Next() {
		var user models.ChatUser
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
		if user.UserName != "mentalisit" {
			users = append(users, user)
		}
		userMap[user.ID] = &users[len(users)-1]
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	// Роль "all"
	var allRoleID int64
	roleQuery := `SELECT id FROM telegram.roles WHERE chat_id = $1 AND name = 'all'`
	err = d.db.QueryRow(roleQuery, chatID).Scan(&allRoleID)
	if err != nil && err != pgx.ErrNoRows {
		return nil, fmt.Errorf("failed to get 'all' role ID: %w", err)
	}

	if allRoleID != 0 {
		for i := range users {
			users[i].Roles[allRoleID] = true
		}
	}

	// Загружаем остальные роли
	if len(users) > 0 {
		rolesQuery := `
            SELECT user_id, role_id 
            FROM telegram.user_roles 
            WHERE chat_id = $1
        `

		roleRows, err := d.db.Query(rolesQuery, chatID)
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
	query := `SELECT is_admin FROM telegram.chat_permissions WHERE chat_id = $1 AND user_id = $2`
	var isAdmin bool
	err := d.db.QueryRow(query, chatID, userID).Scan(&isAdmin)

	if err == nil {
		return isAdmin, nil
	}

	if err != pgx.ErrNoRows {
		return false, fmt.Errorf("failed to check admin status: %w", err)
	}

	query = `SELECT is_admin FROM telegram.chat_members WHERE chat_id = $1 AND user_id = $2`
	err = d.db.QueryRow(query, chatID, userID).Scan(&isAdmin)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("failed to get user from chat_members: %w", err)
	}

	return isAdmin, nil
}

// GetChatAdmins возвращает список администраторов чата
func (d *Db) GetChatAdmins(ctx context.Context, chatID int64) ([]models.ChatUser, error) {
	query := `
        SELECT cm.user_id, cm.first_name, cm.last_name, cm.user_name, cm.is_admin
        FROM telegram.chat_members cm
        LEFT JOIN telegram.chat_permissions cp ON cm.chat_id = cp.chat_id AND cm.user_id = cp.user_id
        WHERE cm.chat_id = $1 AND (cm.is_admin = true OR cp.is_admin = true)
        ORDER BY cm.first_name, cm.last_name
    `

	rows, err := d.db.Query(query, chatID)
	if err != nil {
		return nil, fmt.Errorf("failed to query chat admins: %w", err)
	}
	defer rows.Close()

	var admins []models.ChatUser
	for rows.Next() {
		var admin models.ChatUser
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

// SetChatAdmin назначает пользователя администратором чата
func (d *Db) SetChatAdmin(ctx context.Context, chatID, userID int64, isAdmin bool) error {
	query := `
        INSERT INTO telegram.chat_permissions (chat_id, user_id, is_admin)
        VALUES ($1, $2, $3)
        ON CONFLICT (chat_id, user_id) 
        DO UPDATE SET is_admin = EXCLUDED.is_admin
    `

	_, err := d.db.Exec(query, chatID, userID, isAdmin)
	if err != nil {
		return fmt.Errorf("failed to set chat admin: %w", err)
	}

	return nil
}

// RemoveChatAdmin удаляет права администратора у пользователя
func (d *Db) RemoveChatAdmin(ctx context.Context, chatID, userID int64) error {
	query := `DELETE FROM telegram.chat_permissions WHERE chat_id = $1 AND user_id = $2`

	_, err := d.db.Exec(query, chatID, userID)
	if err != nil {
		return fmt.Errorf("failed to remove chat admin: %w", err)
	}

	return nil
}

// ==================== Corp Members ====================

// GetCorpMembersMyCompendium возвращает участников корпорации из my_compendium
func (d *Db) GetCorpMembersMyCompendium(ctx context.Context, chatID int64) ([]models.CompendiumCorpMember, error) {
	gid, err := d.GetGuildUUIDMyCompendium(chatID)
	if err != nil {
		return nil, err
	}

	query := `SELECT cm.uid, ma.nickname, ma.avatarurl
		FROM my_compendium.corpmember cm
		JOIN my_compendium.multi_accounts ma ON cm.uid = ma.uuid
		WHERE $1 = ANY(cm.guildIds)`

	rows, err := d.db.Query(query, gid)
	if err != nil {
		return nil, fmt.Errorf("failed to query corp members: %w", err)
	}
	defer rows.Close()

	var users []models.CompendiumCorpMember
	for rows.Next() {
		var user models.CompendiumCorpMember
		var userUUID uuid.UUID
		if err := rows.Scan(&userUUID, &user.Username, &user.AvatarURL); err != nil {
			return nil, fmt.Errorf("failed to scan corp member: %w", err)
		}
		user.UserID = userUUID.String()
		user.TableSource = "my_compendium"
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating corp members rows: %w", err)
	}

	return users, nil
}

// GetCorpMembersHSCompendium возвращает участников корпорации из hs_compendium
func (d *Db) GetCorpMembersHSCompendium(ctx context.Context, chatID int64) ([]models.CompendiumCorpMember, error) {
	gid, err := d.GetGuildUUIDMyCompendium(chatID)
	if err != nil {
		return nil, err
	}
	guildId := gid.String()

	sel := "SELECT username,userid,avatarurl FROM hs_compendium.corpmember WHERE guildid = $1"
	results, err := d.db.Query(sel, guildId)
	if err != nil {
		return nil, err
	}
	defer results.Close()

	var mm []models.CompendiumCorpMember
	for results.Next() {
		var t models.CompendiumCorpMember
		err = results.Scan(&t.Username, &t.UserID, &t.AvatarURL)
		if err != nil {
			return nil, err
		}
		t.TableSource = "hs_compendium"
		t.Username = changeCorpMemberName(t)
		mm = append(mm, t)
	}
	return mm, nil
}

func changeCorpMemberName(m models.CompendiumCorpMember) string {
	if len(m.UserID) < 13 {
		return "(TG) " + m.Username
	} else if len(m.UserID) > 16 {
		return "(DS) " + m.Username
	}
	return m.Username
}

// GetGuildUUIDMyCompendium находит UUID гильдии по chatID
func (d *Db) GetGuildUUIDMyCompendium(chatID int64) (gid *uuid.UUID, err error) {
	guildQuery := `SELECT gid FROM my_compendium.guilds WHERE EXISTS (
		SELECT 1 FROM jsonb_object_keys(channels) AS k
		CROSS JOIN jsonb_array_elements_text(channels->k) AS v
		WHERE v = $1
	)`
	err = d.db.QueryRow(guildQuery, strconv.FormatInt(chatID, 10)).Scan(&gid)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("guild not found for chat %d", chatID)
		}
		return nil, fmt.Errorf("failed to find guild: %w", err)
	}
	return gid, nil
}

// RemoveCorpMemberMyCompendium удаляет участника из my_compendium
func (d *Db) RemoveCorpMemberMyCompendium(ctx context.Context, chatID int64, userID string) error {
	gid, err := d.GetGuildUUIDMyCompendium(chatID)
	if err != nil {
		return err
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID format: %w", err)
	}

	var currentGuilds models.UUIDArray
	guildsQuery := `SELECT guildIds FROM my_compendium.corpmember WHERE uid = $1`
	err = d.db.QueryRow(guildsQuery, userUUID).Scan(&currentGuilds)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("corp member not found for uid %s", userID)
		}
		return fmt.Errorf("failed to get current guilds: %w", err)
	}

	var newGuilds models.UUIDArray
	for _, guild := range currentGuilds {
		if guild != *gid {
			newGuilds = append(newGuilds, guild)
		}
	}

	if len(newGuilds) == 0 {
		deleteQuery := `DELETE FROM my_compendium.corpmember WHERE uid = $1`
		_, err = d.db.Exec(deleteQuery, userUUID)
		if err != nil {
			return fmt.Errorf("failed to delete corp member: %w", err)
		}
	} else {
		updateQuery := `UPDATE my_compendium.corpmember SET guildIds = $1 WHERE uid = $2`
		_, err = d.db.Exec(updateQuery, newGuilds, userUUID)
		if err != nil {
			return fmt.Errorf("failed to update corp member: %w", err)
		}
	}

	return nil
}

// RemoveCorpMemberHSCompendium удаляет участника из hs_compendium
func (d *Db) RemoveCorpMemberHSCompendium(ctx context.Context, chatID int64, userID string) error {
	gid, _ := d.GetGuildUUIDMyCompendium(chatID)
	deleteQuery := `DELETE FROM hs_compendium.corpmember WHERE userid = $1 and guildid = $2`
	_, err := d.db.Exec(deleteQuery, userID, gid.String())
	if err != nil {
		return fmt.Errorf("failed to delete corp member: %w", err)
	}
	return nil
}

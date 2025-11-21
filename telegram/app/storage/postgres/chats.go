package postgres

import (
	"context"
	"fmt"
	"telegram/models"
)

// GetChat возвращает информацию о чате
func (d *Db) GetChat(chatID int64) (*models.Chat, error) {
	ctx, cancel := d.getContext()
	defer cancel()

	query := `SELECT chat_id, chat_name FROM telegram.chats WHERE chat_id = $1`

	var chat models.Chat
	err := d.db.QueryRow(ctx, query, chatID).Scan(&chat.ChatID, &chat.ChatName)
	if err != nil {
		return nil, fmt.Errorf("failed to get chat: %w", err)
	}

	return &chat, nil
}

// UpdateChatTitle обновляет название чата если оно изменилось
func (d *Db) UpdateChatTitle(ctx context.Context, chatID int64, chatTitle string) error {
	query := `
        UPDATE telegram.chats 
        SET chat_name = $1, updated_at = NOW() 
        WHERE chat_id = $2 AND chat_name != $1
    `

	result, err := d.db.Exec(ctx, query, chatTitle, chatID)
	if err != nil {
		return fmt.Errorf("failed to update chat title: %w", err)
	}

	// Логируем если название было обновлено
	if result.RowsAffected() > 0 {
		fmt.Printf("Updated chat title for chat %d: %s", chatID, chatTitle)
	}

	return nil
}

// CreateOrUpdateChat создает или обновляет информацию о чате
func (d *Db) CreateOrUpdateChat(ctx context.Context, chatID int64, chatName string) error {
	query := `
        INSERT INTO telegram.chats (chat_id, chat_name, updated_at)
        VALUES ($1, $2, NOW())
        ON CONFLICT (chat_id) 
        DO UPDATE SET 
            chat_name = EXCLUDED.chat_name,
            updated_at = NOW()
    `

	_, err := d.db.Exec(ctx, query, chatID, chatName)
	if err != nil {
		return fmt.Errorf("failed to create/update chat: %w", err)
	}

	return nil
}

// UpdateChatName обновляет название чата
func (d *Db) UpdateChatName(ctx context.Context, chatID int64, newChatName string) error {
	query := `UPDATE telegram.chats SET chat_name = $1, updated_at = NOW() WHERE chat_id = $2`

	result, err := d.db.Exec(ctx, query, newChatName, chatID)
	if err != nil {
		return fmt.Errorf("failed to update chat name: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("chat not found")
	}

	return nil
}

// DeleteChat удаляет чат и все связанные данные
func (d *Db) DeleteChat(ctx context.Context, chatID int64) error {
	query := `DELETE FROM telegram.chats WHERE chat_id = $1`

	result, err := d.db.Exec(ctx, query, chatID)
	if err != nil {
		return fmt.Errorf("failed to delete chat: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("chat not found")
	}

	return nil
}

// GetUserChats возвращает список чатов пользователя
func (d *Db) GetUserChats(ctx context.Context, userID int64) ([]models.Chat, error) {
	query := `
        SELECT DISTINCT c.chat_id, c.chat_name 
        FROM telegram.chats c
        INNER JOIN telegram.chat_members cm ON c.chat_id = cm.chat_id
        WHERE cm.user_id = $1
        ORDER BY c.chat_name
    `

	rows, err := d.db.Query(ctx, query, userID)
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

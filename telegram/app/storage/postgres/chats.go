package postgres

import (
	"fmt"
	"telegram/models2"
)

// GetChat возвращает информацию о чате
func (d *Db) GetChat(chatID int64) (models2.Chat, error) {

	query := `SELECT chat_id, chat_name FROM telegram.chats WHERE chat_id = $1`

	var chat models2.Chat
	err := d.db.QueryRow(query, chatID).Scan(&chat.ChatID, &chat.ChatName)
	if err != nil {
		return chat, fmt.Errorf("failed to get chat: %w", err)
	}

	return chat, nil
}

// GetChats возвращает список всех чатов из базы данных
func (d *Db) GetChats() ([]models2.Chat, error) {

	query := `SELECT chat_id, chat_name FROM telegram.chats`

	rows, err := d.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query chats: %w", err)
	}
	defer rows.Close()

	var chats []models2.Chat
	for rows.Next() {
		var chat models2.Chat
		err := rows.Scan(&chat.ChatID, &chat.ChatName)
		if err != nil {
			return nil, fmt.Errorf("failed to scan chat row: %w", err)
		}
		chats = append(chats, chat)
	}

	// Проверка на ошибки, возникшие во время итерации
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}

	return chats, nil
}

// UpdateChatTitle обновляет название чата если оно изменилось
func (d *Db) UpdateChatTitle(chatID int64, chatTitle string) error {
	query := `
        UPDATE telegram.chats 
        SET chat_name = $1, updated_at = NOW() 
        WHERE chat_id = $2 AND chat_name != $1
    `

	result, err := d.db.Exec(query, chatTitle, chatID)
	if err != nil {
		return fmt.Errorf("failed to update chat title: %w", err)
	}

	// Логируем если название было обновлено
	r, _ := result.RowsAffected()
	if r > 0 {
		fmt.Printf("Updated chat title for chat %d: %s", chatID, chatTitle)
	}

	return nil
}

// CreateOrUpdateChat создает или обновляет информацию о чате
func (d *Db) CreateOrUpdateChat(chatID int64, chatName string) error {
	query := `
        INSERT INTO telegram.chats (chat_id, chat_name, updated_at)
        VALUES ($1, $2, NOW())
        ON CONFLICT (chat_id) 
        DO UPDATE SET 
            chat_name = EXCLUDED.chat_name,
            updated_at = NOW()
    `

	_, err := d.db.Exec(query, chatID, chatName)
	if err != nil {
		return fmt.Errorf("failed to create/update chat: %w", err)
	}

	return nil
}

// UpdateChatName обновляет название чата
func (d *Db) UpdateChatName(chatID int64, newChatName string) error {
	query := `UPDATE telegram.chats SET chat_name = $1, updated_at = NOW() WHERE chat_id = $2`

	result, err := d.db.Exec(query, newChatName, chatID)
	if err != nil {
		return fmt.Errorf("failed to update chat name: %w", err)
	}

	r, _ := result.RowsAffected()
	if r == 0 {
		return fmt.Errorf("chat not found")
	}

	return nil
}

// DeleteChat удаляет чат и все связанные данные
func (d *Db) DeleteChat(chatID int64) error {
	query := `DELETE FROM telegram.chats WHERE chat_id = $1`

	result, err := d.db.Exec(query, chatID)
	if err != nil {
		return fmt.Errorf("failed to delete chat: %w", err)
	}

	r, _ := result.RowsAffected()
	if r == 0 {
		return fmt.Errorf("chat not found")
	}

	return nil
}

func (d *Db) DeleteChatFull(chatID int64) error {
	// Начинаем транзакцию
	tx, err := d.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// В случае ошибки откатываем изменения.
	// Если всё ок, Rollback после Commit ничего не сделает.
	defer tx.Rollback()

	// Список запросов на удаление.
	// Порядок: сначала зависимые данные, потом основная запись чата.
	queries := []struct {
		sql string
		arg any
	}{
		{`DELETE FROM telegram.chat_members WHERE chat_id = $1`, chatID},
		{`DELETE FROM telegram.chat_permissions WHERE chat_id = $1`, chatID},
		{`DELETE FROM telegram.user_roles WHERE chat_id = $1`, chatID},
		{`DELETE FROM telegram.roles WHERE chat_id = $1`, chatID},
		{`DELETE FROM telegram.chats WHERE chat_id = $1`, chatID},
	}

	for _, q := range queries {
		_, err := tx.Exec(q.sql, q.arg)
		if err != nil {
			return fmt.Errorf("failed to execute delete query [%s]: %w", q.sql, err)
		}
	}

	// Фиксируем изменения
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// GetUserChats возвращает список чатов пользователя
func (d *Db) GetUserChats(userID int64) ([]models2.Chat, error) {
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

	var chats []models2.Chat
	for rows.Next() {
		var chat models2.Chat
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

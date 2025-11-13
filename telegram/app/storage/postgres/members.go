package postgres

import (
	"encoding/json"
	"telegram/models"

	tgbotapi "github.com/OvyFlash/telegram-bot-api"
)

func (d *Db) UpsertChatData(m map[*models.Chat]map[int64]tgbotapi.User) {
	for chat, data := range m {
		// Конвертируем map в JSONB
		jsonData, err := json.Marshal(data)
		if err != nil {
			d.log.ErrorErr(err)
		}
		d.upsertChatData(*chat, jsonData)
	}
}
func (d *Db) upsertChatData(chat models.Chat, jsonData []byte) {
	ctx, cancel := d.getContext()
	defer cancel()

	// Используем INSERT ... ON CONFLICT для upsert
	query := `
        INSERT INTO rs_bot.chat_members(chat_id, chat_name, data) 
        VALUES ($1, $2, $3)
        ON CONFLICT (chat_id) 
        DO UPDATE SET chat_name = $2, data = $3, updated_at = CURRENT_TIMESTAMP
    `

	_, err := d.db.Exec(ctx, query, chat.ChatID, chat.ChatName, jsonData)
	if err != nil {
		d.log.ErrorErr(err)
	}
}
func (d *Db) upsertChatRoles(chat models.Chat, jsonData []byte) {
	ctx, cancel := d.getContext()
	defer cancel()

	// Используем INSERT ... ON CONFLICT для upsert
	query := `
        INSERT INTO rs_bot.chat_members(chat_id, chat_name, roles) 
        VALUES ($1, $2, $3)
        ON CONFLICT (chat_id) 
        DO UPDATE SET chat_name = $2 AND roles = $3, updated_at = CURRENT_TIMESTAMP
    `

	_, err := d.db.Exec(ctx, query, chat.ChatID, chat.ChatName, jsonData)
	if err != nil {
		d.log.ErrorErr(err)
	}
}

func (d *Db) ReadAllMembers() map[*models.Chat]map[int64]tgbotapi.User {
	chatMembersAll := make(map[*models.Chat]map[int64]tgbotapi.User)

	ctx, cancel := d.getContext()
	defer cancel()

	// Выполнение запроса
	rows, _ := d.db.Query(ctx, `SELECT chat_id, chat_name, data FROM rs_bot.chat_members;`)

	defer rows.Close()
	for rows.Next() {
		var chat models.Chat
		var data []byte
		var chatMembers map[int64]tgbotapi.User
		err := rows.Scan(&chat.ChatID, &chat.ChatName, &data)
		if err != nil {
			d.log.ErrorErr(err)
		}
		json.Unmarshal(data, &chatMembers)
		if len(chatMembers) != 0 {
			chatMembersAll[&chat] = chatMembers
		}
	}
	return chatMembersAll
}

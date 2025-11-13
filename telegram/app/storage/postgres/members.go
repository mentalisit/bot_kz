package postgres

import (
	"encoding/json"

	tgbotapi "github.com/OvyFlash/telegram-bot-api"
)

func (d *Db) UpsertChatData(m map[int64]map[int64]tgbotapi.User) {
	for chatId, data := range m {
		// Конвертируем map в JSONB
		jsonData, err := json.Marshal(data)
		if err != nil {
			d.log.ErrorErr(err)
		}
		d.upsertChatData(chatId, jsonData)
	}
}
func (d *Db) upsertChatData(chatId int64, jsonData []byte) {
	ctx, cancel := d.getContext()
	defer cancel()

	// Используем INSERT ... ON CONFLICT для upsert
	query := `
        INSERT INTO rs_bot.members(chat_id, data) 
        VALUES ($1, $2)
        ON CONFLICT (chat_id) 
        DO UPDATE SET data = $2, updated_at = CURRENT_TIMESTAMP
    `

	_, err := d.db.Exec(ctx, query, chatId, jsonData)
	if err != nil {
		d.log.ErrorErr(err)
	}
}
func (d *Db) ReadAllMembers() map[int64]map[int64]tgbotapi.User {
	chatMembersAll := make(map[int64]map[int64]tgbotapi.User)

	ctx, cancel := d.getContext()
	defer cancel()

	// Выполнение запроса
	rows, _ := d.db.Query(ctx, `SELECT chat_id, data FROM rs_bot.members;`)

	defer rows.Close()
	for rows.Next() {
		var chatId int64
		var data []byte
		var chatMembers map[int64]tgbotapi.User
		err := rows.Scan(&chatId, &data)
		if err != nil {
			d.log.ErrorErr(err)
		}
		json.Unmarshal(data, &chatMembers)
		if len(chatMembers) != 0 {
			chatMembersAll[chatId] = chatMembers
		}
	}
	return chatMembersAll
}

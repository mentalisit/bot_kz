package postgresV2

import (
	"errors"
	"rs/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// ReadNewsAll - чтение всех записей из таблицы news
func (d *Db) ReadNewsAll() []models.News {

	var news []models.News

	query := "SELECT uid, chat_id, language, messenger_type FROM rs_bot.news ORDER BY chat_id"
	results, err := d.db.Query(query)
	if err != nil {
		d.log.ErrorErr(err)
		return news
	}
	defer results.Close()

	for results.Next() {
		var n models.News
		var uid uuid.UUID
		err = results.Scan(&uid, &n.ChatId, &n.Language, &n.MessengerType)
		if err != nil {
			d.log.ErrorErr(err)
			continue
		}
		n.Uid = uid.String()
		news = append(news, n)
	}

	return news
}

// InsertNews - добавление новой записи
func (d *Db) InsertNews(chatId, language, messengerType string) {

	insert := `INSERT INTO rs_bot.news(chat_id, language, messenger_type) VALUES ($1, $2, $3) ON CONFLICT (chat_id) DO UPDATE SET language = $2, messenger_type = $3`
	_, err := d.db.Exec(insert, chatId, language, messengerType)
	if err != nil {
		d.log.ErrorErr(err)
	}
}

// DeleteNews - удаление записи по chat_id
func (d *Db) DeleteNews(chatId string) {

	deleteNews := `DELETE FROM rs_bot.news WHERE chat_id = $1`
	_, err := d.db.Exec(deleteNews, chatId)
	if err != nil {
		d.log.ErrorErr(err)
	}
}

// IsSubscribedToNews получает данные подписки на новости по chat_id
func (d *Db) IsSubscribedToNews(chatId string) (models.News, bool) {

	query := `SELECT uid, chat_id, language, messenger_type FROM rs_bot.news WHERE chat_id = $1`
	var n models.News
	var uid uuid.UUID
	err := d.db.QueryRow(query, chatId).Scan(&uid, &n.ChatId, &n.Language, &n.MessengerType)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return n, false
		}
		d.log.ErrorErr(err)
		return n, false
	}
	n.Uid = uid.String()
	return n, true
}

package postgres

import (
	"fmt"
	"rs/models"
)

func (d *Db) ConfigWebhookInsert(u models.ConfigWebhook) error {
	ctx, cancel := d.GetContext()
	defer cancel()
	insert := `INSERT INTO rs_bot.configwebhook(name, webhookgame, scorechannel) 
       VALUES ($1,$2,$3)`
	_, err := d.db.Exec(ctx, insert, u.Name, u.WebhookGame, u.ScoreChannel)
	if err != nil {
		return err
	}
	return nil
}

func (d *Db) ConfigWebhookGetAll() ([]models.ConfigWebhook, error) {
	ctx, cancel := d.GetContext()
	defer cancel()
	var cw []models.ConfigWebhook
	selectConfig := `SELECT name,webhookgame,scorechannel FROM rs_bot.configwebhook`

	results, err := d.db.Query(ctx, selectConfig)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса: %w", err)
	}
	defer results.Close()

	for results.Next() {
		var u models.ConfigWebhook
		err = results.Scan(&u.Name, &u.WebhookGame, &u.ScoreChannel)
		if err != nil {
			return nil, fmt.Errorf("ошибка сканирования строки: %w", err)
		}

		cw = append(cw, u)
	}

	// Проверка ошибок после цикла
	if err = results.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при обработке результатов: %w", err)
	}

	return cw, nil
}

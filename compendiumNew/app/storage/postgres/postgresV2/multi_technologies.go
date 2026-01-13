package postgresv2

import (
	"compendium/models"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

// Multi-technologies methods
func (d *Db) TechnologiesGet(uid uuid.UUID, name string) (*models.TechLevels, error) {
	var tech models.TechLevels
	query := `SELECT tech FROM my_compendium.technologies WHERE uid = $1 AND username = $2`

	// sqlx.Get сам вызовет .Scan() у нашего типа TechLevels
	err := d.db.Get(&tech, query, uid, name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &tech, nil
}

func (d *Db) TechnologiesUpdate(uid uuid.UUID, username string, tech models.TechLevels) error {
	query := `
       INSERT INTO my_compendium.technologies (uid, username, tech)
       VALUES ($1, $2, $3)
       ON CONFLICT (uid, username) DO UPDATE SET tech = EXCLUDED.tech`

	// Передаем tech напрямую, sqlx вызовет его метод .Value()
	_, err := d.db.Exec(query, uid, username, tech)
	return err
}

func (d *Db) TechnologiesDelete(uid uuid.UUID, username string) error {
	// Выполняем удаление сразу.
	// Если строки нет, Exec просто вернет 0 затронутых строк без ошибки.
	query := `DELETE FROM my_compendium.technologies WHERE uid = $1 AND username = $2`

	_, err := d.db.Exec(query, uid, username)
	if err != nil {
		d.log.ErrorErr(fmt.Errorf("failed to delete tech for %s: %w", username, err))
		return err
	}
	return nil
}
func (d *Db) TechnologiesUpdateUsername(uid uuid.UUID, oldUsername, newUsername string) error {
	query := `
        UPDATE my_compendium.technologies 
        SET username = $1 
        WHERE username = $2 AND uid = $3`

	result, err := d.db.Exec(query, newUsername, oldUsername, uid)
	if err != nil {
		return err
	}

	// Опционально: проверяем, была ли найдена и обновлена строка
	rows, _ := result.RowsAffected()
	if rows == 0 {
		d.log.Warn("No technology record found to update username" + " uid " + uid.String() + " old " + oldUsername)
	}

	return nil
}

func (d *Db) TechnologiesGetAllCorpMember(cm models.CorpMember) ([]models.CorpMember, error) {
	var results []models.CorpMember

	// ВАЖНО: Мы выбираем данные сразу в структуру.
	// Предполагается, что в CorpMember есть поля с тегами db:"username" и db:"tech"
	query := `SELECT username, tech FROM my_compendium.technologies WHERE uid = $1`

	// Select автоматически создаст слайс, проитерирует строки и распарсит JSON
	err := d.db.Select(&results, query, cm.MAcc.UUID)
	if err != nil {
		return nil, err
	}

	// Если вам нужно, чтобы в каждом объекте были данные из переданного cm (например, MAcc.UUID),
	// Select их не заполнит (он заполняет только то, что в SELECT).
	// Пройдемся циклом, если нужно дозаполнить общие поля:
	for i := range results {
		results[i].MAcc = cm.MAcc
	}

	return results, nil
}

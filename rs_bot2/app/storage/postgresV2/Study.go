package postgresV2

import (
	"database/sql"
	"errors"
	"fmt"
	"rs/models"
	"time"

	"github.com/google/uuid"
)

func (d *Db) GetStudy(uid uuid.UUID, name string) (*models.Study, error) {
	s := models.Study{
		Uuid:    uid,
		Name:    name,
		Studies: nil,
	}
	selectQuery := `SELECT studies FROM my_compendium.study where uid = $1 and name = $2`
	err := d.db.QueryRow(selectQuery, uid, name).Scan(&s.Studies)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		d.log.ErrorErr(err)
		return nil, err
	}
	return &s, nil
}

func (d *Db) InsertStudy(s models.Study) error {
	insertQuery := `
        INSERT INTO my_compendium.study (uid, name, studies)
        VALUES ($1, $2, $3)
        ON CONFLICT (uid, name) DO UPDATE 
        SET studies = EXCLUDED.studies`

	_, err := d.db.Exec(insertQuery, s.Uuid, s.Name, s.Studies)
	if err != nil {
		d.log.ErrorErr(err)
		return err
	}
	return nil
}

func (d *Db) GetAllExpiredStudy() ([]models.Study, error) {
	now := time.Now().UnixMilli()

	// Используем простой и эффективный запрос
	query := `
        SELECT uid, name, studies 
        FROM my_compendium.study 
        WHERE studies @> '[{"endTime": 0}]' -- это "подсказка" для GIN индекса, если он есть
           OR EXISTS (
               SELECT 1 FROM jsonb_array_elements(studies) AS elem 
               WHERE (elem->>'endTime')::bigint < $1
           )`

	rows, err := d.db.Query(query, now)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.Study
	for rows.Next() {
		var s models.Study
		if err := rows.Scan(&s.Uuid, &s.Name, &s.Studies); err != nil {
			return nil, err
		}
		results = append(results, s)
	}

	return results, nil
}

// Удаляет конкретные элементы из массива studies или обновляет запись целиком
func (d *Db) UpdateStudies(s models.Study) error {
	query := `UPDATE my_compendium.study SET studies = $1 WHERE uid = $2 and name = $3`
	_, err := d.db.Exec(query, s.Studies, s.Uuid, s.Name)
	return err
}

// DeleteStudyRecord полностью удаляет строку пользователя из таблицы
func (d *Db) DeleteStudyRecord(s models.Study) error {
	query := `DELETE FROM my_compendium.study WHERE uid = $1 and name = $2`
	_, err := d.db.Exec(query, s.Uuid, s.Name)
	if err != nil {
		return fmt.Errorf("failed to delete study record: %w", err)
	}
	return nil
}

func (d *Db) SyncModuleStatus(s models.Study, m models.Studies) error {
	// Подготавливаем фрагмент JSON для обновления одного модуля
	updateData := map[string]models.TechLevel{
		m.ModuleId: models.TechLevel{
			Level: m.Level,
			Ts:    m.EndTime,
		},
	}

	// Оператор || объединяет текущий JSON с новым, заменяя или добавляя ключи
	query := `
        UPDATE my_compendium.technologies 
        SET tech = COALESCE(tech, '{}'::jsonb) || $3 
        WHERE uid = $1 AND username = $2`

	_, err := d.db.Exec(query, s.Uuid, s.Name, updateData)
	if err != nil {
		d.log.ErrorErr(err)
		return err
	}
	return nil
}

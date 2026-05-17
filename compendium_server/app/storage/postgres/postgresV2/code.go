package postgresv2

import (
	"compendium_s/models"
	"database/sql"
	"errors"
)

func (d *Db) CodeInsert(u models.Code) error {
	// Если Identity реализует Value(), можно просто передать структуру u
	// при условии, что у Code стоят правильные db теги
	query := `INSERT INTO my_compendium.codes (code, time, identity) 
              VALUES (:code, :time, :identity)`

	_, err := d.db.NamedExec(query, u)
	return err
}

func (d *Db) CodeGet(code string) (*models.Code, error) {
	var u models.Code

	// sqlx сопоставит колонки с полями структуры по db-тегам.
	// Если в таблице есть колонка id, которой нет в структуре Code,
	// лучше перечислить нужные колонки явно.
	query := `SELECT code, time, identity FROM my_compendium.codes WHERE code = $1`

	// db.Get сразу записывает результат в структуру u.
	// Метод Scan() для Identity вызовется автоматически.
	err := d.db.Get(&u, query, code)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Запись не найдена
		}
		d.log.ErrorErr(err)
		return nil, err
	}

	return &u, nil
}

func (d *Db) CodeDelete(code string) {
	del := "delete from my_compendium.codes where code = $1"
	_, err := d.db.Exec(del, code)
	if err != nil {
		d.log.ErrorErr(err)
	}
}

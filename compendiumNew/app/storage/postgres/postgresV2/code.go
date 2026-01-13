package postgresv2

import (
	"compendium/models"
)

func (d *Db) CodeInsert(u models.Code) error {
	// Если Identity реализует Value(), можно просто передать структуру u
	// при условии, что у Code стоят правильные db теги
	query := `INSERT INTO my_compendium.codes (code, time, identity) 
              VALUES (:code, :time, :identity)`

	_, err := d.db.NamedExec(query, u)
	return err
}

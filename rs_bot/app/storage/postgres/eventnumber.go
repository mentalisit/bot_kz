package postgres

import "rs/models"

func (d *Db) InsertEventNumber(number, event int, status bool) error {
	ctx, cancel := d.GetContext()
	defer cancel()

	query := `INSERT INTO rs_bot.eventnumber (number, event, status) VALUES ($1, $2, $3)`
	_, err := d.db.Exec(ctx, query, number, event, status)
	if err != nil {
		d.log.ErrorErr(err)
		return err
	}
	return nil
}

func (d *Db) GetEventNumber() (events []models.Events, err error) {
	ctx, cancel := d.GetContext()
	defer cancel()

	query := `SELECT id, number, event, status FROM rs_bot.eventnumber WHERE status = false`
	rows, err := d.db.Query(ctx, query)
	if err != nil {
		d.log.ErrorErr(err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var e models.Events
		if err := rows.Scan(&e.ID, &e.Number, &e.Event, &e.Status); err != nil {
			d.log.ErrorErr(err)
			return nil, err
		}
		events = append(events, e)
	}

	if err := rows.Err(); err != nil {
		d.log.ErrorErr(err)
		return nil, err
	}

	return events, nil
}

func (d *Db) UpdateEventStatus(id int64, newStatus bool) error {
	ctx, cancel := d.GetContext()
	defer cancel()

	query := `UPDATE rs_bot.eventnumber SET status = $1 WHERE id = $2`
	_, err := d.db.Exec(ctx, query, newStatus, id)
	if err != nil {
		d.log.ErrorErr(err)
		return err
	}
	return nil
}

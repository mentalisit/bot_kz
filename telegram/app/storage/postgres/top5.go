package postgres

func (d *Db) ReadTop5Level(corpname string) []string {
	ctx, cancel := d.GetContext()
	defer cancel()
	query := `
        SELECT lvlkz, COUNT(*) AS lvlkz_count
        FROM kzbot.sborkz
        WHERE corpname=$1
          AND date::timestamp >= CURRENT_DATE - INTERVAL '5 days'
        GROUP BY lvlkz
        ORDER BY lvlkz_count DESC
        LIMIT 5;
    `

	// Выполнение запроса
	rows, err := d.db.Query(ctx, query, corpname)
	defer rows.Close()
	if err != nil {
		d.log.ErrorErr(err)
	}

	var levels []string

	// Итерация по результатам запроса
	for rows.Next() {
		var lvlkz string
		var lvlkzCount int
		if err = rows.Scan(&lvlkz, &lvlkzCount); err != nil {
			d.log.ErrorErr(err)
		}
		levels = append(levels, lvlkz)
	}
	if err = rows.Err(); err != nil {
		d.log.ErrorErr(err)
	}
	return levels
}

func (d *Db) ReadTelegramLastMessage(corpname string) int {
	ctx, cancel := d.GetContext()
	defer cancel()
	query := `
        SELECT MAX(tgmesid) FROM kzbot.sborkz
        WHERE corpname=$1;
    `

	var mid int

	// Выполнение запроса
	err := d.db.QueryRow(ctx, query, corpname).Scan(&mid)

	if err != nil {
		d.log.ErrorErr(err)
		return 0
	}
	return mid
}

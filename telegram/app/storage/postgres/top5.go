package postgres

func (d *Db) ReadTop5Level(corpname string) []string {
	ctx, cancel := d.getContext()
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
	if err != nil {
		d.log.ErrorErr(err)
	}
	defer rows.Close()

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
	ctx, cancel := d.getContext()
	defer cancel()
	query := `
        SELECT MAX(tgmesid) FROM kzbot.sborkz
        WHERE corpname=$1 AND active > 0;
    `

	var mid int

	// Выполнение запроса
	err := d.db.QueryRow(ctx, query, corpname).Scan(&mid)

	if err != nil {
		//fmt.Printf("ReadTelegramLastMessage corp:%s err %+v\n", corpname, err)
		return 0
	}
	return mid
}

func (d *Db) ReadTop5LevelForV2(corpUid string) []string {
	ctx, cancel := d.getContext()
	defer cancel()

	query := `
        SELECT 
            data->'Data'->>'lvl_rs' AS lvl,
            COUNT(*) AS lvl_count
        FROM rs_bot.queue_complete
        WHERE data->'Data'->>'corporation_uuid' = $1
          AND (data->'Data'->>'date')::timestamp >= CURRENT_DATE - INTERVAL '5 days'
        GROUP BY lvl
        ORDER BY lvl_count DESC
        LIMIT 5;
    `

	rows, err := d.db.Query(ctx, query, corpUid)
	if err != nil {
		d.log.ErrorErr(err)
		return nil
	}
	defer rows.Close()

	var levels []string

	for rows.Next() {
		var lvl string
		var count int

		if err := rows.Scan(&lvl, &count); err != nil {
			d.log.ErrorErr(err)
			continue
		}

		levels = append(levels, lvl)
	}

	if err := rows.Err(); err != nil {
		d.log.ErrorErr(err)
	}

	return levels
}

func (d *Db) ReadTelegramLastMessageV2(corpUUID string, chatID string) int {
	ctx, cancel := d.getContext()
	defer cancel()

	query := `
        SELECT COALESCE(MAX((msg.value->>'message_id')::bigint), 0) AS max_mid
        FROM rs_bot.queue_complete,
             jsonb_each(data->'Messages') AS msg(key, value)
        WHERE data->'Data'->>'corporation_uuid' = $1
          AND msg.value->>'type_messenger' = 'tg'
          AND msg.key = $2;
    `

	var mid int

	err := d.db.QueryRow(ctx, query, corpUUID, chatID).Scan(&mid)
	if err != nil {
		return 0
	}

	return mid
}

package postgresV2

import "fmt"

// SavePollChannel добавляет канал опросов в JSON поле data таблицы my_compendium.guilds
func (d *Db) SavePollChannel(gid, channelId, channelName, messengerType string) error {

	// Создаем JSON объект для нового канала
	newChannelJSON := fmt.Sprintf(`{"channel_id": "%s", "channel_name": "%s", "type_messenger": "%s"}`,
		channelId, channelName, messengerType)

	// Сначала заменяем null на пустые массивы, чтобы избежать ошибок
	initQuery := `
		UPDATE my_compendium.guilds 
		SET data = COALESCE(data, '{}'::jsonb) || '{"poll_channels": []}'::jsonb
		WHERE gid = $1 AND (data->'poll_channels') IS NULL
	`
	_, initErr := d.db.Exec(initQuery, gid)
	if initErr != nil {
		d.log.ErrorErr(initErr)
		return initErr
	}

	// Теперь проверяем, есть ли уже такой канал
	checkQuery := `
		SELECT COUNT(*) FROM my_compendium.guilds 
		WHERE gid = $1 AND data IS NOT NULL 
		AND jsonb_typeof(data) = 'object' 
		AND data->'poll_channels' IS NOT NULL 
		AND jsonb_typeof(data->'poll_channels') = 'array'
		AND EXISTS (
			SELECT 1 FROM jsonb_array_elements(data->'poll_channels') elem
			WHERE (elem->>'channel_id' = $2 AND elem->>'type_messenger' = $3)
		)`

	var count int
	err := d.db.QueryRow(checkQuery, gid, channelId, messengerType).Scan(&count)
	if err != nil {
		d.log.ErrorErr(err)
		return err
	}

	if count > 0 {
		// Канал существует, обновляем его имя
		updateQuery := `
			UPDATE my_compendium.guilds 
			SET data = jsonb_set(
				data, 
				'{poll_channels}', 
				(
					SELECT jsonb_agg(
						CASE 
							WHEN (elem->>'channel_id' = $2 AND elem->>'type_messenger' = $3) 
							THEN jsonb_set(elem, '{channel_name}', to_jsonb($4::text))
							ELSE elem 
						END
					)
					FROM jsonb_array_elements(data->'poll_channels') elem
				)
			)
			WHERE gid = $1 AND data->'poll_channels' IS NOT NULL AND jsonb_typeof(data->'poll_channels') = 'array'
		`
		_, err = d.db.Exec(updateQuery, gid, channelId, messengerType, channelName)
		if err != nil {
			d.log.ErrorErr(err)
		}
	} else {
		// Канал не существует, добавляем его
		d.log.Info(fmt.Sprintf("SavePollChannel: adding new channel: %s", newChannelJSON))
		updateQuery := `
			UPDATE my_compendium.guilds 
			SET data = jsonb_set(
				COALESCE(data, '{}'::jsonb), 
				'{poll_channels}', 
				COALESCE(data->'poll_channels', '[]'::jsonb) || $2::jsonb
			)
			WHERE gid = $1
		`
		_, err = d.db.Exec(updateQuery, gid, newChannelJSON)
		if err != nil {
			d.log.ErrorErr(err)
		}
	}

	if err != nil {
		d.log.ErrorErr(err)
	}
	return err
}

// DeletePollChannel удаляет канал опросов из JSON поля data таблицы my_compendium.guilds
func (d *Db) DeletePollChannel(gid, channelId, messengerType string) error {

	// Удаляем канал из JSON массива с помощью jsonb_array_elements
	updateQuery := `
		UPDATE my_compendium.guilds 
		SET data = jsonb_set(
			data, 
			'{poll_channels}', 
			(
				SELECT jsonb_agg(elem)
				FROM jsonb_array_elements(data->'poll_channels') elem
				WHERE NOT (elem->>'channel_id' = $2 AND elem->>'type_messenger' = $3)
			)
		)
		WHERE gid = $1 AND data->'poll_channels' IS NOT NULL AND jsonb_typeof(data->'poll_channels') = 'array'
	`

	_, err := d.db.Exec(updateQuery, gid, channelId, messengerType)
	if err != nil {
		d.log.ErrorErr(err)
	}
	return err
}

// SaveBzChannel добавляет канал БЗ в JSON поле data таблицы my_compendium.guilds
func (d *Db) SaveBzChannel(gid, bzKey, channelType, channelId, channelName, messengerType string) error {

	// Создаем JSON объект для нового канала
	newChannelJSON := fmt.Sprintf(`{"channel_id": "%s", "channel_name": "%s", "type_messenger": "%s"}`,
		channelId, channelName, messengerType)

	// Определяем путь в JSON
	var jsonPath string
	switch channelType {
	case "coordination":
		jsonPath = fmt.Sprintf("coordination.%s", bzKey)
	case "discussion":
		jsonPath = fmt.Sprintf("discussion.%s", bzKey)
	default:
		return fmt.Errorf("неизвестный тип канала: %s", channelType)
	}

	// Проверяем, есть ли уже такой канал в массиве
	checkQuery := `
		SELECT COUNT(*) FROM my_compendium.guilds 
		WHERE gid = $1 AND data IS NOT NULL 
		AND jsonb_typeof(data) = 'object' 
		AND data->%s IS NOT NULL 
		AND jsonb_typeof(data->%s) = 'array'
		AND EXISTS (
			SELECT 1 FROM jsonb_array_elements(data->%s) elem
			WHERE (elem->>'channel_id' = $2 AND elem->>'type_messenger' = $3)
		)
	`

	var count int
	err := d.db.QueryRow(fmt.Sprintf(checkQuery, jsonPath, jsonPath, jsonPath), gid, channelId, messengerType).Scan(&count)
	if err != nil {
		d.log.ErrorErr(err)
		return err
	}

	if count > 0 {
		// Канал существует, обновляем его имя
		updateQuery := `
			UPDATE my_compendium.guilds 
			SET data = jsonb_set(
				data, 
				ARRAY[%s]::text[], 
				(
					SELECT jsonb_agg(
						CASE 
							WHEN (elem->>'channel_id' = $2 AND elem->>'type_messenger' = $3) 
							THEN jsonb_set(elem, '{channel_name}', to_jsonb($4::text))
							ELSE elem 
						END
					)
					FROM jsonb_array_elements(data->%s) elem
				)
			)
			WHERE gid = $1 AND data->%s IS NOT NULL AND jsonb_typeof(data->%s) = 'array'
		`
		pathArray := fmt.Sprintf("'%s', '%s'", channelType, bzKey)
		_, err = d.db.Exec(fmt.Sprintf(updateQuery, pathArray, jsonPath), gid, channelId, messengerType, channelName)
	} else {
		// Канал не существует, добавляем его
		updateQuery := `
			UPDATE my_compendium.guilds 
			SET data = jsonb_set(
				COALESCE(data, '{}'::jsonb), 
				ARRAY[%s]::text[], 
				COALESCE(data->%s, '[]'::jsonb) || $2::jsonb
			)
			WHERE gid = $1
		`
		pathArray := fmt.Sprintf("'%s', '%s'", channelType, bzKey)
		_, _ = d.db.Exec(fmt.Sprintf(updateQuery, pathArray, jsonPath), gid, newChannelJSON)
	}

	if err != nil {
		d.log.ErrorErr(err)
	}
	return err
}

// DeleteBzChannel удаляет канал БЗ из JSON поля data таблицы my_compendium.guilds
func (d *Db) DeleteBzChannel(gid, bzKey, channelType, channelId, messengerType string) error {

	// Определяем путь в JSON
	var jsonPath string
	switch channelType {
	case "coordination":
		jsonPath = fmt.Sprintf("coordination.%s", bzKey)
	case "discussion":
		jsonPath = fmt.Sprintf("discussion.%s", bzKey)
	default:
		return fmt.Errorf("неизвестный тип канала: %s", channelType)
	}

	// Удаляем канал из JSON массива с помощью jsonb_array_elements
	updateQuery := `
		UPDATE my_compendium.guilds 
		SET data = jsonb_set(
			COALESCE(data, '{}'::jsonb), 
			ARRAY[%s]::text[], 
			(
				SELECT jsonb_agg(elem)
				FROM jsonb_array_elements(COALESCE(data->%s, '[]'::jsonb)) elem
				WHERE NOT (elem->>'channel_id' = $2 AND elem->>'type_messenger' = $3)
			)
		)
		WHERE gid = $1
	`

	pathArray := fmt.Sprintf("'%s', '%s'", channelType, bzKey)
	_, err := d.db.Exec(fmt.Sprintf(updateQuery, pathArray, jsonPath), gid, channelId, messengerType)
	if err != nil {
		d.log.ErrorErr(err)
	}
	return err
}

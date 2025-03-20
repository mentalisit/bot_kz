package dbpostgres

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"ws/models"
)

// Загружаем все данные в кэш
func (d *Db) LoadAllData() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	rows, err := d.pool.Query(ctx, `SELECT id, data FROM ws.corporations`)
	if err != nil {
		d.log.ErrorErr(err)
		return
	}
	defer rows.Close()

	d.mu.Lock()
	defer d.mu.Unlock()

	d.cache = make(map[string]models.CorporationsData) // Очищаем старый кэш

	for rows.Next() {
		var key string
		var jsonData string
		var corp models.CorporationsData

		if err := rows.Scan(&key, &jsonData); err != nil {
			d.log.ErrorErr(err)
			continue
		}

		if err := json.Unmarshal([]byte(jsonData), &corp); err != nil {
			d.log.ErrorErr(err)
			continue
		}

		d.cache[key] = corp
	}

	fmt.Println("Cache loaded:", len(d.cache), "records")
}

// Чтение данных (сначала из кэша, потом из БД)
func (d *Db) ReadCorpData(key string) *models.CorporationsData {
	d.mu.RLock()
	data, found := d.cache[key]
	d.mu.RUnlock()

	if found {
		return &data
	}

	// Если в кэше нет, ищем в БД
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var jsonData string
	query := `SELECT data FROM ws.corporations WHERE id = $1`
	err := d.pool.QueryRow(ctx, query, key).Scan(&jsonData)
	if err != nil {
		return nil // Не нашли
	}

	var corp models.CorporationsData
	if err := json.Unmarshal([]byte(jsonData), &corp); err != nil {
		d.log.ErrorErr(err)
		return nil
	}

	// Добавляем в кэш
	d.mu.Lock()
	d.cache[key] = corp
	d.mu.Unlock()

	return &corp
}

// Сохранение данных
func (d *Db) SaveCorpData(key string, corp models.CorporationsData) {
	//mapCorp := map[string]interface{}{
	//	"Corporation1Name":  corp.Corporation1Name,
	//	"Corporation2Name":  corp.Corporation2Name,
	//	"Corporation1Score": corp.Corporation1Score,
	//	"Corporation2Score": corp.Corporation2Score,
	//	"Corporation1Id":    corp.Corporation1Id,
	//	"Corporation2Id":    corp.Corporation2Id,
	//	"DateEnded":         corp.DateEnded.Format(time.RFC3339), // Сохраняем время как строку
	//}
	jsonData, err := json.Marshal(corp)
	if err != nil {
		d.log.ErrorErr(err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `INSERT INTO ws.corporations (id, data)
			  VALUES ($1, $2)
			  ON CONFLICT (id) DO NOTHING`
	_, err = d.pool.Exec(ctx, query, key, jsonData)
	if err != nil {
		d.log.ErrorErr(err)
		return
	}

	// Добавляем в кэш
	d.mu.Lock()
	d.cache[key] = corp
	d.mu.Unlock()

	fmt.Printf("Saved %s Data %+v\n", key, corp)
}

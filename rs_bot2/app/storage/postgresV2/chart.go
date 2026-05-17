package postgresV2

import (
	"fmt"
	"rs/models"
)

// GetChartDataByCorpAndPeriod возвращает данные для графиков активности по корпорации.
// period: "week", "month", "quarter", "year", "all"
// level: конкретный уровень ("rs5", "drs10", ...) или "" — все уровни.
func (d *Db) GetChartDataByCorpAndPeriod(corpName, period, level string) (*models.ChartData, error) {

	interval := periodToInterval(period)
	selectDate := periodToSelectDate(period)
	groupDate := periodToGroupDate(period)

	// --- Серии по уровням ---
	seriesQuery := `
        SELECT 
            lvlkz,
            ` + selectDate + ` AS dt,
            COALESCE(SUM(active), 0) AS cnt
        FROM kzbot.sborkz
        WHERE corpname = $1
          AND active > 0
          AND date IS NOT NULL AND date != ''
          ` + interval + `
        GROUP BY lvlkz, ` + groupDate + `
        ORDER BY dt, lvlkz`

	args := []any{corpName}
	if level != "" {
		seriesQuery = `
        SELECT 
            lvlkz,
            ` + selectDate + ` AS dt,
            COALESCE(SUM(active), 0) AS cnt
        FROM kzbot.sborkz
        WHERE corpname = $1
          AND active > 0
          AND date IS NOT NULL AND date != ''
          AND (lvlkz = $2 OR lvlkz = 'd' || $2)
          ` + interval + `
        GROUP BY lvlkz, ` + groupDate + `
        ORDER BY dt, lvlkz`
		args = append(args, level)
	}

	rows, err := d.db.Query(seriesQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("chart series query: %w", err)
	}
	defer rows.Close()

	seriesMap := make(map[string][]models.ChartPoint)
	for rows.Next() {
		var lvl, dt string
		var cnt int
		if err := rows.Scan(&lvl, &dt, &cnt); err != nil {
			return nil, fmt.Errorf("chart series scan: %w", err)
		}
		seriesMap[lvl] = append(seriesMap[lvl], models.ChartPoint{Date: dt, Count: cnt})
	}

	var series []models.ChartSeries
	for lvl, points := range seriesMap {
		series = append(series, models.ChartSeries{Level: lvl, Points: points})
	}

	// --- Итого по всем уровням ---
	totalQuery := `
        SELECT 
            ` + selectDate + ` AS dt,
            COALESCE(SUM(active), 0) AS cnt
        FROM kzbot.sborkz
        WHERE corpname = $1
          AND active > 0
          AND date IS NOT NULL AND date != ''
          ` + interval + `
        GROUP BY ` + groupDate + `
        ORDER BY dt`

	totalRows, err := d.db.Query(totalQuery, corpName)
	if err != nil {
		return nil, fmt.Errorf("chart total query: %w", err)
	}
	defer totalRows.Close()

	var total []models.ChartPoint
	for totalRows.Next() {
		var pt models.ChartPoint
		if err := totalRows.Scan(&pt.Date, &pt.Count); err != nil {
			return nil, fmt.Errorf("chart total scan: %w", err)
		}
		total = append(total, pt)
	}

	return &models.ChartData{
		Period: period,
		Series: series,
		Total:  total,
	}, nil
}

// GetChartDataByUser возвращает данные графиков для конкретного игрока.
func (d *Db) GetChartDataByUser(name, period, level string) (*models.ChartData, error) {

	interval := periodToInterval(period)
	selectDate := periodToSelectDate(period)
	groupDate := periodToGroupDate(period)

	// --- Серии по уровням ---
	seriesQuery := `
        SELECT 
            lvlkz,
            ` + selectDate + ` AS dt,
            COALESCE(SUM(active), 0) AS cnt
        FROM kzbot.sborkz
        WHERE name = $1
          AND active > 0
          AND date IS NOT NULL AND date != ''
          ` + interval

	args := []any{name}
	if level != "" {
		seriesQuery += ` AND (lvlkz = $2 OR lvlkz = 'd' || $2)`
		args = append(args, level)
	}
	seriesQuery += `
        GROUP BY lvlkz, ` + groupDate + `
        ORDER BY dt, lvlkz`

	rows, err := d.db.Query(seriesQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("chart user series query: %w", err)
	}
	defer rows.Close()

	seriesMap := make(map[string][]models.ChartPoint)
	for rows.Next() {
		var lvl, dt string
		var cnt int
		if err := rows.Scan(&lvl, &dt, &cnt); err != nil {
			return nil, fmt.Errorf("chart user series scan: %w", err)
		}
		seriesMap[lvl] = append(seriesMap[lvl], models.ChartPoint{Date: dt, Count: cnt})
	}

	var series []models.ChartSeries
	for lvl, points := range seriesMap {
		series = append(series, models.ChartSeries{Level: lvl, Points: points})
	}

	// --- Итого ---
	totalQuery := `
        SELECT 
            ` + selectDate + ` AS dt,
            COALESCE(SUM(active), 0) AS cnt
        FROM kzbot.sborkz
        WHERE name = $1
          AND active > 0
          AND date IS NOT NULL AND date != ''
          ` + interval + `
        GROUP BY ` + groupDate + `
        ORDER BY dt`

	totalRows, err := d.db.Query(totalQuery, name)
	if err != nil {
		return nil, fmt.Errorf("chart user total query: %w", err)
	}
	defer totalRows.Close()

	var total []models.ChartPoint
	for totalRows.Next() {
		var pt models.ChartPoint
		if err := totalRows.Scan(&pt.Date, &pt.Count); err != nil {
			return nil, fmt.Errorf("chart user total scan: %w", err)
		}
		total = append(total, pt)
	}

	return &models.ChartData{
		Period: period,
		Series: series,
		Total:  total,
	}, nil
}

// GetAvailableLevels возвращает список уникальных уровней для корпорации.
func (d *Db) GetAvailableLevels(corpName string) ([]string, error) {

	query := `
        SELECT DISTINCT lvlkz 
        FROM kzbot.sborkz 
        WHERE corpname = $1 AND active > 0 AND lvlkz IS NOT NULL
        ORDER BY lvlkz`

	rows, err := d.db.Query(query, corpName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var levels []string
	for rows.Next() {
		var lvl string
		if err := rows.Scan(&lvl); err != nil {
			continue
		}
		levels = append(levels, lvl)
	}
	return levels, nil
}

// GetChartPopularityByCorp возвращает сводный график популярности по часам за выбранный период.
// Суммирует активность за каждый час (00-23) в течение всего периода (наложение дней).
func (d *Db) GetChartPopularityByCorp(corpName, period string) ([]models.ChartPoint, error) {

	interval := periodToInterval(period)

	query := `
        SELECT 
            SUBSTRING(time FROM 1 FOR 2) AS hour_str,
            COALESCE(SUM(active), 0) AS cnt
        FROM kzbot.sborkz
        WHERE corpname = $1
          AND active > 0
          AND time IS NOT NULL AND time != ''
          ` + interval + `
        GROUP BY hour_str
        ORDER BY hour_str`

	rows, err := d.db.Query(query, corpName)
	if err != nil {
		return nil, fmt.Errorf("popularity query error: %w", err)
	}
	defer rows.Close()

	// Инициализируем массив 24 часами с нулями, чтобы не было пропусков
	hoursMap := make(map[string]int)
	for i := 0; i < 24; i++ {
		hh := fmt.Sprintf("%02d", i)
		hoursMap[hh] = 0
	}

	for rows.Next() {
		var h string
		var cnt int
		if err := rows.Scan(&h, &cnt); err != nil {
			return nil, fmt.Errorf("popularity scan error: %w", err)
		}
		// Защита от криво сохраненных данных времени
		if len(h) >= 2 {
			hoursMap[h[:2]] = cnt
		}
	}

	// Собираем в отсортированный массив (00, 01, 02 ... 23)
	var result []models.ChartPoint
	for i := 0; i < 24; i++ {
		hh := fmt.Sprintf("%02d", i)
		result = append(result, models.ChartPoint{
			Date:  hh + ":00", // "00:00", "01:00"
			Count: hoursMap[hh],
		})
	}

	return result, nil
}

// GetAvailableCorps возвращает список уникальных корпораций, у которых есть статистика.
func (d *Db) GetAvailableCorps() ([]string, error) {

	query := `
        SELECT DISTINCT corpname 
        FROM kzbot.sborkz 
        WHERE active > 0 AND corpname IS NOT NULL AND corpname != ''
        ORDER BY corpname`

	rows, err := d.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var corps []string
	for rows.Next() {
		var c string
		if err := rows.Scan(&c); err != nil {
			continue
		}
		corps = append(corps, c)
	}
	return corps, nil
}

// periodToInterval возвращает SQL-фрагмент фильтра по дате.
func periodToInterval(period string) string {
	switch period {
	case "day":
		return "AND to_timestamp(date || ' ' || time, 'YYYY-MM-DD HH24:MI') >= NOW() - INTERVAL '24 hours'"
	case "week":
		return "AND to_timestamp(date,'YYYY-MM-DD') >= CURRENT_DATE - INTERVAL '7 days'"
	case "month":
		return "AND to_timestamp(date,'YYYY-MM-DD') >= CURRENT_DATE - INTERVAL '30 days'"
	case "quarter":
		return "AND to_timestamp(date,'YYYY-MM-DD') >= CURRENT_DATE - INTERVAL '90 days'"
	case "year":
		return "AND to_timestamp(date,'YYYY-MM-DD') >= CURRENT_DATE - INTERVAL '365 days'"
	default: // "all"
		return ""
	}
}

// periodToSelectDate возвращает SQL выражение для выборки даты/времени в зависимости от периода
func periodToSelectDate(period string) string {
	if period == "day" {
		return "date || 'T' || SUBSTRING(time FROM 1 FOR 2) || ':00:00'"
	}
	return "date"
}

// periodToGroupDate возвращает SQL выражение для группировки по дате/времени
func periodToGroupDate(period string) string {
	if period == "day" {
		return "date, SUBSTRING(time FROM 1 FOR 2)"
	}
	return "date"
}

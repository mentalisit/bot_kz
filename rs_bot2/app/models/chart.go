package models

// ChartPoint — одна точка на временном графике (дата + количество).
type ChartPoint struct {
	Date  string `json:"date"`  // "2026-03-15"
	Count int    `json:"count"` // количество за этот день
}

// ChartSeries — серия данных для одного уровня КЗ.
type ChartSeries struct {
	Level  string       `json:"level"`  // "rs5", "drs10", etc.
	Points []ChartPoint `json:"points"` // массив точек
}

// ChartData — полный набор данных для графиков.
type ChartData struct {
	Period string        `json:"period"` // "week", "month", "quarter", "year", "all"
	Series []ChartSeries `json:"series"` // данные по каждому уровню
	Total  []ChartPoint  `json:"total"`  // общее количество за период по всем уровням
}

package postgresv2

import (
	"database/sql"
)

// Helper functions
func scanStringArray(rows *sql.Rows, colName string) ([]string, error) {
	var result []string
	for rows.Next() {
		var item sql.NullString
		if err := rows.Scan(&item); err != nil {
			return nil, err
		}
		if item.Valid {
			result = append(result, item.String)
		}
	}
	return result, nil
}

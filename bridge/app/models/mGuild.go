package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

// MultiAccountGuildV2 - версия для V2 с Channels как map[string][]string
type MultiAccountGuildV2 struct {
	GId       uuid.UUID
	GuildName string
	Channels  GuildChannels `db:"channels"` // Наш новый тип
	AvatarUrl string
}
type GuildChannels map[string][]string

// Value преобразует map в JSON для базы данных
func (m GuildChannels) Value() (driver.Value, error) {
	if m == nil {
		return json.Marshal(map[string][]string{})
	}
	return json.Marshal(m)
}

// Scan преобразует JSON из базы данных в map
func (m *GuildChannels) Scan(src interface{}) error {
	if src == nil {
		return nil
	}
	bytes, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, m)
}

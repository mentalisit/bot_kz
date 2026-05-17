package models

import (
	"database/sql/driver"
	"encoding/json"
)

type LocalGuild struct {
	//Gid internal id
	Gid int
	//Name Corporation game name
	Name string
	Data LocalGuildData
}

type LocalGuildData struct {
	GuildIds []LocalGuildsIds
}

type LocalGuildsIds struct {
	TypeMessenger string
	GuildId       string
}

type MultiAccountData struct {
	Merged   []GameAccountData `json:"Merged"`
	Timezone string            `json:"timezone,omitempty"`
	WsPhone  string            `json:"wsPhone,omitempty"`
	NotifyPM bool              `json:"notifyPM,omitempty"`
	// Сюда можно добавлять любые другие настройки в будущем
}

// Scan implements sql.Scanner interface
func (m *MultiAccountData) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	var data []byte
	switch v := value.(type) {
	case []byte:
		data = v
	case string:
		data = []byte(v)
	default:
		return nil
	}
	return json.Unmarshal(data, m)
}

// Value implements driver.Valuer interface
func (m MultiAccountData) Value() (driver.Value, error) {
	return json.Marshal(m)
}

type GameAccountData struct {
	PlayerID     string `json:"id"`
	PlayerName   string `json:"name"`
	OwnerUuid    string `json:"ownerUuid"`
	CurrentOwner string `json:"currentOwner"`
}

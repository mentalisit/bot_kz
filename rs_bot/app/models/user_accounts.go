package models

import (
	"encoding/json"
	"strconv"
)

type UserAccount struct {
	InternalId  int
	GeneralName string
	TgId        string
	DsId        string
	GameId      []string
	ActiveName  string
	Accounts    []string
}

// GetAlt функция получения по номеру
// 0 возвращает GeneralName
// 1-5, возвращает альтов 0-4
// если i больше чем количество аккаунтов возвращаем GeneralName
func (u *UserAccount) GetAlt(i int) string {
	if i == 0 || i >= len(u.Accounts) {
		return u.GeneralName
	}
	return u.Accounts[i-1]
}
func (u *UserAccount) ContainsGameId(i int64) bool {
	s := strconv.FormatInt(i, 10)
	if len(u.GameId) == 0 {
		return false
	}
	for _, ss := range u.GameId {
		if ss == s {
			return true
		}
	}
	return false
}

type PlayerStats struct {
	Player string
	Points int
	Runs   int
	Level  int
}

type Statistic struct {
	EventId int // Номер ивента
	Level   int //Уровень кз
	Points  int //Всего очков
	Runs    int //Количество игр
}

type BattleStats struct {
	Name          string  `json:"name"`
	Level         string  `json:"level"`
	PointsSum     int     `json:"points_sum"`
	RecordsCount  int     `json:"records_count"`
	AveragePoints float64 `json:"average_points"`
	Quality       float64 `json:"quality"`
}

type Participants struct {
	PlayerID   string `json:"PlayerID"`
	PlayerName string `json:"PlayerName"`
}

type ParticipantsInt64 struct {
	PlayerID   int64  `json:"PlayerID"`
	PlayerName string `json:"PlayerName"`
}
type RedStarEvent struct {
	StarSystemID  string         `json:"StarSystemID"`
	StarLevel     int            `json:"StarLevel"`
	DarkRedStar   bool           `json:"DarkRedStar"`
	EventType     string         `json:"EventType"`
	Timestamp     string         `json:"Timestamp"`
	RSEventPoints int            `json:"RSEventPoints,omitempty"`
	Players       []Participants // Это поле будет сериализоваться как Players или PlayersWhoContributed
}

func (r *RedStarEvent) UnmarshalJSON(data []byte) error {
	type Alias RedStarEvent
	aux := &struct {
		*Alias
		PlayersField []Participants `json:"Players,omitempty"`
		ContribField []Participants `json:"PlayersWhoContributed,omitempty"`
	}{
		Alias: (*Alias)(r),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Определяем, какое поле использовать
	if len(aux.PlayersField) > 0 {
		r.Players = aux.PlayersField
	} else if len(aux.ContribField) > 0 {
		r.Players = aux.ContribField
	}

	return nil
}

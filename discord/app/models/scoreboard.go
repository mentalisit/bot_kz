package models

import (
	"encoding/json"
	"fmt"
)

type ScoreboardParams struct {
	Name              string
	ChannelWebhook    string
	ChannelScoreboard string
	LastMessageID     string
}

type RedStarEvent struct {
	StarSystemID  string        `json:"StarSystemID"`
	StarLevel     int           `json:"StarLevel"`
	DarkRedStar   bool          `json:"DarkRedStar"`
	EventType     string        `json:"EventType"`
	Corporation   Corporation   `json:"Corporation"`
	Timestamp     string        `json:"Timestamp"`
	RSEventPoints int           `json:"RSEventPoints,omitempty"`
	Players       []Participant // Это поле будет сериализоваться как Players или PlayersWhoContributed
}

func (r *RedStarEvent) UnmarshalJSON(data []byte) error {
	type Alias RedStarEvent
	aux := &struct {
		*Alias
		PlayersField []Participant `json:"Players,omitempty"`
		ContribField []Participant `json:"PlayersWhoContributed,omitempty"`
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

type WhiteStarStarted struct {
	Slot                 int           `json:"Slot"`
	Opponent             Corporation   `json:"Opponent"`
	EventType            string        `json:"EventType"`
	Timestamp            string        `json:"Timestamp"`
	IsUnderdog           bool          `json:"IsUnderdog"`
	Corporation          Corporation   `json:"Corporation"`
	WhiteStarID          string        `json:"WhiteStarID"`
	OurParticipants      []Participant `json:"OurParticipants"`
	OpponentParticipants []Participant `json:"OpponentParticipants"`
}
type WhiteStarEnded struct {
	WhiteStarID   string      `json:"WhiteStarID"`
	Opponent      Corporation `json:"Opponent"`
	OurScore      int         `json:"OurScore"`
	OpponentScore int         `json:"OpponentScore"`
	XPGained      int         `json:"XPGained"`
	Corporation   Corporation `json:"Corporation"`
	EventType     string      `json:"EventType"`
	Timestamp     string      `json:"Timestamp"`
}

type Corporation struct {
	CorporationID   string `json:"CorporationID"`
	CorporationName string `json:"CorporationName"`
}

func (c *Corporation) GetAvatar() string {
	if c.CorporationID != "" {
		return fmt.Sprintf("https://ws.tsl.rocks/corp/%s/favicon.png", c.CorporationID)
	}
	return ""
}

type Participant struct {
	PlayerID   string `json:"PlayerID"`
	PlayerName string `json:"PlayerName"`
}

type Battles struct {
	EventId  int
	CorpName string
	Name     string
	Level    int
	Points   int
}
type BattlesTop struct {
	Id       int
	CorpName string
	Name     string
	Level    int
	Count    int
}

type PlayerStats struct {
	Player string
	Points int
	Runs   int
	Level  int
}

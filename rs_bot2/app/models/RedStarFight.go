package models

import (
	"encoding/json"
	"fmt"
)

type RedStarFight struct {
	Id                int64
	GameMId           int64
	SolarId           int64
	SendId            string
	Author            string
	Level             string
	Count             int
	Participants      string
	Points            int
	EventId           int
	StartTime         string
	ClientId          int64
	ParticipantsSlice []Participant
}
type Participant struct {
	PlayerId   int64
	PlayerName string
}

func (r *RedStarFight) SaveLevel(DarkOrRed bool, level int) {
	if DarkOrRed {
		r.Level = fmt.Sprintf("drs%d", level)
	} else {
		r.Level = fmt.Sprintf("rs%d", level)
	}
}
func (r *RedStarFight) CountParticipants() int {
	if r.Participants != "" {
		err := json.Unmarshal([]byte(r.Participants), &r.ParticipantsSlice)
		if err != nil {
			fmt.Println(err)
			return 0
		}
		return len(r.ParticipantsSlice)
	}
	return 0
}

package models

import (
	"encoding/json"
	"strings"
)

type QueueActive struct {
	ID            int64                    `json:"id"`
	Messages      map[string]QueueMessages `json:"messages"`
	Data          QueueData                `json:"data"`
	RemainingTime int64                    `json:"remaining_time"`
}

func (q *QueueActive) JsonMarshalMessages() []byte {
	messagesJSON, _ := json.Marshal(q.Messages)
	return messagesJSON
}
func (q *QueueActive) JsonMarshalData() []byte {
	messagesJSON, _ := json.Marshal(q.Data)
	return messagesJSON
}
func (q *QueueActive) JsonUnmarshalMessages(messagesJSON []byte) {
	_ = json.Unmarshal(messagesJSON, &q.Messages)
}
func (q *QueueActive) JsonUnmarshalData(messagesJSON []byte) {
	_ = json.Unmarshal(messagesJSON, &q.Data)
}
func (q *QueueActive) GetFullMap() []byte {
	full := make(map[string]interface{})
	full["Data"] = q.Data
	full["Messages"] = q.Messages
	full["RemainingTime"] = q.RemainingTime
	messagesJSON, _ := json.Marshal(full)
	return messagesJSON
}

type QueueData struct {
	CorporationName string `json:"corporation_name"`
	Name            string `json:"name"`
	UserID          string `json:"user_id"`
	Mention         string `json:"mention"`
	Alt             string `json:"alt,omitempty"`
	Tip             string `json:"tip"`
	Time            string `json:"time"`
	Date            string `json:"date"`
	LvlRS           string `json:"lvl_rs"`
	NumRSName       int    `json:"num_rs_name"`
	NumRSLevel      int    `json:"num_rs_level"`
}

func (q *QueueData) TypeRedStar() (DarkOrRed bool, level string) {
	after, found := strings.CutPrefix(q.LvlRS, "rs")
	if found {
		return false, after
	}
	after, found = strings.CutPrefix(q.LvlRS, "drs")
	if found {
		return true, after
	}

	return false, q.LvlRS
}

type QueueMessages struct {
	TypeMessenger string `json:"type_messenger"`
	//ChatID        string `json:"chat_id"`
	MessageID string `json:"message_id"`
}

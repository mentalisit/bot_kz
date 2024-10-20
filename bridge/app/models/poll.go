package models

type PollStruct struct {
	Author      string            `json:"author"`
	Question    string            `json:"question"`
	Options     []string          `json:"options"`
	CreateTime  int64             `json:"createTime"`
	UrlPoll     string            `json:"urlPoll"`
	Config      BridgeConfig      `json:"config"`
	Votes       []Votes           `json:"votes"`
	PollMessage map[string]string `json:"pollMessage"`
}
type Votes struct {
	Type     string `json:"type"`
	Channel  string `json:"channel"`
	UserName string `json:"userName"`
	Answer   string `json:"answer"`
}
type Request struct {
	Data    map[string]string `json:"data"`
	Options []string          `json:"options"`
}

package models

// TechLevel represents a tech level data structure
type TechLevel struct {
	Level int   `json:"level"`
	Ts    int64 `json:"ts"`
}

// CorpData represents corporation data structure
type CorpData struct {
	Members    []CorpMember `json:"members"`
	Roles      []CorpRole   `json:"roles"`
	FilterID   string       `json:"filterId"`
	FilterName string       `json:"filterName"`
}

// CorpMember represents a member of a corporation.
type CorpMember struct {
	Name         string        `json:"name"`
	UserID       string        `json:"userId"`
	ClientUserID string        `json:"clientUserId"`
	Avatar       string        `json:"avatar"`
	Tech         map[int][]int `json:"tech"`
	AvatarURL    string        `json:"avatarUrl"`
	TimeZone     string        `json:"timeZone"`
	LocalTime    string        `json:"localTime"`
	ZoneOffset   int           `json:"zoneOffset"`
	AfkFor       string        `json:"afkFor"`
	AfkWhen      int           `json:"afkWhen"`
}

// CorpRole represents a corporation role data structure
type CorpRole struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

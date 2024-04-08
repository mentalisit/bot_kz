package models

type Identity struct {
	User  User   `json:"user"`
	Guild Guild  `json:"guild"`
	Token string `json:"token"`
	Type  string `json:"type"`
}

// структура для хранения в будущем игрока и твина и так же нескольких корпораций
type IdentityGET struct {
	User  User    `json:"user"`
	Guild []Guild `json:"guilds"`
	Token string  `json:"token"`
}

type User struct {
	ID            string `json:"id"`
	Username      string `json:"username"`
	Discriminator string `json:"discriminator"`
	Avatar        string `json:"avatar"`
	AvatarURL     string `json:"avatarUrl"`
}

type Guild struct {
	URL  string `json:"url"`
	ID   string `json:"id"`
	Name string `json:"name"`
	Icon string `json:"icon"`
}

type TechLevel struct {
	Level int   `json:"level"`
	Ts    int64 `json:"ts"`
}
type TechLevels map[int]TechLevel

type SyncData struct {
	Ver        int        `json:"ver"`        // Версия данных
	InSync     int        `json:"inSync"`     // Флаг синхронизации
	TechLevels TechLevels `json:"techLevels"` // Коллекция уровней технологии
}
type CorpData struct {
	Members    []CorpMemberint `json:"members"`
	Roles      []CorpRole      `json:"roles"`
	FilterId   string          `json:"filterId"`   // Current filter roleId
	FilterName string          `json:"filterName"` // Name of current filter roleId
}
type CorpMemberint struct {
	Name         string         `json:"name"`
	UserId       string         `json:"userId"`
	ClientUserId string         `json:"clientUserId"`
	Avatar       string         `json:"avatar"`
	Tech         map[int][2]int `json:"tech"`
	AvatarUrl    string         `json:"avatarUrl"`
	TimeZone     string         `json:"timeZone"`
	LocalTime    string         `json:"localTime"`
	ZoneOffset   int            `json:"zoneOffset"` // TZ offset in minutes
	AfkFor       string         `json:"afkFor"`     // readable afk duration
	AfkWhen      int            `json:"afkWhen"`    // Unix Epoch when user returns
}

type CorpMember struct {
	Name         string     `json:"name"`
	UserId       string     `json:"userId"`
	ClientUserId string     `json:"clientUserId"`
	Avatar       string     `json:"avatar"`
	Tech         TechLevels `json:"tech"`
	AvatarUrl    string     `json:"avatarUrl"`
	TimeZone     string     `json:"timeZone"`
	LocalTime    string     `json:"localTime"`
	ZoneOffset   int        `json:"zoneOffset"` // TZ offset in minutes
	AfkFor       string     `json:"afkFor"`     // readable afk duration
	AfkWhen      int        `json:"afkWhen"`    // Unix Epoch when user returns
}

type CorpRole struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

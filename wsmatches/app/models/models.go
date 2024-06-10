package models

import (
	"encoding/xml"
	"time"
)

// Структура для хранения данных о содержимом корзины Amazon S3
type ListBucketResult struct {
	XMLName     xml.Name  `xml:"ListBucketResult"`
	Contents    []Content `xml:"Contents"`
	NextMarker  string    `xml:"NextMarker"`
	IsTruncated bool      `xml:"IsTruncated"`
}

// Структура для хранения данных о файле
type Content struct {
	Key            string `xml:"Key"`
	Generation     int64  `xml:"Generation"`
	MetaGeneration int64  `xml:"MetaGeneration"`
	LastModified   string `xml:"LastModified"`
	ETag           string `xml:"ETag"`
	Size           int64  `xml:"Size"`
}

// Структура для хранения данных о корпорациях
type CorporationsData struct {
	Corporation1Name  string    `json:"Corporation1Name"`
	Corporation1Id    string    `json:"Corporation1Id"`
	Corporation2Name  string    `json:"Corporation2Name"`
	Corporation2Id    string    `json:"Corporation2Id"`
	Corporation1Score int       `json:"Corporation1Score"`
	Corporation2Score int       `json:"Corporation2Score"`
	DateEnded         time.Time `json:"DateEnded"`
}
type Match struct {
	Corporation1Name  string    `json:"Corporation1Name"`
	Corporation1Id    string    `json:"Corporation1Id"`
	Corporation2Name  string    `json:"Corporation2Name"`
	Corporation2Id    string    `json:"Corporation2Id"`
	Corporation1Score int       `json:"Corporation1Score"`
	Corporation2Score int       `json:"Corporation2Score"`
	Corporation1Elo   int       `json:"Corporation1Elo"`
	Corporation2Elo   int       `json:"Corporation2Elo"`
	DateEnded         time.Time `json:"DateEnded"`
	MatchId           string    `json:"MatchId"`
}

func (data *CorporationsData) SortWin() *CorporationsData {
	if data.Corporation2Score > data.Corporation1Score {
		corp := CorporationsData{
			Corporation1Name:  data.Corporation2Name,
			Corporation1Id:    data.Corporation2Id,
			Corporation2Name:  data.Corporation1Name,
			Corporation2Id:    data.Corporation1Id,
			Corporation1Score: data.Corporation2Score,
			Corporation2Score: data.Corporation1Score,
			DateEnded:         data.DateEnded,
		}
		return &corp
	}
	return data
}

//type LevelCorp struct {
//	CorpName string
//	Level    int
//	EndDate  time.Time
//	HCorp    string
//	Percent  int
//}

type LevelCorps struct {
	CorpName   string
	Level      int
	EndDate    time.Time
	HCorp      string
	Percent    int
	LastUpdate time.Time
	Relic      int
}

type FilterCorps struct {
	Corp []Corporation `json:"Corp"`
}
type Corp struct {
	MaxPage int           `json:"MaxPage"`
	Matches []Corporation `json:"matches"`
}
type Ws struct {
	MaxPage int     `json:"MaxPage"`
	Matches []Match `json:"matches"`
}
type Corporation struct {
	Name string `json:"Name"`
	Id   string `json:"Id"`
	Win  int    `json:"Win"`
	Loss int    `json:"Loss"`
	Draw int    `json:"Draw"`
	Elo  int    `json:"Elo"`
}

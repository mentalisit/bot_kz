package corpPercent

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

type CorpsData struct {
	Corp1Name  string    `json:"Corp1Name"`
	Corp2Name  string    `json:"Corp2Name"`
	Corp1Score int       `json:"Corp1Score"`
	Corp2Score int       `json:"Corp2Score"`
	DateEnded  time.Time `json:"DateEnded"`
}

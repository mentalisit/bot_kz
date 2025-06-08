package rsq

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"queue/models"
	"sort"
	"strconv"
	"strings"
)

func GetDataCaprican() []models.QueueStruct {
	sborkz, err := fetchDataSborkz()
	if err != nil || len(sborkz) == 0 {
		return []models.QueueStruct{}
	}
	return sortData(sborkz)
}

type DataCaprican struct {
	Category    string  `json:"category"`
	ChannelID   string  `json:"channelId"`
	DisplayName string  `json:"displayName"`
	GuildID     string  `json:"guildId"`
	GuildInvite *string `json:"guildInvite"`
	GuildName   string  `json:"guildName"`
	Invite      *string `json:"invite"`
	JoinTime    float64 `json:"joinTime"`
	Locale      *string `json:"locale"`
	MessageSent int     `json:"messageSent"`
	RSN         float64 `json:"rsn"`
	Share       int     `json:"share"`
	UserID      string  `json:"userId"`
}

func sortData(data []DataCaprican) []models.QueueStruct {
	var queueSlice []models.QueueStruct
	for _, datum := range data {
		if datum.Share != 1 {
			continue
		}
		var st models.QueueStruct
		st.CorpName = datum.GuildName
		st.Level = fmt.Sprintf("%.1f", datum.RSN)
		st.Count = 1
		queueSlice = append(queueSlice, st)
	}

	// Сортировка по полям CorpName и Level
	sort.Slice(queueSlice, func(i, j int) bool {
		if queueSlice[i].CorpName == queueSlice[j].CorpName {
			return queueSlice[i].Level < queueSlice[j].Level
		}
		return queueSlice[i].CorpName < queueSlice[j].CorpName
	})

	// Объединение структур с одинаковыми CorpName и Level
	merged := merging(queueSlice)

	for i := range merged {
		// Проверяем, содержит ли Level точку
		if strings.Contains(merged[i].Level, ".") {
			parts := strings.Split(merged[i].Level, ".")
			if len(parts) == 2 {
				// Преобразуем часть после точки в число
				afterDot, err := strconv.Atoi(parts[1])
				if err == nil {
					if afterDot == 5 {
						// Заменяем Level на символ до точки + "d"
						merged[i].Level = "drs" + parts[0]
					} else if afterDot == 6 {
						// Заменяем Level на символ до точки + "i"
						merged[i].Level = "iDRS" + parts[0]
					} else if afterDot == 0 {
						merged[i].Level = "rs" + parts[0]
					}
				}
			}
		}
	}

	return merged
}

func fetchDataSborkz() ([]DataCaprican, error) {
	url := "https://api.tsl.rocks/rsq"
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("ошибка запроса: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения тела ответа: %w", err)
	}

	var sborkz []DataCaprican

	//fmt.Println(string(body))

	if err := json.Unmarshal(body, &sborkz); err != nil {
		return nil, fmt.Errorf("ошибка парсинга JSON: %w", err)
	}

	return sborkz, nil
}
func merging(queueSlice []models.QueueStruct) []models.QueueStruct {
	sort.Slice(queueSlice, func(i, j int) bool {
		if queueSlice[i].CorpName == queueSlice[j].CorpName {
			return queueSlice[i].Level < queueSlice[j].Level
		}
		return queueSlice[i].CorpName < queueSlice[j].CorpName
	})
	var merged []models.QueueStruct
	for i := 0; i < len(queueSlice); i++ {
		if i == 0 || queueSlice[i].CorpName != queueSlice[i-1].CorpName || queueSlice[i].Level != queueSlice[i-1].Level {
			// Если это новая комбинация CorpName и Level, добавляем новую структуру
			merged = append(merged, queueSlice[i])
		} else {
			// Если комбинация совпадает с предыдущей, суммируем Count
			merged[len(merged)-1].Count += queueSlice[i].Count
		}
	}
	return merged
}

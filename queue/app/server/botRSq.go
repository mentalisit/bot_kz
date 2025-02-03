package server

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

var merged1 []QueueStruct

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

func sortData(data []DataCaprican) []QueueStruct {
	var queueSlice []QueueStruct
	for _, datum := range data {
		if datum.Share != 1 {
			continue
		}
		var st QueueStruct
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
	merged := Merging(queueSlice)

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

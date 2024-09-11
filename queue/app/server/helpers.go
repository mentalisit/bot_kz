package server

import (
	"queue/models"
	"sort"
	"strconv"
	"strings"
)

func deleteChannelName(s models.Sborkz) string {
	split := strings.Split(s.Corpname, ".")
	if len(split) == 2 {
		return split[0]
	}
	return s.Corpname
}

func Merging(queueSlice []QueueStruct) []QueueStruct {
	var merged []QueueStruct
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
func ConvertingTumchaToQueueStruct(rs map[string][]models.Tumcha) (q []QueueStruct) {
	lvl := func(t models.Tumcha) (s string) {
		if t.Vid == "black" {
			s = "d"
		}
		s += strconv.Itoa(t.Level)
		return s
	}
	for cName, tumchas := range rs {
		r := QueueStruct{
			CorpName: cName,
		}
		for i, tumcha := range tumchas {
			r.Level = lvl(tumcha)
			r.Count = i + 1

			q = append(q, r)
		}
	}
	return Merging(q)
}

func ConvertingSborkzToQueueStruct(m []models.Sborkz) (q []QueueStruct) {
	for _, sborkz := range m {
		corp := deleteChannelName(sborkz)
		q = append(q, QueueStruct{
			CorpName: corp,
			Level:    sborkz.Lvlkz,
			Count:    1,
		})
	}
	sort.Slice(q, func(i, j int) bool {
		if q[i].CorpName == q[j].CorpName {
			return q[i].Level < q[j].Level
		}
		return q[i].CorpName < q[j].CorpName
	})
	return Merging(q)
}

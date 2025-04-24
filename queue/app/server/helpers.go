package server

import (
	"queue/models"
	"sort"
	"strings"
)

func deleteChannelName(s models.Sborkz) string {
	split := strings.Split(s.Corpname, ".")
	if len(split) == 2 {
		return split[0]
	} else {
		i := strings.Split(s.Corpname, "/")
		if len(i) == 2 {
			return i[0]
		}
	}

	return s.Corpname
}

func Merging(queueSlice []models.QueueStruct) []models.QueueStruct {
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
func ConvertingTumchaToQueueStruct(rs map[string][]models.QueueStruct) (q []models.QueueStruct) {
	for _, structs := range rs {
		q = append(q, structs...)
	}
	return Merging(q)
}

//func ConvertingTumchaToQueueStruct(rs map[string][]models.Tumcha) (q []models.QueueStruct) {
//	lvl := func(t models.Tumcha) (s string) {
//		if t.Vid == "black" {
//			s = "drs"
//		} else {
//			s = "rs"
//		}
//		s += strconv.Itoa(t.Level)
//		return s
//	}
//	for cName, tumchas := range rs {
//		r := models.QueueStruct{
//			CorpName: cName,
//		}
//		for _, tumcha := range tumchas {
//			r.Level = lvl(tumcha)
//			r.Count = 1
//
//			q = append(q, r)
//		}
//	}
//	return Merging(q)
//}

func ConvertingSborkzToQueueStruct(m []models.Sborkz) (q []models.QueueStruct) {
	for _, sborkz := range m {
		corp := deleteChannelName(sborkz)
		q = append(q, models.QueueStruct{
			CorpName: corp,
			Level:    sborkz.Lvlkz,
			Count:    1,
		})
	}
	return Merging(q)
}

package hspublic

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"time"
	"ws/models"
)

func (h *HS) DownloadFile(fileName string, newContent []models.Content) {
	var count int
	var corpsdata []models.Match

	for _, cont := range newContent {
		count++
		corpData := h.r.ReadCorpData(cont.Key)
		if corpData == nil {
			corpData = h.GetCorporationsData(cont.Key)
			h.r.SaveCorpDate(cont.Key, *corpData)
		}
		corpData = corpData.SortWin()
		mid := models.Match{
			Corporation1Name:  corpData.Corporation1Name,
			Corporation1Id:    corpData.Corporation1Id,
			Corporation2Name:  corpData.Corporation2Name,
			Corporation2Id:    corpData.Corporation2Id,
			Corporation1Score: corpData.Corporation1Score,
			Corporation2Score: corpData.Corporation2Score,
			DateEnded:         corpData.DateEnded,
			MatchId:           cont.Key,
		}

		corpsdata = append(corpsdata, mid)
	}
	fmt.Println("count", count)

	h.SaveCorpList(corpsdata)

	var match []models.Match
	for _, corpsdatum := range corpsdata {
		var current models.Match
		current = corpsdatum
		current.Corporation1Elo = elo[current.Corporation1Id]
		current.Corporation2Elo = elo[current.Corporation2Id]

		match = append(match, current)
	}
	sort.Slice(match, func(i, j int) bool {
		return match[i].DateEnded.After(match[j].DateEnded)
	})

	h.SaveFile(fileName, match)
}

func (h *HS) SaveFile(fileName string, corpsdata any) {
	marshal, err := json.Marshal(corpsdata)
	if err != nil {
		h.log.ErrorErr(err)
		return
	}
	path := fmt.Sprintf("docker/ws/%s.json", fileName)
	err = os.WriteFile(path, marshal, 0644)
	if err != nil {
		h.log.ErrorErr(err)
		return
	}
	fmt.Printf("%s файл сохранен %s\n", time.Now().Format(time.DateTime), path)
}
func (h *HS) SaveCorpList(corps []models.Match) {
	corpMap := make(map[string]models.Corporation)
	for _, corp := range corps {
		c := models.Corporation{
			Name: corp.Corporation1Name,
			Id:   corp.Corporation1Id,
		}
		if corpMap[corp.Corporation1Id] != c {
			corpMap[corp.Corporation1Id] = c
		}

		c = models.Corporation{
			Name: corp.Corporation2Name,
			Id:   corp.Corporation2Id,
		}
		if corpMap[corp.Corporation2Id] != c {
			corpMap[corp.Corporation2Id] = c
		}
	}
	h.SaveCorpListCount(corpMap, corps)
}
func (h *HS) SaveCorpListCount(m map[string]models.Corporation, corps []models.Match) {
	var cc []models.Corporation
	for _, corporation := range m {
		var c models.Corporation
		c.Name = corporation.Name
		c.Id = corporation.Id
		for _, match := range corps {
			if corporation.Id == match.Corporation1Id {
				if match.Corporation1Score > match.Corporation2Score {
					c.Win += 1
				} else if match.Corporation1Score < match.Corporation2Score {
					c.Loss += 1
				} else if match.Corporation1Score == match.Corporation2Score {
					c.Draw += 1
				}
			} else if corporation.Id == match.Corporation2Id {
				if match.Corporation2Score > match.Corporation1Score {
					c.Win += 1
				} else if match.Corporation2Score < match.Corporation1Score {
					c.Loss += 1
				} else if match.Corporation1Score == match.Corporation2Score {
					c.Draw += 1
				}
			}
		}
		cc = append(cc, c)
	}
	sort.Slice(corps, func(i, j int) bool {
		return corps[i].DateEnded.Before(corps[j].DateEnded)
	})
	EloLogic(corps, cc)
}

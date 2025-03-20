package hspublic

import (
	"sort"
	"ws/models"
)

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

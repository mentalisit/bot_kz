package hspublic

import (
	"ws/models"
)

func (h *HS) SavePercent(newContent []models.Content) {

	var count int

	listHCorp := make(map[string]models.LevelCorp)

	all, err := h.p.ReadCorpLevelAll()
	if err != nil {
		h.log.ErrorErr(err)
		return
	}

	for _, corp := range all {
		listHCorp[corp.HCorp] = corp
	}

	for _, cont := range newContent {
		count++
		corpData := h.r.ReadCorpData(cont.Key)
		if corpData == nil {
			corpData = h.GetCorporationsData(cont.Key)
			h.r.SaveCorpDate(cont.Key, *corpData)
		}

		if listHCorp[corpData.Corporation1Name].HCorp != "" {
			c := listHCorp[corpData.Corporation1Name]
			if c.EndDate.Before(corpData.DateEnded) {
				c.EndDate = corpData.DateEnded
				c.Percent = c.Level - 1

				h.p.InsertUpdateCorpLevel(c)
			}
		}
	}
}

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
		f := corpData.SortWin()

		if listHCorp[f.Corporation1Name].HCorp != "" {
			c := listHCorp[f.Corporation1Name]
			if c.EndDate.Before(f.DateEnded) {
				c.EndDate = f.DateEnded
				c.Percent = c.Level - 1

				h.p.InsertUpdateCorpLevel(c)
			}
		}
	}
}

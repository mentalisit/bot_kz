package hspublic

import (
	"fmt"
	"time"
	"ws/models"
)

func (h *HS) SavePercent(newContent []models.Content) {

	listHCorps := h.getMapLevelCorps()

	for _, cont := range newContent {
		corpData := h.r.ReadCorpData(cont.Key)
		if corpData == nil {
			corpData = h.GetCorporationsData(cont.Key)
			h.r.SaveCorpDate(cont.Key, *corpData)
		}

		f := corpData.SortWin()

		if listHCorps[f.Corporation1Name].HCorp != "" {
			c := listHCorps[f.Corporation1Name]
			if c.EndDate.Before(f.DateEnded) {
				c.EndDate = f.DateEnded
				c.Percent = c.Level - 1

				h.p.InsertUpdateCorpsLevel(c)
			}
		}
		h.Relic(listHCorps, corpData)
	}

	h.recalculateCorpLevel()
}
func (h *HS) Relic(list map[string]models.LevelCorps, corpData *models.CorporationsData) {
	if list[corpData.Corporation1Name].HCorp != "" {
		c, _ := h.p.ReadCorpsLevel(corpData.Corporation1Name)
		if c.LastUpdate.Before(corpData.DateEnded) {
			c.LastUpdate = corpData.DateEnded
			c.Relic = c.Relic + corpData.Corporation1Score

			h.p.InsertUpdateCorpsLevel(c)
			fmt.Printf("InsertUpdateCorpsLevel %+v\n", c)
			time.Sleep(1 * time.Second)
		}
	}
	if list[corpData.Corporation2Name].HCorp != "" {
		c, _ := h.p.ReadCorpsLevel(corpData.Corporation2Name)
		if c.LastUpdate.Before(corpData.DateEnded) {
			c.LastUpdate = corpData.DateEnded
			c.Relic = c.Relic + corpData.Corporation2Score

			h.p.InsertUpdateCorpsLevel(c)
			fmt.Printf("InsertUpdateCorpsLevel %+v\n", c)
			time.Sleep(1 * time.Second)
		}
	}
}
func (h *HS) getMapLevelCorps() map[string]models.LevelCorps {
	listHCorp := make(map[string]models.LevelCorps)

	all, err := h.p.ReadCorpsLevelAll()
	if err != nil {
		h.log.ErrorErr(err)
		return listHCorp
	}

	for _, corp := range all {
		listHCorp[corp.HCorp] = corp
	}
	return listHCorp
}
func (h *HS) recalculateCorpLevel() {
	level := func(r int) int {
		levels := []struct {
			relic int
			level int
		}{
			{60000, 21},
			{50000, 20},
			{40000, 19},
			{32000, 18},
			{25000, 17},
			{20000, 16},
			{16000, 15},
			{13000, 14},
			{11000, 13},
			{9000, 12},
			{7000, 11},
			{5000, 10},
			{3000, 9},
			{2000, 8},
			{1000, 7},
			{500, 6},
			{250, 5},
			{100, 4},
			{30, 3},
			{1, 2},
		}

		for _, l := range levels {
			if r > l.relic {
				return l.level
			}
		}
		return 1
	}

	list := h.getMapLevelCorps()

	for _, corps := range list {
		if level(corps.Relic) != corps.Level {
			corps.Level = level(corps.Relic)
			corps.Percent = corps.Level - 1
			h.p.InsertUpdateCorpsLevel(corps)
			h.log.InfoStruct("recalculateCorpLevel", corps)
		}
	}
}

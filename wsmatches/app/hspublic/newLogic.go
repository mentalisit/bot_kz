package hspublic

import (
	"fmt"
	"time"
	"ws/models"
)

func (h *HS) getMapLevelCorpsV2() map[string]models.CorpInfo {
	listCorp := make(map[string]models.CorpInfo)

	all, err := h.p.GetAllCorpInfo()
	if err != nil {
		h.log.ErrorErr(err)
		return listCorp
	}

	for _, corp := range all {
		if corp.CorpID != "" && !corp.Webhook {
			listCorp[corp.CorpID] = corp
		}
	}
	return listCorp
}

func (h *HS) RelicV2(list map[string]models.CorpInfo, corpData *models.CorporationsData) {
	if list[corpData.Corporation1Id].CorpName != "" {
		c, _ := h.p.ReadCorpInfoByCorpID(corpData.Corporation1Id)
		if c.LastUpdate.Before(corpData.DateEnded) {
			c.LastUpdate = corpData.DateEnded
			c.XP = c.XP + 100

			_ = h.p.UpdateCorpInfo(*c)
			fmt.Printf("UpdateCorpInfo %+v\n", c)
			time.Sleep(1 * time.Second)
		}
	}
	if list[corpData.Corporation2Id].CorpName != "" {
		c, _ := h.p.ReadCorpInfoByCorpID(corpData.Corporation2Id)
		if c.LastUpdate.Before(corpData.DateEnded) {
			c.LastUpdate = corpData.DateEnded
			c.XP = c.XP + 40

			_ = h.p.UpdateCorpInfo(*c)
			fmt.Printf("UpdateCorpInfo %+v\n", c)
			time.Sleep(1 * time.Second)
		}
	}
}

func (h *HS) recalculateCorpLevelV2() {
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

	listCorpId := h.getMapLevelCorpsV2()

	for _, info := range listCorpId {
		if level(info.XP) != info.Level {
			info.Level = level(info.XP)
			_ = h.p.UpdateCorpInfo(info)
			h.log.InfoStruct("recalculateCorpInfo", info)
		}
	}

}

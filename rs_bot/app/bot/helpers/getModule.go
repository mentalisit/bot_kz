package helpers

import (
	"rs/models"
	"strings"
)

func (h *Helpers) Get2TechDataUserId(name, userID string, mAcc *models.MultiAccount) (genesis, enrich, rsextender int) {
	technicalData, err := h.storage.Postgres.TechnologiesGetAll(mAcc.UUID)
	if err != nil {
		h.log.ErrorErr(err)
	} else {
		if len(technicalData) < 1 {
			return
		}

		if (mAcc.DiscordID == userID && mAcc.DiscordUsername == name) ||
			(mAcc.TelegramID == userID && mAcc.TelegramUsername == name) ||
			(mAcc.WhatsappID == userID && mAcc.WhatsappUsername == name) {
			name = mAcc.Nickname
		}
	}

	for _, datum := range technicalData {
		if strings.ToLower(datum.Name) == strings.ToLower(name) {
			rsextender = datum.Tech[603].Level
			enrich = datum.Tech[503].Level
			genesis = datum.Tech[508].Level
		}
	}
	if rsextender == 0 && enrich == 0 && genesis == 0 && len(technicalData) != 0 {
		rsextender = technicalData[0].Tech[603].Level
		enrich = technicalData[0].Tech[503].Level
		genesis = technicalData[0].Tech[508].Level
	}

	if rsextender != 0 && enrich != 0 && genesis != 0 {
		module := h.storage.Postgres.ModuleReadUUID(mAcc.UUID, name)
		newModule := models.Module{
			Uid:  mAcc.UUID,
			Name: name,
			Gen:  genesis,
			Enr:  enrich,
			Rse:  rsextender,
		}
		if module == nil {
			h.storage.Postgres.ModuleInsertUUID(newModule)
		} else {
			h.storage.Postgres.ModuleUpdateUUID(newModule)
		}
	}

	return genesis, enrich, rsextender
}

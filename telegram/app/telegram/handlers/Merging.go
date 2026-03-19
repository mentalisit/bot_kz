package handlers

import (
	"github.com/mentalisit/restapi/models"
)

func (h *WebAppHandler) mergeAccounts(existingAccount, discordAcc *models.MultiAccount) (*models.MultiAccount, error) {
	h.log.InfoStruct("Merging accounts", map[string]interface{}{
		"existing_account": existingAccount,
		"discord_account":  discordAcc,
	})

	// Объединяем alts из обоих аккаунтов
	if existingAccount.Alts == nil {
		existingAccount.Alts = []string{}
	}
	if discordAcc.Alts == nil {
		discordAcc.Alts = []string{}
	}

	// Добавляем все alts из Discord аккаунта в основной
	allAlts := append(existingAccount.Alts, discordAcc.Alts...)

	// Убираем дубликаты
	uniqueAlts := make([]string, 0)
	seen := make(map[string]bool)
	for _, alt := range allAlts {
		if !seen[alt] {
			uniqueAlts = append(uniqueAlts, alt)
			seen[alt] = true
		}
	}
	existingAccount.Alts = uniqueAlts

	// Обновляем существующий аккаунт с объединенными данными
	updatedAccount, err := h.storage.Db.UpdateMultiAccount(*existingAccount)
	if err != nil {
		h.log.Error("Error updating account with alts: " + err.Error())
		return nil, err
	}

	corpMemberDS, _ := h.storage.Db.CorpMemberByUId(discordAcc.UUID)
	h.log.InfoStruct("corpMemberDS need delete ", corpMemberDS)

	tech, _ := h.storage.Db.TechnologiesGetAll(discordAcc.UUID)
	h.log.InfoStruct("tech need delete ", tech)

	// Удаляем старую Discord запись после успешного объединения
	//_ = h.storage.Db.DeleteMultiAccount(discordAcc.UUID)
	//
	//h.log.Info(fmt.Sprintf("Accounts merged successfully - main_uuid: %s, merged_discord_uuid: %s",
	//	existingAccount.UUID.String(), discordAcc.UUID.String()))

	return updatedAccount, nil
}

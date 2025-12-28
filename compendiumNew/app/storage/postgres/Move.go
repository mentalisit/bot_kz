package postgres

import (
	"compendium/models"
)

func (d *Db) ReadAndMoveToNewLogic(m models.IncomingMessage) {
	multiAccount, _ := d.Multi.FindMultiAccountByUserId(m.NameId)
	if multiAccount != nil && multiAccount.Nickname != "" {
		corpMember, _ := d.Multi.CorpMemberByUId(multiAccount.UUID)
		if corpMember != nil {
			//d.Multi.TechnologiesGetAllCorpMember()
		}
	}
}

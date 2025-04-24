package helpers

import (
	"fmt"
	"rs/models"
)

func moduleToEmodji(genesis, enrich, rsextender int) (gen, enr, rse string) {
	gen = fmt.Sprintf("<:genesis:1199068748280242237> %d ", genesis)
	enr = fmt.Sprintf("<:enrich:1199068793633251338> %d ", enrich)
	rse = fmt.Sprintf("<:rse:1199068829511335946> %d ", rsextender)
	return gen, enr, rse
}
func (h *Helpers) ReadNameModulesUUID(in models.InMessage, name string) {
	var multiAccount *models.MultiAccount
	if name == "" {
		name = in.Username
		multiAccount, _ = h.storage.Postgres.FindMultiAccountByUserId(in.UserId)
		if multiAccount != nil {
			name = multiAccount.Nickname
		}
	}
	if multiAccount == nil {
		h.log.Error("multiAccount==nil")
		return
	}
	gen, enr, rse := Get2TechDataUserId(name, in.UserId, in.Ds.Guildid)

	m := models.Module{
		Uid:  multiAccount.UUID,
		Name: name,
		Gen:  gen,
		Enr:  enr,
		Rse:  rse,
	}
	moduleReadUUID := h.storage.Postgres.ModuleReadUUID(multiAccount.UUID, name)
	if moduleReadUUID == nil {
		h.storage.Postgres.ModuleInsertUUID(m)
	} else {
		if moduleReadUUID.Gen != gen || moduleReadUUID.Enr != enr || moduleReadUUID.Rse != rse {
			h.storage.Postgres.ModuleUpdateUUID(m)
		}
	}
}

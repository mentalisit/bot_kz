package logic

import (
	"compendium/models"

	"github.com/google/uuid"
)

func (c *Hs) removeMember(m models.IncomingMessage) bool {
	if m.Text == "%GuildMemberRemove" && m.MentionName == "" {
		err := c.corpMember.CorpMemberDelete(m.MGuild.GuildId(), m.NameId)
		if err != nil {
			c.log.ErrorErr(err)
		}
		corpMember, err := c.db.Multi.CorpMemberByUId(m.MultiAccount.UUID)
		if err == nil && corpMember != nil {
			if corpMember.Exist(m.MGuild.GId) {
				corpMember.GuildIds = removeGuild(corpMember.GuildIds, m.MGuild.GId)
				err = c.db.Multi.CorpMemberUpdateGuildIds(*corpMember)
				if err != nil {
					c.log.ErrorErr(err)
				}
			}
		}
		memberByUId, err := c.db.V2.CorpMemberByUId(m.MAcc.UUID)
		if err == nil && memberByUId != nil {
			if memberByUId.Exist(m.MGuild.GId) {
				memberByUId.GuildIds = removeGuild(memberByUId.GuildIds, m.MGuild.GId)
				err = c.db.V2.CorpMemberUpdate(*memberByUId)
				if err != nil {
					c.log.ErrorErr(err)
				}
			}
		}
		return true
	}
	return false
}
func removeGuild(g []uuid.UUID, gud uuid.UUID) []uuid.UUID {
	var gNew []uuid.UUID
	for _, u := range g {
		if u != gud {
			gNew = append(gNew, u)
		}
	}
	return gNew
}

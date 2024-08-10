package logic

import "compendium/models"

func (c *Hs) removeMember(m models.IncomingMessage) bool {
	if m.Text == "%GuildMemberRemove" && m.MentionName == "" {
		err := c.corpMember.CorpMemberDelete(m.GuildId, m.NameId)
		if err != nil {
			c.log.ErrorErr(err)
		}
		return true
	}
	return false
}

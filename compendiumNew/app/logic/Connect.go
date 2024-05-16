package logic

import (
	"compendium/logic/generate"
	"compendium/models"
	"fmt"
)

func (c *Hs) connect(m models.IncomingMessage) {
	err := c.sendDM(m, fmt.Sprintf(c.getText(m, "CODE_FOR_CONNECT"), m.GuildName))
	if err != nil && err.Error() == "forbidden" {
		c.sendChat(m, fmt.Sprintf(c.getText(m, "ERROR_SEND"), m.MentionName))
		return
	} else if err != nil {
		c.log.ErrorErr(err)
	}
	c.sendChat(m, fmt.Sprintf(c.getText(m, "INSTRUCTIONS_SEND"), m.MentionName))
	newIdentify, cm := c.generate(m)
	code := generate.GenerateFormattedString(newIdentify)
	err = c.guilds.GuildInsert(newIdentify.Guild)
	if err != nil {
		c.log.ErrorErr(err)
		return
	}
	err = c.users.UsersInsert(newIdentify.User)
	if err != nil {
		c.log.ErrorErr(err)
		return
	}
	err = c.corpMember.CorpMemberInsert(cm)
	if err != nil {
		c.log.ErrorErr(err)
		return
	}
	err = c.listUser.ListUserInsert(newIdentify.Token, newIdentify.User.ID, newIdentify.Guild.ID)
	if err != nil {
		c.log.ErrorErr(err)
		return
	}
	err = c.sendDM(m, code)
	if err != nil {
		c.log.ErrorErr(err)
		return
	}
	urlLink := "https://mentalisit.github.io/HadesSpace/"
	urlLinkAdd := "compendiumTech?c2=" + code + "&lang=" + m.Language
	err = c.sendDM(m, fmt.Sprintf(c.getText(m, "PLEASE_PASTE_CODE"), urlLink, urlLink+urlLinkAdd))
	if err != nil {
		c.log.ErrorErr(err)
		return
	}
}

func (c *Hs) generate(m models.IncomingMessage) (models.Identity, models.CorpMember) {
	//проверить если есть NameId то предложить соединить для двух корпораций
	identity := models.Identity{
		User: models.User{
			ID:       m.NameId,
			Username: m.Name,
			//Discriminator: "",
			Avatar:    m.AvatarF,
			AvatarURL: m.Avatar,
			Alts:      []string{},
		},
		Guild: models.Guild{
			URL:  m.GuildAvatar,
			ID:   m.GuildId,
			Name: m.GuildName,
			Icon: m.GuildAvatarF,
		},
		Token: generate.GenerateToken(),
		//Type:  c.in.Type,
	}
	cm := models.CorpMember{
		Name:      m.Name,
		UserId:    m.NameId,
		GuildId:   m.GuildId,
		Avatar:    m.AvatarF,
		Tech:      map[int][2]int{},
		AvatarUrl: m.Avatar,
	}
	return identity, cm
}

package logic

import (
	"compendium/logic/generate"
	"compendium/models"
	"fmt"
)

func (c *Hs) connect() {
	err := c.sendDM(fmt.Sprintf("Код для подключения приложения к серверу %s.", c.in.GuildName))
	if err != nil && err.Error() == "forbidden" {
		c.sendChat(c.in.MentionName +
			" пожалуйста отправьте мне команду старт в личных сообщениях, " +
			"я как бот не могу первый отправить вам личное сообщение. " +
			"И после повторите команду  ")
		return
	} else if err != nil {
		c.log.ErrorErr(err)
	}
	c.sendChat(c.in.MentionName + ", Инструкцию отправили вам в Директ.")
	newIdentify, cm := c.generate()
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
	err = c.sendDM(code)
	if err != nil {
		c.log.ErrorErr(err)
		return
	}
	err = c.sendDM("Пожалуйста, вставьте код в приложение\nhttps://mentalisit.github.io/HadesSpace/ \n" +
		"или просто перейдите по ссылке для автоматической авторизации  \n" +
		"https://mentalisit.github.io/HadesSpace/compendiumTech?c2=" + code)
	if err != nil {
		c.log.ErrorErr(err)
		return
	}
}

func (c *Hs) generate() (models.Identity, models.CorpMember) {
	//проверить если есть NameId то предложить соединить для двух корпораций
	identity := models.Identity{
		User: models.User{
			ID:       c.in.NameId,
			Username: c.in.Name,
			//Discriminator: "",
			Avatar:    c.in.AvatarF,
			AvatarURL: c.in.Avatar,
			Alts:      []string{},
		},
		Guild: models.Guild{
			URL:  c.in.GuildAvatar,
			ID:   c.in.GuildId,
			Name: c.in.GuildName,
			Icon: c.in.GuildAvatarF,
		},
		Token: generate.GenerateToken(),
		//Type:  c.in.Type,
	}
	cm := models.CorpMember{
		Name:      c.in.Name,
		UserId:    c.in.NameId,
		GuildId:   c.in.GuildId,
		Avatar:    c.in.AvatarF,
		Tech:      map[int][2]int{},
		AvatarUrl: c.in.Avatar,
	}
	return identity, cm
}

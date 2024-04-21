package Compendium

import (
	"compendium/Compendium/generate"
	"compendium/models"
	"fmt"
)

func (c *Compendium) connect() {

	err := c.sendDM(fmt.Sprintf("Код для подключения приложения к серверу %s.", c.in.GuildName))
	if err != nil && err.Error() == "forbidden" {
		c.sendChat(c.in.MentionName +
			" пожалуйста отправьте мне команду старт в личных сообщениях, " +
			"я как бот не могу первый отправить вам личное сообщение. " +
			"И после повторите команду  ")
		return
	}
	c.sendChat(c.in.MentionName + ", Инструкцию отправили вам в Директ.")
	code := generate.GenerateFormattedString(c.generate())
	err = c.sendDM(code)
	if err != nil {
		return
	}
	err = c.sendDM("Пожалуйста, вставьте код в приложение\nhttps://mentalisit.github.io/HadesSpace/ \n" +
		"или просто перейдите по ссылке для автоматической авторизации  \n" +
		"https://mentalisit.github.io/HadesSpace/compendiumCorp?c=" + code)
	if err != nil {
		return
	}
}

func (c *Compendium) generate() models.Identity {
	//проверить если есть NameId то предложить соединить для двух корпораций
	identity := models.Identity{
		User: models.User{
			ID:            c.in.NameId,
			Username:      c.in.Name,
			Discriminator: "",
			Avatar:        c.in.AvatarF,
			AvatarURL:     c.in.Avatar,
		},
		Guild: models.Guild{
			URL:  c.in.GuildAvatar,
			ID:   c.in.GuildId,
			Name: c.in.GuildName,
			Icon: c.in.GuildAvatarF,
		},
		Token: generate.GenerateToken(),
		Type:  c.in.Type,
	}
	return identity
}

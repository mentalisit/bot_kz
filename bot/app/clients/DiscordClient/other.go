package DiscordClient

import (
	"fmt"
	gt "github.com/bas24/googletranslatefree"
	"github.com/bwmarrin/discordgo"
	"kz_bot/clients/restapi"
	"kz_bot/models"
	"kz_bot/pkg/utils"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func (d *Discord) CheckAdmin(nameid string, chatid string) bool {
	perms, err := d.S.UserChannelPermissions(nameid, chatid)
	if err != nil {
		d.log.ErrorErr(err)
	}
	if perms&discordgo.PermissionAdministrator != 0 {
		return true
	} else {
		return false
	}
}
func (d *Discord) RoleToIdPing(rolePing, guildid string) (string, error) {
	if guildid == "" {
		d.log.Panic("почему то нет гуилд ид")
	}
	if d.roles[guildid] == nil {
		d.roles[guildid] = make(map[string]string)
	}
	if d.roles[guildid][rolePing] != "" {
		return d.roles[guildid][rolePing], nil
	}

	g, err := d.S.Guild(guildid)
	if err != nil {
		return rolePing, err
	}
	exist, role := d.roleExists(g, rolePing)
	if !exist {
		//создаем роль и возврашаем пинг
		role, err = d.createRole(rolePing, guildid)
		if err != nil {
			d.roles[guildid][rolePing] = rolePing
			return rolePing, err
		}
		d.roles[guildid][rolePing] = role.Mention()
		return role.Mention(), nil
	} else {
		d.roles[guildid][rolePing] = role.Mention()
		return role.Mention(), nil
	}
}

//func (d *Discord) TextToRoleRsPing(rolePing, guildid string) string {
//
//	if guildid == "" {
//		d.log.Panic("почему то нет гуилд ид")
//		panic("почему то нет гуилд ид")
//	}
//	g, err := d.S.Guild(guildid)
//	if err != nil {
//		d.log.ErrorErr(err)
//	}
//	exist, role := d.roleExists(g, rolePing)
//	if !exist {
//		return fmt.Sprintf("`роль %s не найдена в %s`", rolePing, g.Name)
//	} else {
//		return role.Mention()
//	}
//}

//func (d *Discord) DMchannel(AuthorID string) (chatidDM string) {
//	create, err := d.S.UserChannelCreate(AuthorID)
//	if err != nil {
//		return ""
//	}
//	chatidDM = create.ID
//	return chatidDM
//}

func (d *Discord) CleanChat(chatid, mesid, text string) {
	res := strings.HasPrefix(text, ".")
	if !res { //если нет префикса  то удалить через 3 минуты
		go d.DeleteMesageSecond(chatid, mesid, 180)
	}
}

// получаем есть ли роль и саму роль
func (d *Discord) roleExists(g *discordgo.Guild, nameRoles string) (bool, *discordgo.Role) {
	nameRoles = strings.ToLower(nameRoles)

	for _, role := range g.Roles {
		if role.Name == "@everyone" {
			continue
		}
		if strings.ToLower(role.Name) == nameRoles {
			return true, role
		}
	}
	return false, nil
}
func (d *Discord) GuildChatName(chatid, guildid string) string {
	g, err := d.S.Guild(guildid)
	if err != nil {
		d.log.ErrorErr(err)
	}
	chatName := g.Name
	channels, _ := d.S.GuildChannels(guildid)

	for _, r := range channels {
		if r.ID == chatid {
			chatName = chatName + "." + r.Name
		}
	}
	return chatName
}
func (d *Discord) createRole(rolPing, guildid string) (*discordgo.Role, error) {
	t := true
	perm := int64(37080064)
	create, err := d.S.GuildRoleCreate(guildid, &discordgo.RoleParams{
		Name:        rolPing,
		Permissions: &perm,
		Mentionable: &t,
	})
	if err != nil {
		return nil, err
	}
	return create, nil
}

func (d *Discord) getLanguage(lang, key string) string {
	return d.storage.Dictionary.GetText(lang, key)
}

func (d *Discord) CleanOldMessageChannel(chatId, lim string) {
	limit, _ := strconv.Atoi(lim)
	if limit == 0 {
		return
	}
	messages, err := d.S.ChannelMessages(chatId, limit, "", "", "")
	if err != nil {
		d.log.ErrorErr(err)
		return
	}
	for _, message := range messages {
		if message.WebhookID == "" {
			if !message.Author.Bot {
				d.DeleteMessage(chatId, message.ID)
				continue
			}
			if !strings.HasPrefix(message.Content, ".") {
				d.DeleteMessage(chatId, message.ID)
				continue
			}
		}
	}
}

func (d *Discord) avatar(m *discordgo.MessageCreate) bool {
	str, ok := strings.CutPrefix(m.Content, ". ")
	if ok {
		arg := strings.Split(strings.ToLower(str), " ")
		if len(arg) == 2 {
			if arg[0] == "ава" || arg[0] == "ava" {
				mentionIds := userMentionRE.FindAllStringSubmatch(arg[1], -1)
				if len(mentionIds) > 0 {
					members, err := d.S.GuildMembers(m.GuildID, "", 999)
					if err != nil {
						d.log.ErrorErr(err)
					}
					for _, member := range members {
						if member.User.ID == mentionIds[0][1] {
							aname := m.Author.Username
							if m.Member.Nick != "" {
								aname = m.Member.Nick
							}
							name := member.User.Username
							if member.Nick != "" {
								name = member.Nick
							}
							em := &discordgo.MessageEmbed{
								Title: fmt.Sprintf("Аватар %s по запросу %s", name, aname),
								Color: 14232643,
								Image: &discordgo.MessageEmbedImage{
									URL: member.AvatarURL("2048"),
								},
								Author: nil,
							}
							embed, err := d.S.ChannelMessageSendEmbed(m.ChannelID, em)
							if err != nil {
								fmt.Println(err.Error())
								return false
							}
							go d.DeleteMesageSecond(m.ChannelID, embed.ID, 183)
							go d.DeleteMesageSecond(m.ChannelID, m.ID, 30)
							return true
						}
					}
				}
			}
		}
		if arg[0] == "стек" {
			utils.PrintGoroutinesStack()
			return true
		}
	}
	return false
}

func (d *Discord) getAuthorName(m *discordgo.MessageCreate) string {
	username := m.Author.Username
	if m.Member != nil && m.Member.Nick != "" {
		username = m.Member.Nick
	}
	return username
}

func (d *Discord) latinOrNot(m *discordgo.MessageCreate) {
	cyrillicPattern := regexp.MustCompile(`[а-яА-ЯґҐєЄіІїЇ]`)
	if len(m.Content) >= 1 {
		if !cyrillicPattern.MatchString(m.Content) {
			channel, err := d.S.Channel(m.ChannelID)
			if err != nil {
				return
			}
			gostPattern := regexp.MustCompile(`гост`)

			if gostPattern.MatchString(channel.Name) {
				//возможно нужно доп условие
				go func() {
					text2, _ := gt.Translate(m.Content, "auto", "ru")
					mes := d.SendWebhook(text2, m.Author.Username, m.ChannelID, m.Author.AvatarURL("128"))
					d.DeleteMesageSecond(m.ChannelID, mes, 90)
				}()
			}
		}
	}
}
func (d *Discord) transtale(m *discordgo.Message, lang string, r *discordgo.MessageReactionAdd) {
	text2, _ := gt.Translate(m.Content, "auto", lang)
	go func() {
		time.Sleep(30 * time.Second)
		err := d.S.MessageReactionRemove(r.ChannelID, r.MessageID, r.Emoji.Name, r.UserID)
		if err != nil {
			fmt.Println("Ошибка удаления реакции", err)
		}
	}()
	mes := d.SendWebhook(text2, m.Author.Username, m.ChannelID, m.Author.AvatarURL("128"))
	d.DeleteMesageSecond(m.ChannelID, mes, 90)
}
func (d *Discord) dmChannel(AuthorID string) (chatidDM string) {
	create, err := d.S.UserChannelCreate(AuthorID)
	if err != nil {
		return ""
	}
	chatidDM = create.ID
	return chatidDM
}
func (d *Discord) GetRoles(guildId string) []models.CorpRole {
	roles, err := d.S.GuildRoles(guildId)
	if err != nil {
		d.log.ErrorErr(err)
		return nil
	}
	var guildRole []models.CorpRole
	for _, role := range roles {
		r := models.CorpRole{
			Name: role.Name,
			ID:   role.ID,
		}
		if r.Name == "@everyone" {
			r.ID = ""
		}

		guildRole = append(guildRole, r)
	}
	return guildRole
}
func (d *Discord) CheckRole(guildId, memderId, roleid string) bool {
	if roleid == "" {
		return true
	}
	member, err := d.S.GuildMember(guildId, memderId)
	if err != nil {
		d.log.ErrorErr(err)
		d.log.Info(fmt.Sprintf("CheckRole guildId %s, memderId %s, roleid %s \n", guildId, memderId, roleid))

		i := models.IncomingMessage{
			Text:    "%GuildMemberRemove",
			NameId:  memderId,
			GuildId: guildId,
			Type:    "ds",
		}

		_ = restapi.SendCompendiumApp(i)
		return false
	}
	for _, role := range member.Roles {
		if roleid == role {
			return true
		}
	}
	return false
}
func (d *Discord) GetMembersRoles(guildid string) (mm []models.DsMembersRoles) {
	members, err := d.S.GuildMembers(guildid, "", 1000)
	if err != nil {
		d.log.ErrorErr(err)
	}
	for _, member := range members {
		m := models.DsMembersRoles{
			Userid:  member.User.ID,
			RolesId: member.Roles,
		}
		mm = append(mm, m)
	}

	return mm
}
func (d *Discord) CleanRsBotOtherMessage() {
	defer func() {
		if r := recover(); r != nil {
			d.log.Info(fmt.Sprintf("recover() %+v", r))
		}
	}()
	for _, config := range d.corpConfigRS {
		if config.DsChannel != "" {
			channelMessages, err := d.S.ChannelMessages(config.DsChannel, 100, "", "", "")
			if err != nil {
				restErr, _ := err.(*discordgo.RESTError)
				if restErr.Message != nil && restErr.Message.Code == discordgo.ErrCodeUnknownChannel {
					d.log.Info("нужно сделать удаление этого канала : " + config.CorpName)
				} else {
					d.log.ErrorErr(err)
				}
				continue
			}
			if len(channelMessages) > 0 {
				t := time.Now().Unix()
				for _, message := range channelMessages {
					if message.Author.String() == "Rs_bot#9945" {
						if len(message.Embeds) > 0 {
							if IsDifferenceMoreThan5Minutes(message.Embeds[0].Timestamp) {
								d.log.Info(fmt.Sprintf("Rs_bot#9945 message.Embeds.Title: %+v\ndelete \n", message.Embeds[0].Title))
								_ = d.S.ChannelMessageDelete(message.ChannelID, message.ID)
							}
						} else if time.Now().Sub(message.Timestamp).Hours() < 96 && !strings.Contains(message.Content, "Черга") {
							d.log.Info(fmt.Sprintf("message hours%.1f %+v\n", time.Now().Sub(message.Timestamp).Hours(), message))
						} else {
							fmt.Printf("MESSAGE: %+v\n", message)
						}

					}
					if message.Author.String() != "RsBot#0000" && message.Author.String() != "Rs_bot#9945" && message.Author.String() != "КзБот#0000" {

						if t-message.Timestamp.Unix() < 1209600 && t-message.Timestamp.Unix() > 180 {
							if message.Content == "" || !strings.HasPrefix(message.Content, ".") {
								_ = d.S.ChannelMessageDelete(config.DsChannel, message.ID)
							}
						}
						if t-message.Timestamp.Unix() < 1209600 && t-message.Timestamp.Unix() > 260000 {
							if strings.HasPrefix(message.Content, ".") {
								_ = d.S.ChannelMessageDelete(config.DsChannel, message.ID)

							}
						}
					}
				}
				fmt.Println("clean OK " + config.CorpName)
			}
		}
	}
	fmt.Println("clean OK")
}
func IsDifferenceMoreThan5Minutes(timeStr string) bool {
	// Получаем текущее время
	currentTime := time.Now()

	// Парсинг переданного времени из строкового представления
	parsedTime, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		return false
	}

	// Вычисление разницы во времени
	diff := currentTime.Sub(parsedTime)
	if diff < 0 {
		diff = -diff
	}

	// Проверка, превышает ли разница 5 минут
	if diff > 5*time.Minute {
		return true
	}

	return false
}

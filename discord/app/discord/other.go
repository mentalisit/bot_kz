package DiscordClient

import (
	"discord/models"
	"fmt"
	gt "github.com/bas24/googletranslatefree"
	"github.com/bwmarrin/discordgo"
	"log/slog"
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
	exist, role := d.roleExists(d.re.GetGuildRoles(guildid), rolePing)
	if !exist {
		//создаем роль и возврашаем пинг
		newRole, err := d.createRole(rolePing, guildid)
		if err != nil {
			d.log.ErrorErr(err)
			return rolePing, err
		}
		return newRole.Mention(), nil
	} else {
		return role.Mention(), nil
	}
}

func (d *Discord) CleanChat(chatid, mesid, text string) {
	res := strings.HasPrefix(text, ".")
	if !res { //если нет префикса  то удалить через 3 минуты
		go d.DeleteMesageSecond(chatid, mesid, 180)
	}
}

func (d *Discord) roleExists(Roles []*discordgo.Role, nameRoles string) (bool, *discordgo.Role) {
	nameRoles = strings.ToLower(nameRoles)

	for _, role := range Roles {
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
	d.re.guildsRoles[guildid], _ = d.S.GuildRoles(guildid)
	return create, nil
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
			if strings.HasPrefix(message.Content, "Hades' Star Official") {
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

// for compendium
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

// for compendium
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

		d.api.SendCompendiumAppRecover(i)
		return false
	}
	for _, role := range member.Roles {
		if roleid == role {
			return true
		}
	}
	return false
}

// for compendium
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
	configRs, err := d.storage.Db.ReadConfigRs()
	if err != nil {
		slog.Error(err.Error())
		return
	}
	for _, config := range configRs {
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
						if strings.HasPrefix(message.Content, "Hades' Star Official") {
							continue
						}
						if len(message.Embeds) > 0 {
							if IsDifferenceMoreThan5Minutes(message.Embeds[0].Timestamp) {
								d.log.Info(fmt.Sprintf("Rs_bot#9945 message.Embeds.Title: %+v\ndelete \n", message.Embeds[0].Title))
								_ = d.S.ChannelMessageDelete(message.ChannelID, message.ID)
							}
						} else if time.Now().Sub(message.Timestamp).Hours() < 96 && !strings.Contains(message.Content, "Черга") {
							d.log.Info(fmt.Sprintf("message hours%.1f %+v\n", time.Now().Sub(message.Timestamp).Hours(), message))
						} else if len(message.Attachments) > 0 {
							fmt.Printf("удалено сообщение с вложением %+v\n", message.Attachments[0])
							_ = d.S.ChannelMessageDelete(message.ChannelID, message.ID)
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

func (d *Discord) ReadNewsChannel() (en, ru, ua string) {
	// Определяем временную метку 5 минут назад
	fiveMinutesAgo := time.Now().UTC().Add(-5 * time.Minute)

	messages, _ := d.S.ChannelMessages("1305333971269324851", 5, "", "", "")
	for i, message := range messages {
		if message.Timestamp.After(fiveMinutesAgo) {
			if message.Content != "" {
				fmt.Printf("Now message %d Timestamp %s %s %s\n", i, message.Timestamp, message.Author, message.Content)
			} else if len(message.Embeds) > 0 {
				fmt.Printf("Now message %d Timestamp %s %s %s\n", i, message.Embeds[0].Timestamp, message.Embeds[0].Author, message.Embeds[0].Description)
			}

			if message.MessageReference != nil && message.MessageReference.GuildID == "355101373483712513" {
				d.log.InfoStruct("news message ", message)
			}

			if message.Author.String() == "Hades' Star Official #announcements#0000" {
				if strings.Contains(message.Content, "Red Star event starts") {
					d.log.Info(message.Content)
					d.storage.Db.SaveEventDate(message.Content)
				}
				en = message.Content
				ru, _ = gt.Translate(message.Content, "auto", "ru")
				ua, _ = gt.Translate(message.Content, "auto", "uk")
				return
			}
			if len(message.Embeds) == 1 {
				return d.filterNewsMessage(message.Embeds[0].Description)
			}
		}
	}
	return
}
func (d *Discord) filterNewsMessage(msg string) (en, ru, ua string) {
	//message = `Season ${season} of the Corporation Red Star event has just started! Help your Corporation climb the event leaderboard throughout the weekend, and earn free Crystals in the end. For more information, see the in game Leaderboards window.`;
	reRsEvent := regexp.MustCompile(`Season (\d+) of the Corporation Red Star event has just started! Help your Corporation climb the event leaderboard throughout the weekend, and earn free Crystals in the end. For more information, see the in game Leaderboards window.`)
	match := reRsEvent.FindStringSubmatch(msg)
	if len(match) > 1 {
		season, _ := strconv.Atoi(match[1])
		d.log.Info(msg)
		d.storage.Db.SaveEventDate(msg)
		en = msg
		ru = fmt.Sprintf("Сезон %d события Корпорации Красная Звезда только что начался! Помогите своей Корпорации подняться в таблице лидеров события в течение выходных и получите бесплатные Кристаллы в конце. Для получения дополнительной информации см. игровое окно Таблицы лидеров.", season)
		ua = fmt.Sprintf("%d-й сезон події «Червона Зірка корпорації» щойно розпочався! Допоможіть своїй корпорації піднятися в таблиці лідерів події протягом вихідних і заробляйте безкоштовні кристали в кінці. Для отримання додаткової інформації дивіться вікно таблиці лідерів у грі.", season)
		return
	}
	if msg == "White Star event is on now! For all White Stars that start in the next 4 days, your Corporation will be awarded significantly more XP." {
		en = msg
		ru = "Событие «Белая звезда» уже началось! За все Белые Звезды, которые начнутся в течение следующих 4 дней, ваша Корпорация получит значительно больше опыта."
		ua = "Подія «Біла Зірка» вже розпочалася! За всі Білі Зірки, які розпочнуться протягом наступних 4 днів, ваша Корпорація отримає значно більше досвіду."
		return
	}
	if msg == "2x Credit Asteroid event is on now! For the next 3 days, all Rich Asteroid Fields in Red Stars will yield twice the credits." {
		en = msg
		ru = "Специальное мероприятие богатых астероидов красных звёзд активно! \nВ течение следующих 3 дней все богатые астероидные поля в красных звездах дадут вдвое больше кредитов."
		ua = "Спеціальний захід багатих астероїдів червоних зірок активно!\nПротягом наступних 3 днів усі багаті астероїдні поля у червоних зірках дадуть удвічі більше кредитів."
		return
	}
	if msg == "Blue Star special event is active now! For the next 3 days, all Blue Star credit rewards are doubled." {
		en = msg
		ru = "Специальное мероприятие голубых звёзд активно!\nВ течение следующих 3 дней награды - кредиты в мероприятии голубых звёзд  удваиваются."
		ua = "Спеціальний захід блакитних зірок активно!\nПротягом наступних 3 днів нагороди – кредити у заході блакитних зірок подвоюються."
		return
	}
	d.log.Info(msg)
	return "", "", ""
}

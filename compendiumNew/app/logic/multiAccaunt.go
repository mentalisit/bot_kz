package logic

import (
	"compendium/logic/generate"
	"compendium/models"
	"fmt"
	"github.com/google/uuid"
	"strings"
	"time"
)

var dm map[string]map[string]string

func (c *Hs) multiConnect(m models.IncomingMessage) bool {
	if m.Text == "%создать мульти аккаунт" || m.Text == "%create multi account" {
		findMultiAccount, _ := c.db.Multi.FindMultiAccountByUserId(m.NameId)
		if findMultiAccount != nil {
			c.sendChat(m, "у тебя уже есть мультиаккаунт\n "+findMultiAccount.GetTextUsername()) //+данные аккаунта
			return true
		}
		multiAccount, err := c.db.Multi.CreateMultiAccountWithPlatform(m.NameId, m.Name, m.Type, m.Name)
		if err != nil {
			c.log.ErrorErr(err)
		}
		if multiAccount != nil {
			c.sendChat(m, "создан мультиаккаунт для подключения другого мессенджера отправь команду %соединить ")
			//todo c.db.DB.MultiListUserInsert(generate.GenerateToken(180), m.GuildId, multiAccount.UUID)
		}
		return true
	}
	if m.Text == "%соединить" {
		defer func() {
			if r := recover(); r != nil {
				c.log.Info(fmt.Sprintf("recover() %+v", r))
			}
		}()

		fMA, err := c.db.Multi.FindMultiAccountByUserId(m.NameId)
		if err != nil {
			c.log.ErrorErr(err)
		}
		if fMA != nil && fMA.TelegramID != "" && fMA.DiscordID != "" {
			accountCorpMember, _ := c.db.Multi.CorpMemberByUId(fMA.UUID)
			for _, gid := range accountCorpMember.GuildIds {
				c.newLink(m, gid)
				time.Sleep(1 * time.Second)
			}
			return true
		}
		if fMA == nil {
			fMA, err = c.db.Multi.CreateMultiAccountWithPlatform(m.NameId, m.Name, m.Type, m.Name)
			if err != nil {
				c.log.ErrorErr(err)
				return true
			}
		}
		if dm == nil {
			dm = make(map[string]map[string]string)
		}

		if dm[fMA.UUID.String()] == nil {
			dm[fMA.UUID.String()] = make(map[string]string)
		}
		dm[fMA.UUID.String()][m.Type] = m.DmChat

		linkAccount, err := c.db.Multi.GenerateLinkCode(fMA.UUID)
		if err != nil {
			c.log.ErrorErr(err)
			return true
		}
		if linkAccount != nil {
			c.sendChat(m, "проверь личные сообщения ")
			text := "Отправь мне в личные сообщения "
			if m.Type == "ds" {
				text = text + "в телеграм боте https://t.me/gote1st_bot"
			} else if m.Type == "tg" {
				text = text + "в дискорд боте https://discord.com/users/909526127305953290"
			}
			_, err = c.sendDM(m, text)
			if err != nil && err.Error() == "forbidden" {
				if c.checkMoron(m) {
					c.log.InfoStruct("moron", m)
				} else {
					c.sendChat(m, fmt.Sprintf(c.getText(m, "ERROR_SEND"), m.MentionName))
				}
				return true
			}
			_, err = c.sendDM(m, "%код соединения "+linkAccount.Code)
			if err != nil {
				c.log.ErrorErr(err)
			}
		}
		return true
	}
	cutPrefix, found := strings.CutPrefix(m.Text, "%код соединения ")
	if found {

		m.Type = m.Type[:2]
		uid, err := c.db.Multi.ValidateLinkCode(cutPrefix)
		if err != nil {
			c.log.ErrorErr(err)
			_, _ = c.sendDM(m, err.Error())
			return true
		}
		if uid != nil {
			err = c.db.Multi.DeleteLinkCodesByUUID(*uid)
			if err != nil {
				c.log.ErrorErr(err)
			}
			multiAccount, err := c.db.Multi.FindMultiAccountUUID(*uid)
			if err != nil {
				c.log.ErrorErr(err)
			}
			if multiAccount != nil {
				if dm == nil {
					dm = make(map[string]map[string]string)
				}

				if dm[multiAccount.UUID.String()] == nil {
					dm[multiAccount.UUID.String()] = make(map[string]string)
				}
				dm[multiAccount.UUID.String()][m.Type] = m.DmChat
				text := ""
				if (m.Type == "ds" && multiAccount.DiscordID == "") || (m.Type == "tg" && multiAccount.TelegramID == "") {
					multiAccount, err = c.db.Multi.UpdateMultiAccountInfo(multiAccount.UUID, m.Type, m.NameId, m.Name)
					if err != nil {
						c.log.ErrorErr(err)
						return true
					}

					text = fmt.Sprintf("соединено удачно.\n%s", multiAccount.GetTextUsername())
					m.MultiAccount = multiAccount
					c.migrationDataMultiUser(m)

				} else {
					text = "уже соединено \n" + multiAccount.GetTextUsername()
					c.log.Info(fmt.Sprintf("tip %s multiAccount %+v", m.Type, multiAccount))
				}
				_, _ = c.sendDM(m, text)
				//todo need send url + new token

			}
		}

		return true
	}
	if m.Text == "%миграция тест" {
		c.migrationDataMultiUser(m)
		return true
	}
	return false
}
func (c *Hs) newLink(m models.IncomingMessage, Gid uuid.UUID) {
	if m.MultiAccount != nil {
		jwtGenerateToken, errJwt := generate.JWTGenerateToken(m.MultiAccount.UUID, Gid, m.Name)
		if errJwt != nil {
			c.log.ErrorErr(errJwt)
		}
		getUUID, _ := c.db.Multi.GuildGetUUID(Gid)
		if getUUID != nil {
			sendDM, _ := c.sendDM(m, "Подготавливаю секретную ссылку")
			time.Sleep(2 * time.Second)

			links := "https://mentalisit.github.io/HadesSpace/compendiumTech?secretToken=" + jwtGenerateToken
			text := fmt.Sprintf(c.getText(m, "SECRET_LINK"), links, getUUID.GuildName)
			err := c.editMessage(m, m.DmChat, sendDM, text, "MarkdownV2")
			if err != nil {
				c.log.ErrorErr(err)
				return
			}
		}
	}
}
func (c *Hs) migrationDataMultiUser(m models.IncomingMessage) {
	defer func() {
		if r := recover(); r != nil {
			c.log.Info(fmt.Sprintf("recover() %+v", r))
		}
	}()
	//user
	fmt.Println("user")
	var alts map[string]int
	alts = make(map[string]int)
	var gameName string
	getUser := func(user *models.User) {
		if user == nil {
			return
		}
		if user.GameName != "" {
			gameName = user.GameName
		}
		if len(user.Alts) > 0 {
			for _, alt := range user.Alts {
				alts[alt]++
			}
		}
	}
	byUserIdTg, _ := c.users.UsersGetByUserId(m.MultiAccount.TelegramID)
	getUser(byUserIdTg)
	byUserIdDs, _ := c.users.UsersGetByUserId(m.MultiAccount.DiscordID)
	getUser(byUserIdDs)
	if gameName != "" && gameName != m.MultiAccount.Nickname {
		m.MultiAccount.Nickname = gameName
		_, err := c.db.Multi.UpdateMultiAccountNickname(*m.MultiAccount)
		if err != nil {
			c.log.ErrorErr(err)
		}
		fmt.Printf("UpdateMultiAccountNickname %s\n", m.MultiAccount.Nickname)
	}
	if len(alts) > 0 {
		for s, _ := range alts {
			m.MultiAccount.Alts = append(m.MultiAccount.Alts, s)
		}
		_, err := c.db.Multi.UpdateMultiAccountAlts(*m.MultiAccount)
		if err != nil {
			c.log.ErrorErr(err)
		}
		fmt.Printf("UpdateMultiAccountAlts %+v\n", m.MultiAccount.Alts)
	}
	AltContains := func(s string) {
		contains := false
		for _, alt := range m.MultiAccount.Alts {
			if alt == s {
				contains = true
			}
		}
		if !contains {
			m.MultiAccount.Alts = append(m.MultiAccount.Alts, s)
			c.db.Multi.UpdateMultiAccountAlts(*m.MultiAccount)
		}
	}
	//user end

	fmt.Println("member ")
	//member
	member := models.MultiAccountCorpMember{
		Uid: m.MultiAccount.UUID,
	}
	GetMembers := func(userId string) {
		members, err := c.corpMember.CorpMembersReadByUserId(userId)
		if err != nil {
			c.log.ErrorErr(err)
		}
		if len(members) > 0 {
			for _, i := range members {
				if i.AfkFor != "" {
					member.AfkFor = i.AfkFor
				}
				if i.TimeZone != "" {
					member.TimeZona = i.TimeZone
				}
				if i.ZoneOffset != 0 {
					member.ZonaOffset = i.ZoneOffset
				}
				getGuild, _ := c.db.Multi.GuildGet(i.GuildId)
				if getGuild == nil {
					c.log.Error("getGuild == nil")
					return
				}
				found := false
				for _, id := range member.GuildIds {
					if id == getGuild.GId {
						found = true
					}
				}
				if !found {
					member.GuildIds = append(member.GuildIds, getGuild.GId)
				}

			}
		}
	}
	if m.MultiAccount.TelegramID != "" {
		GetMembers(m.MultiAccount.TelegramID)
	}
	if m.MultiAccount.DiscordID != "" {
		GetMembers(m.MultiAccount.DiscordID)
	}
	for _, id := range member.GuildIds {
		c.newLink(m, id)
		time.Sleep(1 * time.Second)
	}
	err := c.db.Multi.CorpMemberInsert(member)
	if err != nil {
		c.log.ErrorErr(err)
	}
	fmt.Printf("CorpMemberInsert %+v\n", member)
	//member end

	fmt.Println("Tech")
	//Tech
	var tech map[string][]models.TechTable
	tech = make(map[string][]models.TechTable)
	getTech := func(userid string) {
		techTables, _ := c.tech.TechGetAllUserId(userid)
		for _, table := range techTables {
			if m.Name == table.Name {
				table.Name = m.MultiAccount.Nickname
			}
			tech[table.Name] = append(tech[table.Name], table)
		}
	}
	if m.MultiAccount.TelegramID != "" {
		getTech(m.MultiAccount.TelegramID)
	}
	if m.MultiAccount.DiscordID != "" {
		getTech(m.MultiAccount.DiscordID)
	}
	for _, tables := range tech {
		var techOne map[string]models.TechTable
		techOne = make(map[string]models.TechTable)
		var t models.TechTable
		if len(tables) == 1 {
			t = tables[0]
		} else {
			for _, table := range tables {
				if techOne[table.Name].Name == "" {
					techOne[table.Name] = table
				} else {
					if len(techOne[table.Name].Tech) < len(table.Tech) {
						techOne[table.Name] = table
					}
				}
			}
			for _, table := range techOne {
				t = table
			}
		}
		AltContains(t.Name)
		if t.Name == m.Name {
			t.Name = m.MultiAccount.Nickname
		}
		err = c.db.Multi.TechnologiesInsert(m.MultiAccount.UUID, t.Name, t.Tech)
		if err != nil {
			c.log.ErrorErr(err)
		}
		fmt.Printf("TechnologiesInsert %s %+v\n", t.Name, t.Tech)
	}
	//tech end
	fmt.Println("Tech end ")
	c.db.Multi.RemoveUserIdAllTable(m.MultiAccount)

	text := fmt.Sprintf("перенос данных на новую версию завершён \nВаши твины:%+v", m.MultiAccount.Alts)
	c.sendDM(m, text)

	if dm[m.MultiAccount.UUID.String()] != nil {
		for s, s2 := range dm[m.MultiAccount.UUID.String()] {
			if s == "ds" {
				fmt.Println(s2)
			}
		}
	}

}

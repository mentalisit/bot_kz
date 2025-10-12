package logic

import (
	"compendium/logic/generate"
	"compendium/models"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

const urlLink = "https://mentalisit.github.io/HadesSpace/"

func (c *Hs) connect(m models.IncomingMessage) (conn bool) {
	if helperCommand(m.Text, "connect") ||
		helperCommand(m.Text, "подключить") ||
		helperCommand(m.Text, "підключити") {
		conn = true
	}
	if !conn {
		return false
	}

	text := fmt.Sprintf(c.getText(m, "CODE_FOR_CONNECT"), m.MultiGuild.GuildName)
	mid1, err := c.sendDM(m, text)
	if err != nil && err.Error() == "forbidden" {
		if c.checkMoron(m) {
			c.log.InfoStruct("moron", m)
		} else {
			c.sendChat(m, fmt.Sprintf(c.getText(m, "ERROR_SEND"), m.MentionName))
			return
		}

	} else if err != nil {
		c.log.ErrorErr(err)
		c.log.InfoStruct("connect error ", err)
		return
	}

	c.sendChat(m, fmt.Sprintf(c.getText(m, "INSTRUCTIONS_SEND"), m.MentionName))

	newIdentify, cm := generate.GenerateIdentity(m)
	if m.MultiAccount != nil {
		jwtGenerateToken, errJwt := generate.JWTGenerateToken(m.MultiAccount.UUID, m.MultiGuild.GId, m.Name)
		if errJwt != nil {
			c.log.ErrorErr(errJwt)
			return false
		}
		newIdentify.Token = jwtGenerateToken

		sendDM, _ := c.sendDM(m, "Подготавливаю секретную ссылку")
		go func() {
			time.Sleep(5 * time.Second)
			_ = c.deleteMessage(m, m.DmChat, mid1)

			links := "https://mentalisit.github.io/HadesSpace/compendiumTech?secretToken=" + newIdentify.Token
			text = fmt.Sprintf(c.getText(m, "SECRET_LINK"), links, m.MultiGuild.GuildName)
			err = c.editMessage(m, m.DmChat, sendDM, text, "MarkdownV2")
			if err != nil {
				c.log.ErrorErr(err)
				return
			}
		}()
		memberByUId, _ := c.db.Multi.CorpMemberByUId(m.MultiAccount.UUID)
		if memberByUId == nil {
			_ = c.db.Multi.CorpMemberInsert(models.MultiAccountCorpMember{
				Uid:        m.MultiAccount.UUID,
				GuildIds:   []uuid.UUID{m.MultiGuild.GId},
				TimeZona:   "",
				ZonaOffset: 0,
				AfkFor:     "",
			})
			return
		}
		if len(memberByUId.GuildIds) > 0 {
			var found bool
			for _, id := range memberByUId.GuildIds {
				if id == m.MultiGuild.GId {
					found = true
				}
			}
			if !found {
				memberByUId.GuildIds = append(memberByUId.GuildIds, m.MultiGuild.GId)
				_ = c.db.Multi.CorpMemberUpdateGuildIds(*memberByUId)
			}
		}

	} else {
		//tokenOld, _ := c.listUser.ListUserGetToken(m.NameId, m.GuildId)
		//if tokenOld == "" {
		//	tokenOld, _ = c.listUser.ListUserGetToken(m.NameId, m.MultiGuild.GId.String())
		//}
		//
		//
		//prefixToken := fmt.Sprintf("%s.%s-mg.%s", m.Type, m.NameId, m.MultiGuild.GuildId())
		//if tokenOld != "" && strings.Contains(tokenOld, prefixToken) {
		//	newIdentify.Token = tokenOld
		//} else {
		//	newIdentify.Token = prefixToken + generate.GenerateToken(174)
		//	if tokenOld != "" {
		//		err = c.listUser.ListUserDelete(tokenOld)
		//		if err != nil {
		//			c.log.ErrorErr(err)
		//		}
		//	}
		//	err = c.listUser.ListUserInsert(newIdentify.Token, m.NameId, m.MultiGuild.GuildId())
		//	if err != nil {
		//		c.log.ErrorErr(err)
		//	}
		//}

		code := c.generateCodeAndSave(newIdentify)

		err = c.users.UsersInsert(newIdentify.User) //insert or update
		if err != nil {
			c.log.ErrorErr(err)
			return
		}

		err = c.corpMember.CorpMemberInsert(cm)
		if err != nil {
			c.log.ErrorErr(err)
			return
		}

		//getCountByGuildIdByUserId, _ := c.listUser.ListUserGetCountByGuildIdByUserId(m.GuildId, m.NameId)
		//if getCountByGuildIdByUserId != 0 {
		//	_ = c.listUser.ListUserDeleteByUserIdByGuildId(m.NameId, m.GuildId)
		//}
		//
		//err = c.listUser.ListUserInsert(newIdentify.Token, newIdentify.User.ID, newIdentify.MultiGuild.GuildId())
		//if err != nil {
		//	c.log.ErrorErr(err)
		//	return
		//}

		mid2, errs1 := c.sendDM(m, code)
		if errs1 != nil {
			c.log.ErrorErr(errs1)
			return
		}

		urlLinkAdd := "compendiumTech?c=" + code + "&lang=" + m.Language + "&client=1"
		mid3, errs2 := c.sendDM(m, fmt.Sprintf(c.getText(m, "PLEASE_PASTE_CODE"), urlLink, urlLink+urlLinkAdd))
		if errs2 != nil {
			c.log.ErrorErr(errs2)
			return
		}

		go c.timerEditMessage(m, mid1, mid2, mid3, newIdentify.Token)
	}
	return
}

func (c *Hs) generateCodeAndSave(Identity models.Identity) string {
	segments := []string{generate.RandString(4), generate.RandString(4), generate.RandString(4)}

	m := models.Code{
		Code:      strings.Join(segments, "-"),
		Timestamp: time.Now().Unix(),
		Identity:  Identity,
	}

	go func() {
		err := c.db.DB.CodeInsert(m)
		if err != nil {
			c.log.ErrorErr(err)
			c.log.InfoStruct("CodeInsert", m)
		}
	}()

	return m.Code
}

func (c *Hs) timerEditMessage(m models.IncomingMessage, mid1, mid2, mid3, token string) {
	time.Sleep(3 * time.Minute)

	//token, _ := c.listUser.ListUserGetToken(m.NameId, m.MultiGuild.GuildId())
	//if token == "" {
	//	c.log.InfoStruct("token is empty", m)
	//	return
	//}
	links := "https://mentalisit.github.io/HadesSpace/compendiumTech?secretToken=" + token
	text := fmt.Sprintf(c.getText(m, "SECRET_LINK"), links, m.MultiGuild.GuildName)
	if m.Type == "ds" || m.Type == "tg" {
		err := c.editMessage(m, m.DmChat, mid1, text, "MarkdownV2")
		if err != nil {
			c.log.ErrorErr(err)
			return
		}
	} else if m.Type == "wa" {
		c.sendDM(m, text)
	}
	err := c.deleteMessage(m, m.DmChat, mid2)
	if err != nil {
		c.log.ErrorErr(err)
		return
	}
	err = c.deleteMessage(m, m.DmChat, mid3)
	if err != nil {
		c.log.ErrorErr(err)
		return
	}
}

func (c *Hs) checkMoron(in models.IncomingMessage) bool {
	if len(c.moron) > 0 {
		if c.moron[in] != 0 {
			return true
		}
	}
	c.moron[in] += 1
	return false
}

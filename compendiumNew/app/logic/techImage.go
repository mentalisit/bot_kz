package logic

import (
	"compendium/logic/imageGenerator"
	"compendium/models"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"regexp"
	"strings"
)

func (c *Hs) BytesToTechLevel(b []byte) (map[int]models.TechLevel, models.TechLevelArray) {
	var m map[int]models.TechLevel
	m = make(map[int]models.TechLevel)
	err := json.Unmarshal(b, &m)
	if err != nil {
		fmt.Println(err)
		m[701] = models.TechLevel{
			Ts:    0,
			Level: 0,
		}
	}
	var mi = make(models.TechLevelArray)
	for i, le := range m {
		mi[i] = [2]int{le.Level}
	}
	return m, mi
}
func (c *Hs) sendPic(m models.IncomingMessage, userAvatar, guildAvatar, picName, guildName string, mBytes []byte) {
	_, mt := c.BytesToTechLevel(mBytes)
	userPic := imageGenerator.GenerateUser(userAvatar, guildAvatar, picName, guildName, mt)
	c.sendChatPic(m, "", userPic)
}
func (c *Hs) techImage(m models.IncomingMessage) (tech bool) {
	if m.Text == "%т и" || m.Text == "%t i" || m.Text == "%т і" || m.Text == "%tech" || m.Text == "%техи" {
		tech = true
	}
	if !tech {
		return false
	}

	var mBytes []byte
	picName := m.Name
	guildName := m.GuildName
	guildAvatar := m.GuildAvatar
	userAvatar := m.Avatar

	if m.MultiGuild != nil {
		guildName = m.MultiGuild.GuildName
		guildAvatar = m.MultiGuild.AvatarUrl
	}

	if m.MultiAccount != nil {
		mBytesTech, err := c.db.Multi.TechnologiesGet(m.MultiAccount.UUID, m.MultiAccount.Nickname)
		if err != nil {
			c.log.ErrorErr(err)
			c.sendChat(m, c.getText(m, "DATA_NOT_FOUND"))
			return
		}

		if mBytesTech != nil && len(mBytesTech) > 0 {
			mBytes = mBytesTech
			picName = m.MultiAccount.Nickname
		}
	}

	if len(mBytes) == 0 {
		mBytesTech, err := c.tech.TechGet(m.Name, m.NameId, m.GuildId)
		if err != nil {
			c.log.ErrorErr(err)
			c.sendChat(m, c.getText(m, "DATA_NOT_FOUND"))
			return
		}
		mBytes = mBytesTech

		user, err := c.users.UsersGetByUserId(m.NameId)
		if err != nil {
			if !errors.Is(err, pgx.ErrNoRows) {
				c.log.ErrorErr(err)
				c.log.InfoStruct("techImage", m)
			}
			c.sendChat(m, c.getText(m, "DATA_NOT_FOUND"))
			return
		}

		if user != nil && user.GameName != "" {
			picName = user.GameName
		}
	}

	if len(mBytes) == 0 {
		c.sendChat(m, c.getText(m, "DATA_NOT_FOUND"))
		return true
	}
	c.sendPic(m, userAvatar, guildAvatar, picName, guildName, mBytes)
	return
}

func (c *Hs) techImageName(m models.IncomingMessage) bool {
	after, _ := strings.CutPrefix(m.Text, "%")
	re := regexp.MustCompile(`^[тt] +<@(\d{17,20})> +[иiі]$`)
	retg := regexp.MustCompile(`^[tт] +@(\S+) +[iиі]$`)
	reNew := regexp.MustCompile(`^(?:tech|техи) +(?:<@(\d{17,20})>|@(\S+))$`)

	matches := re.FindStringSubmatch(after)
	matchestg := retg.FindStringSubmatch(after)
	matchesNew := reNew.FindStringSubmatch(after)

	if len(matches) > 0 || len(matchestg) > 0 || len(matchesNew) > 0 {
		var mBytes []byte
		picName := ""
		guildName := m.GuildName
		guildAvatar := m.GuildAvatar
		userAvatar := ""

		if m.MultiGuild != nil {
			guildName = m.MultiGuild.GuildName
			guildAvatar = m.MultiGuild.AvatarUrl
		}

		var err error
		userID := ""
		userName := ""
		if len(matches) > 0 { // id discord
			userID = matches[1]
		}
		if len(matchestg) > 0 {
			userName = matchestg[1]
		}
		if len(matchesNew) > 0 {
			if matchesNew[1] != "" {
				userID = matchesNew[1] // Если Discord ID найден, он здесь
			} else if matchesNew[2] != "" {
				userName = matchesNew[2] // Если Telegram username найден, он здесь
			}
		}

		var multiAcc *models.MultiAccount
		if userID != "" {
			multiAcc, _ = c.db.Multi.FindMultiAccountByUserId(userID)
		} else if userName != "" {
			multiAcc, _ = c.db.Multi.FindMultiAccountByUsername(userName)
		}

		if multiAcc != nil {
			mBytesTech, _, err := c.db.Multi.TechnologiesGetName(multiAcc.Nickname)
			if err != nil {
				return false
			}
			if len(mBytesTech) != 0 {
				mBytes = mBytesTech
				picName = multiAcc.Nickname
				userAvatar = multiAcc.AvatarURL
			}

		}
		if len(mBytes) == 0 {
			var user *models.User
			if userID != "" {
				user, err = c.users.UsersGetByUserId(userID)
			} else if userName != "" {
				user, err = c.users.UsersGetByUserName(userName)
			}
			if err != nil || user == nil {
				if !errors.Is(err, pgx.ErrNoRows) {
					c.log.Info(err.Error())
				}
				c.sendChat(m, c.getText(m, "DATA_NOT_FOUND"))
				return true
			}
			techBytes, err := c.tech.TechGet(user.Username, user.ID, m.GuildId)
			if err != nil {
				c.log.ErrorErr(err)
				c.sendChat(m, c.getText(m, "DATA_NOT_FOUND"))
				return true
			}
			mBytes = techBytes
			picName = user.Username
			userAvatar = user.AvatarURL
			if user.GameName != "" {
				picName = user.GameName
			}

		}
		if len(mBytes) == 0 {
			c.sendChat(m, c.getText(m, "DATA_NOT_FOUND"))
			return true
		}
		c.sendPic(m, userAvatar, guildAvatar, picName, guildName, mBytes)
		return true
	}
	return false
}

func (c *Hs) techImageNameAlt(m models.IncomingMessage) bool {
	after, _ := strings.CutPrefix(m.Text, "%")

	reNew := regexp.MustCompile(`^(?:[tт] +(\S+(?: +\S+)?) +[iиі]|tech +(\S+(?: +\S+)?)|техи +(\S+(?: +\S+)?))$`)

	userName := ""
	match := reNew.FindStringSubmatch(after)

	if len(match) > 1 {
		// Поиск имени в одной из групп
		for i := 1; i < len(match); i++ {
			if match[i] != "" {
				userName = match[i]
				break
			}
		}
	} else {
		re := regexp.MustCompile(`^[tт] +(\S+(?: +\S+)?) +[iиі]$`)
		matches := re.FindStringSubmatch(after)
		if len(matches) > 0 {
			userName = matches[1]
		}
	}

	if userName != "" {
		var mBytes []byte
		picName := userName
		guildName := m.GuildName
		guildAvatar := m.GuildAvatar
		userAvatar := ""

		if m.MultiGuild != nil {
			guildName = m.MultiGuild.GuildName
			guildAvatar = m.MultiGuild.AvatarUrl
		}

		multiAccount, _ := c.db.Multi.FindMultiAccountByUsername(userName)
		if multiAccount != nil {
			userAvatar = multiAccount.AvatarURL
			mBytesTech, _ := c.db.Multi.TechnologiesGet(multiAccount.UUID, userName)
			if mBytesTech != nil {
				mBytes = mBytesTech
			}
		}
		if multiAccount == nil {
			techBytes, userID, err := c.tech.TechGetName(userName, m.GuildId)
			if err != nil || userID == "" {
				usersGetNick, _ := c.users.UsersFindByGameName(userName)
				if len(usersGetNick) > 1 {
					c.log.InfoStruct("techImage usersGetNick", usersGetNick)
				}
				if usersGetNick != nil && len(usersGetNick) > 0 {
					for _, userGetNick := range usersGetNick {
						techBytes, userID, err = c.tech.TechGetName(userGetNick.Username, m.GuildId)
						if userID != "" || techBytes != nil {
							break
						}
					}
				}
				if userID == "" || techBytes == nil {
					compatible, _ := c.listOfCompatible(&models.Guild{ID: m.GuildId, Type: m.Type, Name: m.GuildName})
					if compatible != nil {
						techBytes, userID, _ = c.tech.TechGetName(userName, compatible.ID)
					}
					if userID == "" || techBytes == nil {
						c.sendChat(m, c.getText(m, "DATA_NOT_FOUND"))
						return true
					}
				}
			}

			user, err := c.users.UsersGetByUserId(userID)
			if err != nil || user == nil {
				c.sendChat(m, c.getText(m, "DATA_NOT_FOUND"))
				c.log.ErrorErr(err)
				return true
			}
			userAvatar = user.AvatarURL
			mBytes = techBytes
		}
		if len(mBytes) == 0 {
			c.sendChat(m, c.getText(m, "DATA_NOT_FOUND"))
			return true
		}
		c.sendPic(m, userAvatar, guildAvatar, picName, guildName, mBytes)
		return true
	}

	return false
}

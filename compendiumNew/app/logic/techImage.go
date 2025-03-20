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
func (c *Hs) techImage(m models.IncomingMessage) (tech bool) {
	if m.Text == "%т и" || m.Text == "%t i" || m.Text == "%т і" || m.Text == "%tech" || m.Text == "%техи" {
		tech = true
	}
	if !tech {
		return false
	}
	mBytes, err := c.tech.TechGet(m.Name, m.NameId, m.GuildId)
	if err != nil {
		c.log.ErrorErr(err)
		c.sendChat(m, c.getText(m, "DATA_NOT_FOUND"))
		return
	}
	user, err := c.users.UsersGetByUserId(m.NameId)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			c.log.ErrorErr(err)
			c.log.InfoStruct("techImage", m)
		}
		c.sendChat(m, c.getText(m, "DATA_NOT_FOUND"))
		return
	}
	picName := m.Name
	if user != nil && user.GameName != "" {
		picName = user.GameName
	}
	_, mt := c.BytesToTechLevel(mBytes)
	userPic := imageGenerator.GenerateUser(
		m.Avatar,
		m.GuildAvatar,
		picName,
		m.GuildName,
		mt)
	c.sendChatPic(m, "", userPic)
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
		var user *models.User
		var err error
		if len(matches) > 0 {
			userID := matches[1]
			user, err = c.users.UsersGetByUserId(userID)
		} else if len(matchestg) > 0 {
			name := matchestg[1]
			user, err = c.users.UsersGetByUserName(name)
		} else if len(matchesNew) > 0 {
			discordID := matchesNew[1]  // Если Discord ID найден, он здесь
			tgUsername := matchesNew[2] // Если Telegram username найден, он здесь
			if discordID != "" {
				user, err = c.users.UsersGetByUserId(discordID)
			} else if tgUsername != "" {
				user, err = c.users.UsersGetByUserName(tgUsername)
			}
		}
		if err != nil {
			if !errors.Is(err, pgx.ErrNoRows) {
				c.log.Info(err.Error())
			}
			c.sendChat(m, c.getText(m, "DATA_NOT_FOUND"))
			return true
		}
		if user == nil {
			c.sendChat(m, c.getText(m, "DATA_NOT_FOUND"))
			return true
		}
		techBytes, err := c.tech.TechGet(user.Username, user.ID, m.GuildId)
		if err != nil {
			c.log.ErrorErr(err)
			c.sendChat(m, c.getText(m, "DATA_NOT_FOUND"))
			return true
		}
		picName := user.Username
		if user.GameName != "" {
			picName = user.GameName
		}
		_, mt := c.BytesToTechLevel(techBytes)
		userPic := imageGenerator.GenerateUser(
			user.AvatarURL,
			m.GuildAvatar,
			picName,
			m.GuildName,
			mt)
		c.sendChatPic(m, "", userPic)
		return true
	}
	return false
}
func (c *Hs) techImageNameAlt(m models.IncomingMessage) bool {
	after, _ := strings.CutPrefix(m.Text, "%")
	re := regexp.MustCompile(`^[tт] +(\S+(?: +\S+)?) +[iиі]$`)
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
	}

	matches := re.FindStringSubmatch(after)
	if len(matches) > 0 {
		userName = matches[1]
	}

	if userName != "" {
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
		if user.Username == userName && user.GameName != "" {
			userName = user.GameName
		}

		_, mt := c.BytesToTechLevel(techBytes)
		userPic := imageGenerator.GenerateUser(
			user.AvatarURL,
			m.GuildAvatar,
			userName,
			m.GuildName,
			mt)
		c.sendChatPic(m, "", userPic)
		return true
	}

	return false
}

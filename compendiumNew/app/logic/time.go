package logic

import (
	"compendium/models"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func (c *Hs) TzTime(m models.IncomingMessage) bool {
	if c.TzTimeSet(m) {
		return true
	} else if c.TzGet(m) {
		return true
	} else if c.TzGetTime(m) {
		return true
	}
	return false
}
func (c *Hs) TzTimeSet(m models.IncomingMessage) bool {
	// Регулярное выражение для поиска подстроки "%tz set", необязательного параметра и часового пояса или смещения
	re := regexp.MustCompile(`%tz set\s*(<@\d+>|@\S+)?\s*([+-]?\d+(\.\d+)?|[A-Za-z/_]+)`)

	// Ищем совпадения в строке
	matches := re.FindStringSubmatch(m.Text)
	if len(matches) < 3 {
		return false
	}

	optionalParam := matches[1]
	timezone := matches[2]

	// Проверяем, если это числовое смещение с дробной частью
	if offset, err := strconv.ParseFloat(timezone, 64); err == nil {
		c.TzTimeSetTime(offset, optionalParam, m)
		return true
	}

	// Пытаемся загрузить местоположение
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		text := fmt.Sprintf(c.getText(m, "I_COULD_NOT_FIND_ANY"), m.MentionName, timezone)
		c.sendChat(m, text)
	}

	// Получаем текущее время в указанном часовом поясе
	now := time.Now().In(loc)
	_, offset := now.Zone()

	// Преобразуем смещение из секунд в часы с дробной частью
	offsetHours := float64(offset) / 3600
	c.TzTimeSetTime(offsetHours, optionalParam, m)
	return true
}
func (c *Hs) TzTimeSetTime(offset float64, mentionName string, m models.IncomingMessage) {
	offsetInt := int(offset * 60)
	timeZona := fmt.Sprintf("UTC+%+v", offset)
	if offset < 0 {
		timeZona = fmt.Sprintf("UTC%+v", offset)
	}

	cm := models.CorpMember{
		Name:       m.Name,
		GuildId:    m.GuildId,
		Avatar:     m.AvatarF,
		Tech:       map[int][2]int{},
		AvatarUrl:  m.Avatar,
		TimeZone:   timeZona,
		ZoneOffset: offsetInt,
	}
	u := models.User{
		Username:  m.Name,
		AvatarURL: m.Avatar,
		Alts:      []string{},
	}

	if mentionName == "" {
		err := c.corpMember.CorpMemberTZUpdate(m.NameId, m.GuildId, timeZona, offsetInt)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				cm.UserId = m.NameId
				_ = c.corpMember.CorpMemberInsert(cm)
				u.ID = m.NameId
				_ = c.users.UsersInsert(u)
			} else {
				c.log.ErrorErr(err)
				return
			}
		}
		text := fmt.Sprintf(c.getText(m, "TIMEZONA_SET"), m.MentionName, m.Name, timeZona)
		c.sendChat(m, text)
	} else {
		re := regexp.MustCompile(`<@(\d+)>|@(\S+)`)

		// Ищем совпадения в строке
		matches := re.FindStringSubmatch(mentionName)
		text := ""
		if matches[1] != "" {
			//ds nameid
			err := c.corpMember.CorpMemberTZUpdate(matches[1], m.GuildId, timeZona, offsetInt)
			if err != nil {
				if errors.Is(err, pgx.ErrNoRows) {
					cm.UserId = matches[1]
					_ = c.corpMember.CorpMemberInsert(cm)
					u.ID = matches[1]
					_ = c.users.UsersInsert(u)
				} else {
					c.log.ErrorErr(err)
					return
				}
			}
			text = fmt.Sprintf(c.getText(m, "TIMEZONA_SET"), m.MentionName, mentionName, timeZona)
		} else if matches[2] != "" {
			//tg name
			user, err := c.users.UsersGetByUserName(matches[2])
			if err != nil {
				c.log.ErrorErr(err)
				return
			}
			err = c.corpMember.CorpMemberTZUpdate(user.ID, m.GuildId, timeZona, offsetInt)
			if err != nil {
				if errors.Is(err, pgx.ErrNoRows) {
					cm.UserId = user.ID
					_ = c.corpMember.CorpMemberInsert(cm)
					u.ID = user.ID
					_ = c.users.UsersInsert(u)
				} else {
					c.log.ErrorErr(err)
					return
				}
			}
			text = fmt.Sprintf(c.getText(m, "TIMEZONA_SET"), m.MentionName, matches[2], timeZona)
		}
		if text != "" {
			c.sendChat(m, text)
		}
	}
}
func (c *Hs) TzGet(m models.IncomingMessage) bool {
	if strings.HasPrefix(m.Text, "%tz get") {
		members, err := c.corpMember.CorpMembersRead(m.GuildId)
		if err != nil {
			c.log.ErrorErr(err)
			c.sendChat(m, fmt.Sprintf(c.getText(m, "TIMEZONA_IS_CURRENTLY"), m.MentionName, m.Name, "Not set"))
			return false
		}
		for _, member := range members {
			if member.UserId == m.NameId {
				text := fmt.Sprintf(c.getText(m, "TIMEZONA_IS_CURRENTLY"), m.MentionName, m.Name, member.TimeZone)
				c.sendChat(m, text)
				return true
			}
		}
	}
	return false
}
func (c *Hs) TzGetTime(m models.IncomingMessage) bool {
	if strings.HasPrefix(m.Text, "%time") || strings.HasPrefix(m.Text, "%время") || strings.HasPrefix(m.Text, "%час") {
		members, err := c.corpMember.CorpMembersRead(m.GuildId)
		if err != nil {
			c.log.ErrorErr(err)
			return false
		}

		// Исходные данные
		data := [][]string{
			{"Local Time", "User", ""},
		}

		for _, member := range members {
			if member.TimeZone != "" {
				newRow := []string{member.LocalTime, member.Name, ""}
				data = append(data, newRow)
			}
		}
		text := fmt.Sprintf(c.getText(m, "LOCAL_TIME_FOR_EVERYONE"), m.MentionName)
		c.sendFormatedText(m, text, data)
		c.sendChat(m, c.getText(m, "UNLISTED_MEMBERS"))
		return true
	}
	return false
}

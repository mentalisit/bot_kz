package logic

import (
	"compendium/models"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
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
	if m.MAcc != nil {
		member, err := c.db.V2.CorpMemberByUId(m.MAcc.UUID)
		if err == nil && member != nil {
			member.TimeZona = timeZona
			member.ZonaOffset = offsetInt
			err = c.db.V2.CorpMemberUpdate(*member)
			if err != nil {
				c.log.ErrorErr(err)
			}
		}
	}

	cm := models.CorpMember{
		Name:         m.Name,
		GuildId:      m.MGuild.GuildId(),
		Tech:         models.TechLevels{},
		AvatarUrl:    m.Avatar,
		TimeZone:     timeZona,
		ZoneOffset:   offsetInt,
		MGuild:       m.MGuild,
		MultiAccount: m.MultiAccount,
	}
	u := models.User{
		Username:  m.Name,
		AvatarURL: m.Avatar,
		Alts:      []string{},
	}

	if mentionName == "" {
		var err error
		if m.MultiAccount != nil {
			err = c.db.Multi.CorpMemberTZUpdate(m.MultiAccount.UUID, timeZona, offsetInt)
		} else {
			err = c.corpMember.CorpMemberTZUpdate(m.NameId, m.MGuild.GuildId(), timeZona, offsetInt)
		}

		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				cm.UserId = m.NameId
				if m.MultiAccount == nil {
					_ = c.corpMember.CorpMemberInsert(cm)
					u.ID = m.NameId
					_ = c.users.UsersInsert(u)
				}
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
			if m.MAcc != nil {
				corpMember, err := c.db.V2.CorpMemberByUId(m.MAcc.UUID)
				if err == nil && corpMember != nil && matches[1] == m.MAcc.DiscordID {
					corpMember.TimeZona = timeZona
					corpMember.ZonaOffset = offsetInt
					err = c.db.V2.CorpMemberUpdate(*corpMember)
					if err != nil {
						c.log.ErrorErr(err)
					}
				}
			}
			err := c.corpMember.CorpMemberTZUpdate(matches[1], m.MGuild.GuildId(), timeZona, offsetInt)
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
			if m.MAcc != nil {
				corpMember, err := c.db.V2.CorpMemberByUId(m.MAcc.UUID)
				if err == nil && corpMember != nil && matches[2] == m.MAcc.TelegramUsername {
					corpMember.TimeZona = timeZona
					corpMember.ZonaOffset = offsetInt
					err = c.db.V2.CorpMemberUpdate(*corpMember)
					if err != nil {
						c.log.ErrorErr(err)
					}
				}
			}
			user, err := c.users.UsersGetByUserName(matches[2])
			if err != nil {
				c.log.ErrorErr(err)
				return
			}
			err = c.corpMember.CorpMemberTZUpdate(user.ID, m.MGuild.GuildId(), timeZona, offsetInt)
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
		if m.MAcc != nil {
			corpMember, err := c.db.V2.CorpMemberByUId(m.MAcc.UUID)
			if err == nil && corpMember != nil {
				text := fmt.Sprintf(c.getText(m, "TIMEZONA_IS_CURRENTLY"), m.MentionName, m.Name, corpMember.TimeZona)
				c.sendChat(m, text)
				return true
			}
		}
		members, err := c.corpMember.CorpMembersRead(m.MGuild.GuildId())
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
		// Исходные данные
		data := [][]string{
			{"Local Time", "User", ""},
		}

		members, err := c.corpMember.CorpMembersRead(m.MGuild.GuildId())
		if err != nil {
			c.log.ErrorErr(err)
		}
		if m.MGuild != nil {
			membersRead, _ := c.db.Multi.CorpMembersRead(m.MGuild.GId)
			if len(membersRead) != 0 {
				for _, member := range membersRead {
					accountUUID, _ := c.db.Multi.FindMultiAccountUUID(member.Uid)
					if accountUUID == nil {
						continue
					}
					mmember := models.CorpMember{
						Name:       accountUUID.Nickname,
						UserId:     member.Uid.String(),
						AvatarUrl:  accountUUID.AvatarURL,
						TimeZone:   member.TimeZona,
						ZoneOffset: member.ZonaOffset,
						AfkFor:     member.AfkFor,
					}
					if mmember.TimeZone != "" {
						t12, t24 := getTimeStrings(mmember.ZoneOffset)
						mmember.LocalTime = t12
						mmember.LocalTime24 = t24
					}

					members = append(members, mmember)
				}
			}
			memberS, _ := c.db.V2.CorpMembersReadMulti(&m.MGuild.GId)
			if len(memberS) != 0 {
				for _, member := range memberS {
					mmember := models.CorpMember{
						Name:       member.Multi.Nickname,
						UserId:     member.Multi.UUID.String(),
						AvatarUrl:  member.AvatarUrl,
						TimeZone:   member.TimeZone,
						ZoneOffset: member.ZoneOffset,
						AfkFor:     member.AfkFor,
					}
					if mmember.TimeZone != "" {
						t12, t24 := getTimeStrings(mmember.ZoneOffset)
						mmember.LocalTime = t12
						mmember.LocalTime24 = t24
					}

					members = append(members, mmember)
				}
			}
		}

		for _, member := range members {
			if member.TimeZone != "" {
				newRow := []string{member.LocalTime24, member.Name, ""}
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

func getTimeStrings(offset int) (string, string) {
	// Получаем текущее время в UTC
	now := time.Now().UTC()

	// Применяем смещение к текущему времени в UTC
	offsetDuration := time.Duration(offset) * time.Minute
	timeWithOffset := now.Add(offsetDuration)

	// Форматируем время в 12-часовом формате с AM/PM
	time12HourFormat := timeWithOffset.Format("03:04 PM")

	// Форматируем время в 24-часовом формате
	time24HourFormat := timeWithOffset.Format("15:04")

	return time12HourFormat, time24HourFormat
}

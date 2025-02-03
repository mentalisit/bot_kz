package logic

import (
	"compendium/models"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func (c *Hs) wskill(m models.IncomingMessage) bool {
	after, found := strings.CutPrefix(m.Text, "%wskill")
	if !found {
		after, found = strings.CutPrefix(m.Text, "%бзоткат")
	}
	if !found {
		return false
	}
	re := regexp.MustCompile(`(\S+|@\S+|<@\d+>)(\s+)?(\S+)?`)
	matches := re.FindStringSubmatch(after)
	if len(matches) == 0 || len(after) < 3 {
		c.wskillList(m)
		return true
	}
	name := matches[1]
	ship := matches[3]
	afterMatches := re.FindStringIndex(after)
	endIndex := afterMatches[1]
	textAfterMatch := after[endIndex:]
	c.wskillNameShip(m, name, ship, textAfterMatch)
	return true
}
func (c *Hs) wskillList(m models.IncomingMessage) {
	ws, err := c.db.DB.WsKillReadByGuildId(m.GuildId)
	if err != nil {
		c.log.ErrorErr(err)
		c.log.Info(m.GuildId)
		return
	}
	if len(ws) > 0 {
		// Исходные данные
		data := [][]string{
			{"Returns in", "Name", "Ship"},
		}

		for _, w := range ws {
			newRow := []string{c.timeUntil(m, w.TimestampEnd), w.UserName, w.ShipName}
			data = append(data, newRow)
		}
		c.sendFormatedText(m, fmt.Sprintf(c.getText(m, "SCHEDULED_RETURNS"), m.MentionName), data)
	} else {
		c.sendChat(m, fmt.Sprintf(c.getText(m, "NO_SHIP_ARE_SCHEDULED"), m.MentionName))
	}
}

// todo translate
func (c *Hs) wskillNameShip(m models.IncomingMessage, name, ship, afterText string) {
	if name == "" || ship == "" {
		text := fmt.Sprintf("%s, The wskill command accepts any of the following formats: `wskill name ship`, `wskill name ship <time>`, or `wskill name ship delete`.", m.MentionName)
		c.sendChat(m, text)
		return
	}
	if len(afterText) <= 2 {
		now := time.Now().UTC()
		add := now.Add(18 * time.Hour)
		timestamp := add.Unix()
		wskill := models.WsKill{
			GuildId:      m.GuildId,
			ChatId:       m.ChannelId,
			UserName:     c.getNameText(name),
			Mention:      name,
			ShipName:     ship,
			TimestampEnd: timestamp,
			Language:     m.Language,
		}
		count, _ := c.guilds.GuildGetCountByGuildId(m.GuildId)
		if count == 0 {
			err := c.guilds.GuildInsert(models.Guild{
				URL:  m.GuildAvatar,
				ID:   m.GuildId,
				Name: m.GuildName,
				Type: m.Type,
			})
			if err != nil {
				c.log.ErrorErr(err)
			}
		}

		err := c.db.DB.WsKillInsert(wskill)
		if err != nil {
			c.log.ErrorErr(err)
			c.log.InfoStruct("WsKillInsert", wskill)
			return
		} else {
			addLocal := add.In(c.getTimeLocation(m.NameId))
			text := fmt.Sprintf(c.getText(m, "IS_DUE_TO_RETURN"),
				m.MentionName, name, ship, addLocal.Format("15:04:05 -07:00"), c.timeUntil(m, add.Unix()))
			c.sendChat(m, text)
		}
	} else {
		re := regexp.MustCompile(` (delete|back in)`)
		matches := re.FindStringSubmatch(afterText)
		if len(matches) > 0 {
			if matches[1] == "delete" {
				wskill := models.WsKill{
					GuildId:  m.GuildId,
					UserName: c.getNameText(name),
					ShipName: ship,
				}
				err := c.db.DB.WsKillDelete(wskill)
				text := fmt.Sprintf("%s, %s %s Kill ", m.MentionName, m.Name, ship)
				if err != nil {
					c.sendChat(m, text+"Not Found")
					c.log.ErrorErr(err)
					c.log.InfoStruct("WsKillDelete", wskill)
				} else {
					c.sendChat(m, text+"Removed")
				}

				return
			} else if matches[1] == "back in" {
				afterMatches := re.FindStringIndex(afterText)
				endIndex := afterMatches[1]
				textAfterMatch := afterText[endIndex:]
				re = regexp.MustCompile(` (\d{1,2})h ?(\d{1,2})?`)
				matches = re.FindStringSubmatch(textAfterMatch)

				hour := 0
				if len(matches) > 1 && matches[1] != "" {
					hour, _ = strconv.Atoi(matches[1])
				}

				minute := 0
				if len(matches) > 1 && matches[2] != "" {
					minute, _ = strconv.Atoi(matches[2])
				}

				now := time.Now().UTC()

				timeHour := now.Add(time.Duration(hour) * time.Hour)
				timeMinute := timeHour.Add(time.Duration(minute) * time.Minute)
				timestamp := timeMinute.Unix()

				wskill := models.WsKill{
					GuildId:      m.GuildId,
					ChatId:       m.ChannelId,
					UserName:     c.getNameText(name),
					Mention:      name,
					ShipName:     ship,
					TimestampEnd: timestamp,
					Language:     m.Language,
				}
				err := c.db.DB.WsKillInsert(wskill)
				if err != nil {
					c.log.ErrorErr(err)
					c.log.InfoStruct("WsKillInsert", wskill)
					return
				} else {
					addLocal := timeMinute.In(c.getTimeLocation(m.NameId))
					text := fmt.Sprintf(c.getText(m, "IS_DUE_TO_RETURN"),
						m.MentionName, name, ship, addLocal.Format("15:04:05 -07:00"), c.timeUntil(m, timestamp))
					c.sendChat(m, text)
				}
			}
		}
		re = regexp.MustCompile(` (\d{1,2}(?:am|pm)) ?(\d{1,2})?m?`)
		matches = re.FindStringSubmatch(afterText)
		if len(matches) > 0 {
			Hour := matches[1]
			Minute := matches[2]
			c.wskillDeadIn(m, Hour, Minute, name, ship)
		}
	}
}
func (c *Hs) wskillDeadIn(m models.IncomingMessage, hour, minute, name, ship string) {
	// Парсинг времени
	layout := "3pm"
	t, err := time.Parse(layout, hour)
	if err != nil {
		c.log.ErrorErr(err)
		c.log.Info(hour)
		return
	}
	minutes := 0
	// Извлечение количества минут
	if minute != "" {
		minutes, _ = strconv.Atoi(minute)
	}
	// Текущее время для создания времени с сегодняшней датой
	location := c.getTimeLocation(m.NameId)
	now := time.Now().In(location)
	parsedTime := time.Date(now.Year(), now.Month(), now.Day(), t.Hour(), minutes, 0, 0, location)

	wskill := models.WsKill{
		GuildId:      m.GuildId,
		ChatId:       m.ChannelId,
		UserName:     c.getNameText(name),
		Mention:      name,
		ShipName:     ship,
		TimestampEnd: parsedTime.Unix(),
		Language:     m.Language,
	}
	err = c.db.DB.WsKillInsert(wskill)
	if err != nil {
		c.log.ErrorErr(err)
		c.log.InfoStruct("WsKillInsert", wskill)
	} else {
		text := fmt.Sprintf(c.getText(m, "IS_DUE_TO_RETURN"),
			m.MentionName, name, ship, parsedTime.Format("15:04:05 -07:00"), c.timeUntil(m, parsedTime.Unix()))
		c.sendChat(m, text)
	}
}

func (c *Hs) getTimeLocation(userid string) *time.Location {
	mem, errt := c.corpMember.CorpMemberByUserId(userid)
	if errt != nil {
		c.log.ErrorErr(errt)
		return time.UTC
	}
	if mem.TimeZone == "" && mem.ZoneOffset == 0 {
		return time.UTC
	}
	return time.FixedZone(mem.TimeZone, mem.ZoneOffset*60)
}

// Возвращает остаток времени
func (c *Hs) timeUntil(m models.IncomingMessage, wstimestamp int64) string {
	timestamp := time.Unix(wstimestamp, 0)
	// Получаем текущее время
	now := time.Now().UTC()

	// Вычисляем разницу между указанным таймштампом и текущим временем
	duration := timestamp.Sub(now)

	// Проверяем, если разница меньше нуля, значит указанное время уже прошло
	if duration < 0 {
		return c.getText(m, "TIME_HAS_ALREADY_PASSED")
	}

	// Извлекаем часы, минуты и секунды из разницы
	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) % 60
	seconds := int(duration.Seconds()) % 60

	// Форматируем строку в формате "17h 10m 45s"
	return fmt.Sprintf(c.getText(m, "H_M_S"), hours, minutes, seconds)
}
func (c *Hs) wsKillTimer() {
	for {
		now := time.Now().UTC()
		if now.Second() == 0 {
			all, err := c.db.DB.WsKillReadAll()
			if err != nil {
				c.log.ErrorErr(err)
				return
			}
			if all != nil && len(all) > 0 {
				for _, kill := range all {
					h, m := getTimeLeft(kill.TimestampEnd)
					if h == 0 && (m == 15 || m == 0) {
						guild, errs := c.guilds.GuildGet(kill.GuildId)
						if errs != nil {
							c.log.ErrorErr(errs)
							c.log.InfoStruct("kill. ", kill)
							continue
						}
						in := models.IncomingMessage{
							ChannelId: kill.ChatId,
							Type:      guild.Type,
							Language:  kill.Language,
						}
						if m == 15 {
							text := fmt.Sprintf(c.getText(in, "WILL_BE_ABLE_TO_RETURN"),
								kill.Mention, kill.ShipName)
							c.sendChat(in, text)
						} else if m == 0 {
							err = c.db.DB.WsKillDelete(kill)
							if err != nil {
								c.log.ErrorErr(err)
								c.log.InfoStruct("WsKillDelete", kill)
								return
							}
							text := fmt.Sprintf(c.getText(in, "IS_NOW_ABLE_TO_RETURN"),
								kill.Mention, kill.ShipName)
							c.sendChat(in, text)
						}
					}
				}
			}
			time.Sleep(20 * time.Second)
		} else if now.Hour() == 23 && now.Minute() == 45 && now.Second() == 30 {
			go c.updateAvatars()
			time.Sleep(time.Second)
		} else {
			time.Sleep(time.Second)
		}
	}
}
func getTimeLeft(timestamp int64) (int, int) {
	now := time.Now().UTC()
	ws := time.Unix(timestamp, 0)

	// Вычисляем разницу между указанным таймштампом и текущим временем
	duration := ws.Sub(now)

	// Проверяем, если разница меньше нуля, значит указанное время уже прошло
	if duration < 0 {
		return 0, 0
	}

	// Извлекаем часы, минуты из разницы
	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) % 60
	return hours, minutes
}
func (c *Hs) getNameText(name string) string {
	re := regexp.MustCompile(`<@(\d+)>`)
	matches := re.FindStringSubmatch(name)
	if len(matches) > 0 {
		user, err := c.users.UsersGetByUserId(matches[1])
		if err != nil {
			return name
		} else {
			return user.Username
		}
	}
	if name[:1] == "@" {
		return name[1:]
	}
	return name
}

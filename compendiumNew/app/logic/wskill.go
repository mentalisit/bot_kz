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
	re := regexp.MustCompile(`(\S+|@\S+|<@\d+>) ?(\w+)?`)
	matches := re.FindStringSubmatch(after)
	if len(matches) == 0 || len(after) < 3 {
		c.wskillList(m)
		return true
	}
	name := matches[1]
	ship := matches[2]
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
			{"Returns in", "NameNameName", "Ship"},
		}

		for _, w := range ws {
			newRow := []string{timeUntil(w.TimestampEnd), w.UserName, w.ShipName}
			data = append(data, newRow)
		}
		// Определяем максимальную длину для каждого столбца
		colWidths := make([]int, len(data[0]))
		for _, row := range data {
			for i, col := range row {
				if len(col) > colWidths[i] {
					colWidths[i] = len(col)
				}
			}
		}

		// Формируем формат для печати строк
		format := ""
		for _, width := range colWidths {
			format += fmt.Sprintf("%%-%ds  ", width)
		}
		format = format[:len(format)-2] // Убираем последний лишний пробел

		// Печатаем строки с выравниванием
		text := ""
		for _, row := range data {
			text += fmt.Sprintf(format+"\n", row[0], row[1], row[2])
		}
		c.sendChat(m, fmt.Sprintf("%s, Scheduled WS Returns\n%s", m.MentionName, text))
	} else {
		c.sendChat(m, m.MentionName+", No ships are scheduled to return.")
	}
}
func (c *Hs) wskillNameShip(m models.IncomingMessage, name, ship, afterText string) {
	if len(afterText) <= 2 {
		now := time.Now().UTC()
		add := now.Add(18 * time.Hour)
		timestamp := add.Unix()
		wskill := models.WsKill{
			GuildId:      m.GuildId,
			ChatId:       m.ChannelId,
			UserName:     m.Name,
			Mention:      name,
			ShipName:     ship,
			TimestampEnd: timestamp,
		}
		err := c.db.DB.WsKillInsert(wskill)
		if err != nil {
			c.log.ErrorErr(err)
			c.log.InfoStruct("WsKillInsert", wskill)
			return
		} else {
			// Установка временной зоны (в данном случае -05:00)
			//location := time.FixedZone("EST", -5*3600)
			add = add.In(time.UTC)
			text := fmt.Sprintf("%s %s's %s is due to return at %s (%s)",
				m.MentionName, name, ship, add.Format("15:04:05 -07:00"), timeUntil(add.Unix()))
			c.sendChat(m, text)
		}
	} else {
		re := regexp.MustCompile(` (delete|back in)`)
		matches := re.FindStringSubmatch(afterText)
		if len(matches) > 0 {
			if matches[1] == "delete" {
				wskill := models.WsKill{
					GuildId:  m.GuildId,
					UserName: m.Name,
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
				if matches[1] != "" {
					hour, _ = strconv.Atoi(matches[1])
				}

				minute := 0
				if matches[2] != "" {
					minute, _ = strconv.Atoi(matches[2])
				}

				now := time.Now().UTC()

				timeHour := now.Add(time.Duration(hour) * time.Hour)
				timeMinute := timeHour.Add(time.Duration(minute) * time.Minute)
				timestamp := timeMinute.Unix()

				wskill := models.WsKill{
					GuildId:      m.GuildId,
					ChatId:       m.ChannelId,
					UserName:     m.Name,
					Mention:      name,
					ShipName:     ship,
					TimestampEnd: timestamp,
				}
				err := c.db.DB.WsKillInsert(wskill)
				if err != nil {
					c.log.ErrorErr(err)
					c.log.InfoStruct("WsKillInsert", wskill)
					return
				} else {
					add := timeMinute.In(time.UTC)
					text := fmt.Sprintf("%s %s's %s is due to return at %s (%s)",
						m.MentionName, name, ship, add.Format("15:04:05 -07:00"), timeUntil(timestamp))
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
	now := time.Now().UTC()
	parsedTime := time.Date(now.Year(), now.Month(), now.Day(), t.Hour(), minutes, 0, 0, now.UTC().Location())

	wskill := models.WsKill{
		GuildId:      m.GuildId,
		ChatId:       m.ChannelId,
		UserName:     m.Name,
		Mention:      name,
		ShipName:     ship,
		TimestampEnd: parsedTime.Unix(),
	}
	err = c.db.DB.WsKillInsert(wskill)
	if err != nil {
		c.log.ErrorErr(err)
		c.log.InfoStruct("WsKillInsert", wskill)
	} else {
		add := parsedTime.In(time.UTC)
		text := fmt.Sprintf("%s %s's %s is due to return at %s (%s)",
			m.MentionName, name, ship, add.Format("15:04:05 -07:00"), timeUntil(parsedTime.Unix()))
		c.sendChat(m, text)
	}
}
func timeUntil(wstimestamp int64) string {
	timestamp := time.Unix(wstimestamp, 0)
	// Получаем текущее время
	now := time.Now().UTC()

	// Вычисляем разницу между указанным таймштампом и текущим временем
	duration := timestamp.Sub(now)

	// Проверяем, если разница меньше нуля, значит указанное время уже прошло
	if duration < 0 {
		return "time has already passed"
	}

	// Извлекаем часы, минуты и секунды из разницы
	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) % 60
	seconds := int(duration.Seconds()) % 60

	// Форматируем строку в формате "17h 10m 45s"
	return fmt.Sprintf("%dh %dm %ds", hours, minutes, seconds)
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
							continue
						}
						in := models.IncomingMessage{
							ChannelId: kill.ChatId,
							Type:      guild.Type,
							Language:  "",
						}
						if m == 15 {
							text := fmt.Sprintf("%s's %s will be able to return to the White Star in 15 minutes.",
								kill.Mention, kill.ShipName)

							c.sendChat(in, text)
						} else if m == 0 {
							err = c.db.DB.WsKillDelete(kill)
							if err != nil {
								c.log.ErrorErr(err)
								c.log.InfoStruct("WsKillDelete", kill)
								return
							}
							text := fmt.Sprintf("%s's %s is now able to return to the White Star",
								kill.Mention, kill.ShipName)
							c.sendChat(in, text)
						}
					}
				}
			}
			time.Sleep(50 * time.Second)
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

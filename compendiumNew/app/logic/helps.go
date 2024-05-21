package logic

import (
	"compendium/models"
	"strings"
)

func (c *Hs) Help(m models.IncomingMessage) bool {
	after, found := strings.CutPrefix(m.Text, "%")
	if found {
		split := strings.Split(after, " ")
		lenSplit := len(split)
		if split[0] == "help" || split[0] == "помощь" || split[0] == "допомога" {
			if lenSplit == 1 {
				c.helpGeneral(m)
				return true
			} else if lenSplit == 2 {
				if split[1] == "tech" || split[1] == "техи" {
					//send help tech
				} else if split[1] == "alts" || split[1] == "альт" {
					//send help alts
				} else if split[1] == "wskill" || split[1] == "бзоткат" {
					//send help wskill
				} else if split[1] == "tz" {
					//send help TZ
				} else if split[1] == "nick" || split[1] == "ник" {
					//send help Nick
				} else {
					c.helpGeneral(m)
				}
			}

		}

	}
	//регулярное выражение на хелп wskill time tz

	return false
}
func (c *Hs) helpGeneral(m models.IncomingMessage) {
	if m.Type == "ds" {
		c.sendChat(m, c.getText(m, "HELP_TEXT_DS"))
	} else {
		c.sendChat(m, c.getText(m, "HELP_TEXT_TG"))
	}

}

//на текуший момент доступны команды:\n
//'**%помощь**' для получения общей справки о командах \n
//'**%подключить**' для подключения приложения\n\n
//'**%помощь техи**' справка по тех. модулям\n
//'**%помощь альт**' справка по настройке альтов\n
//'**%помощь бзоткат**' справка по откату после смерти\n
//'**%помощь альт**' справка по настройке времени\n
//'**%помощь ник**' справка по замене имени\n
//
//на текущий момент доступно только вывод изображения с модулями\n
//'**%т и**' для получения изображения с вашими модулями\n
//'**%т @имя и**' для получения изображения с модулями другого игрока\n
//'**%т имя и**' для получения изображения с модулями альта\n
//
//
//
//
//
//'**%alts add NameAlt**' для создания альта для технологий\n
//'**%alts del NameAlt**' для удаления альта\n
//

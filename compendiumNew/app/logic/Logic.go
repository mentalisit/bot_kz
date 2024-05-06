package logic

import (
	"compendium/models"
	"fmt"
	"strings"
)

func (c *Hs) logic(m models.IncomingMessage) {
	fmt.Printf("logic %+v\n", m)
	cutPrefix, found := strings.CutPrefix(m.Text, "%")

	if found {
		c.in = m
		switch cutPrefix {
		case "help", "помощь", "справка":
			c.help()
		case "подключить", "connect":
			c.connect()
		case "т и", "t i":
			c.techImage()
		default:

			c.regular(cutPrefix)
		}
	}
}
func (c *Hs) regular(text string) {
	if c.techImageName() {
	} else if c.logicRoles() {
	} else if c.createAlt() {
	} else {
		c.log.Info("else " + text)
	}
}
func (c *Hs) help() {
	if c.in.Type == "ds" {
		c.sendChat("на текуший момент доступны команды:\n" +
			"'%help' или '%помощь', '%справка' для получения текущей справки \n" +
			"'%connect' или '%подключить' для подключения приложения\n" +
			"'%t i' или '%т и' для получения изображения с вашими модулями\n" +
			"'%t @Name i' или '%т @имя и' для получения изображения с модулями другого игрока\n" +
			"'%alts add NameAlt' для создания альта для технологий\n" +
			"'%alts del NameAlt' для удаления альта")
	} else {
		c.sendChat("на текуший момент доступны команды:\n" +
			"'%help' или '%помощь', '%справка' для получения текущей справки \n" +
			"'%connect' или '%подключить' для подключения приложения\n" +
			"'%t i' или '%т и' для получения изображения с вашими модулями\n" +
			"'%t @Name i' или '%т @имя и' для получения изображения с модулями другого игрока\n" +
			"'%role create RoleName' создание роли для телеграм\n" +
			"'%role delete Rolename' удаление роли для телеграм\n" +
			"'%role s RoleName' для подписки на роль\n" +
			"'%role u RoleName' для удаления подписки на роль\n" +
			"'%alts add NameAlt' для создания альта для технологий\n" +
			"'%alts del NameAlt' для удаления альта")
	}

}

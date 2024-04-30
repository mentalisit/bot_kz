package Compendium

import (
	"compendium/models"
	"fmt"
	"strings"
)

func (c *Compendium) logic(m models.IncomingMessage) {
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
func (c *Compendium) regular(text string) {
	if c.techImageName() {
	} else if c.logicRoles() {
	} else {
		c.log.Info(text)
	}
}
func (c *Compendium) help() {
	c.sendChat("на текуший момент доступны команды:\n" +
		"'%help' или '%помощь', '%справка' для получения текущей справки \n" +
		"'%connect' или '%коннект' для подключения приложения\n" +
		"'%t i' или '%т и' для получения изображения с вашими модулями\n" +
		"'%t @Name i' или '%т @имя и' для получения изображения с модулями другого игрока\n")
}

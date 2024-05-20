package logic

import (
	"compendium/models"
	"fmt"
	"strings"
)

func (c *Hs) logic(m models.IncomingMessage) {
	fmt.Printf("logic: %s %+v\n", m.Language, m)
	cutPrefix, found := strings.CutPrefix(m.Text, "%")

	if found {
		switch cutPrefix {
		case "help", "помощь", "допомога":
			c.help(m)
		case "подключить", "connect", "підключити":
			c.connect(m)
		case "т и", "t i", "т і":
			c.techImage(m)
		default:

			c.regular(m)
		}
	}
}
func (c *Hs) regular(m models.IncomingMessage) {
	if c.techImageName(m) {
	} else if c.techImageNameAlt(m) {
	} else if c.logicRoles(m) {
	} else if c.createAlt(m) {
	} else if c.wskill(m) {
	} else {
		c.log.Info("else " + m.Text)
	}
}
func (c *Hs) help(m models.IncomingMessage) {
	if m.Type == "ds" {
		c.sendChat(m, c.getText(m, "HELP_TEXT_DS"))
	} else {
		c.sendChat(m, c.getText(m, "HELP_TEXT_TG"))
	}

}

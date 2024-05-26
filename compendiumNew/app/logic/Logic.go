package logic

import (
	"compendium/models"
	"fmt"
	"strings"
	"time"
)

func (c *Hs) logic(m models.IncomingMessage) {
	fmt.Printf("logic: %+v %+v\n", time.Now().Format(time.RFC3339), m)
	cutPrefix, found := strings.CutPrefix(m.Text, "%")

	if found {
		switch cutPrefix {
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
	if c.Help(m) {
	} else if c.techImageName(m) {
	} else if c.techImageNameAlt(m) {
	} else if c.logicRoles(m) {
	} else if c.createAlt(m) {
	} else if c.wskill(m) {
	} else if c.TzTime(m) {
	} else if c.setGameName(m) {
	} else {
		c.log.Info("else " + m.Text)
	}
}

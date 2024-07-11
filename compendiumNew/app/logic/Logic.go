package logic

import (
	"compendium/models"
	"fmt"
	"time"
)

func (c *Hs) logic(m models.IncomingMessage) {
	fmt.Printf("logic: %+v %+v\n", time.Now().Format(time.RFC3339), m)
	if c.connect(m) {
	} else if c.Help(m) {
	} else if c.techImage(m) {
	} else if c.techImageName(m) {
	} else if c.techImageNameAlt(m) {
	} else if c.logicRoles(m) {
	} else if c.createAlt(m) {
	} else if c.wskill(m) {
	} else if c.TzTime(m) {
	} else if c.setGameName(m) {
	} else {
		c.log.Info(fmt.Sprintf("else %+v\n", m))
	}
}

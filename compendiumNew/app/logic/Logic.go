package logic

import (
	"compendium/models"
	"fmt"
	"strings"
	"time"
)

func (c *Hs) logic(m models.IncomingMessage) {
	startTime := time.Now()
	defer func() {
		fmt.Printf("Команда %s Отправитель %s Время выполнения: %v\n", m.Text, m.Name, time.Since(startTime))
	}()

	c.PrintGoroutine()

	if strings.Contains(m.Type, "DM") && !strings.HasPrefix(m.Text, "%") {
		return
	}
	if m.MGuild == nil {
		c.log.InfoStruct("m.MGuild==nil ", m)
	}

	fmt.Printf("logic: %+v %+v\n", time.Now().Format(time.RFC3339), m)

	if c.connect(m) {
	} else if c.Help(m) {
	} else if c.techImage(m) {
	} else if c.techImageName(m) {
	} else if c.techImageNameAlt(m) {
	} else if c.createAlt(m) {
	} else if c.wsKill(m) {
	} else if c.TzTime(m) {
	} else if c.setGameName(m) {
	} else if c.removeMember(m) {
	} else {
		c.sendChat(m, c.getText(m, "ErrorRequest"))
		fmt.Printf("else Corp:%s %+v\n", m.MGuild.GuildName, m)
	}
}

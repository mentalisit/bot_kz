package NewLogic

import (
	"compendium/models"
	"fmt"
	"strings"
	"time"
)

func (c *HsLogic) Logic(m models.IncomingMessage) {
	if strings.Contains(m.Type, "DM") && !strings.HasPrefix(m.Text, "%") {
		return
	}

	fmt.Printf("logic: %+v %+v\n", time.Now().Format(time.RFC3339), m)

	if c.connect(m) {
		//} else if c.multiConnect(m) {
		//} else if c.Help(m) {
		//} else if c.techImage(m) {
		//} else if c.techImageName(m) {
		//} else if c.techImageNameAlt(m) {
		//} else if c.logicRoles(m) {
		//} else if c.createAlt(m) {
		//} else if c.wskill(m) {
		//} else if c.TzTime(m) {
		//} else if c.setGameName(m) {
		//} else if c.removeMember(m) {
	} else {
		c.sendChat(m, c.getText(m, "ErrorRequest"))
		fmt.Printf("else Corp:%s %+v\n", m.MultiGuild.GuildName, m)
	}
}

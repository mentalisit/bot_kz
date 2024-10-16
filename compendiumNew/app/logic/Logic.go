package logic

import (
	"compendium/models"
	"fmt"
	"runtime"
	"time"
)

func (c *Hs) logic(m models.IncomingMessage) {
	c.PrintGoroutine()
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
	} else if c.removeMember(m) {
	} else {
		c.log.Info(fmt.Sprintf("else %+v\n", m))
	}
}
func (c *Hs) PrintGoroutine() {
	goroutine := runtime.NumGoroutine()
	tm := time.Now()
	mdate := (tm.Format("2006-01-02"))
	mtime := (tm.Format("15:04"))
	text := fmt.Sprintf(" %s %s Горутин  %d\n", mdate, mtime, goroutine)
	if goroutine > 120 {
		c.log.Info(text)
		c.log.Panic(text)
	} else if goroutine > 50 && goroutine%10 == 0 {
		c.log.Info(text)
	}

	fmt.Println(text)
}

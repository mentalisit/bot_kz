package main

import (
	"sort"
	"time"
	"ws/hspublic"
	"ws/server"

	"github.com/mentalisit/conf"
	"github.com/mentalisit/conf/logger"
)

var log *logger.Logger

func main() {
	_, logg, db := conf.InitConf("WS")
	log = logg
	hs := hspublic.NewHS(log, db)

	go server.NewSrv(log, hs.Db)

	newContent := hs.GetContentAll()
	hs.DownloadFile("wsAll", newContent)
	sort.Slice(newContent, func(i, j int) bool {
		return newContent[i].LastModified < newContent[j].LastModified
	})
	hs.SavePercent(newContent)

	var count int

	for {
		now := time.Now()

		if now.Second() == 0 && (now.Minute() == 0 || now.Minute()%5 == 0) {
			newContent = hs.GetContentSevenDays()
			sort.Slice(newContent, func(i, j int) bool {
				return newContent[i].LastModified < newContent[j].LastModified
			})
			hs.SavePercent(newContent)
			//hs.DownloadFile("ws", newContent)

			newContent = hs.GetContentAll()

			if count == len(newContent) {
				continue
			}

			hs.DownloadFile("wsAll", newContent)
			count = len(newContent)
		}

		time.Sleep(1 * time.Second)
	}
}

package main

import (
	"github.com/mentalisit/logger"
	"sort"
	"time"
	"ws/config"
	"ws/hspublic"
	"ws/server"
)

var log *logger.Logger

func main() {
	cfg := config.InitConfig()
	log = logger.LoggerZap(cfg.Logger.Token, cfg.Logger.ChatId, cfg.Logger.Webhook, "WS")
	hs := hspublic.NewHS(log, cfg.Postgress.Password)

	go server.NewSrv(log)

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

package main

import (
	"github.com/mentalisit/logger"
	"time"
	"ws/config"
	"ws/hspublic"
)

var log *logger.Logger

func main() {
	cfg := config.InitConfig()
	log = logger.LoggerZap(cfg.Logger.Token, cfg.Logger.ChatId, cfg.Logger.Webhook)
	hs := hspublic.NewHS(log)

	newContent := hs.GetContentSevenDays()
	hs.DownloadFile("ws", newContent)

	newContent = hs.GetContentAll()
	hs.DownloadFile("wsAll", newContent)

	for {
		now := time.Now()

		if now.Second() == 0 && now.Minute() == 0 {
			newContent = hs.GetContentSevenDays()
			hs.DownloadFile("ws", newContent)

			newContent = hs.GetContentAll()
			hs.DownloadFile("wsAll", newContent)
		}

		time.Sleep(1 * time.Second)
	}
}

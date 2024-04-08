package main

import (
	"github.com/mentalisit/logger"
	"os"
	"os/signal"
	"syscall"
	"telegram/telegram"
)

func main() {
	cfg := InitConfig()

	log := logger.LoggerZap(cfg.Logger.Token, cfg.Logger.ChatId, cfg.Logger.Webhook)

	_ = telegram.NewTelegram(log, cfg.TokenTelegram)

	log.Info("Service telegram load")
	//ожидаем сигнала завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

}

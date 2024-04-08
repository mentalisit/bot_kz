package main

import (
	"bridge/config"
	"bridge/server"
	"github.com/mentalisit/logger"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.InitConfig()

	log := logger.LoggerZap(cfg.Logger.Token, cfg.Logger.ChatId, cfg.Logger.Webhook)

	server.NewBridge(log)

	log.Info("Service bridge load")
	//ожидаем сигнала завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
}

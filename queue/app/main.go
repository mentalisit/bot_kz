package main

import (
	"github.com/mentalisit/logger"
	"os"
	"os/signal"
	"queue/config"
	"queue/server"
	"syscall"
)

func main() {
	cfg := config.InitConfig()
	log := logger.LoggerZap(cfg.Logger.Token, cfg.Logger.ChatId, cfg.Logger.Webhook, "queue")

	server.NewServer(log, cfg)

	log.Info("Service queue load")

	//ожидаем сигнала завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
}

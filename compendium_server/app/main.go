package main

import (
	"compendium_s/config"
	"compendium_s/server"
	"compendium_s/storage"
	"github.com/mentalisit/logger"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.InitConfig()
	log := logger.LoggerZap(cfg.Logger.Token, cfg.Logger.ChatId, cfg.Logger.Webhook)
	st := storage.NewStorage(log, cfg)

	server.NewServer(log, cfg, st)

	log.Info("Service compendium server load")
	//ожидаем сигнала завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
}

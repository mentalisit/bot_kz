package main

import (
	"compendium/config"
	"compendium/logic"
	"compendium/server"
	"compendium/storage"
	"github.com/mentalisit/logger"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.InitConfig()
	log := logger.LoggerZap(cfg.Logger.Token, cfg.Logger.ChatId, cfg.Logger.Webhook, "CompendiumLogic")
	st := storage.NewStorage(log, cfg)

	s := server.NewServer(log, st)
	logic.NewCompendium(log, s.In, st)

	log.Info("Service compendiumNew load")
	//ожидаем сигнала завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
}

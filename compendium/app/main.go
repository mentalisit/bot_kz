package main

import (
	"compendium/Compendium"
	"compendium/config"
	"compendium/server"
	"compendium/storage"
	"github.com/mentalisit/logger"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.InitConfig()
	log := logger.LoggerZapDiscord(cfg.Logger.Webhook)
	st := storage.NewStorage(log, cfg)

	//d := ds.NewDiscord(log, st, cfg)
	//t := tg.NewTelegram(log, cfg, st)

	s := server.NewServer(log, cfg, st)
	Compendium.NewCompendium(log, s.In, st)

	log.Info("Service compendium load")
	//ожидаем сигнала завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
}

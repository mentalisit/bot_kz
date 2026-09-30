package main

import (
	"compendium_s/server"
	"compendium_s/storage"
	"os"
	"os/signal"
	"syscall"

	"github.com/mentalisit/conf"
)

func main() {
	_, log, db := conf.InitConf("CompendiumServer")
	st := storage.NewStorage(log, db)

	server.NewServer(log, st)

	log.Info("Service compendium server load")
	//ожидаем сигнала завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
}

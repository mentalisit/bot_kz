package main

import (
	"compendium/logic"
	"compendium/server"
	"compendium/storage"
	"os"
	"os/signal"
	"syscall"

	"github.com/mentalisit/conf"
)

func main() {
	_, log, db := conf.InitConf("COMPENDIUM")
	st := storage.NewStorage(log, db)

	s := server.NewServer(log, st)
	logic.NewCompendium(log, s.In, st)

	log.Info("Service compendiumNew load")
	//ожидаем сигнала завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
}

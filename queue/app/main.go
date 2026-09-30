package main

import (
	"os"
	"os/signal"
	"queue/server"
	"syscall"

	"github.com/mentalisit/conf"
)

func main() {
	_, log, db := conf.InitConf("QUEUE")

	server.NewServer(log, db)

	log.Info("Service queue load")

	//ожидаем сигнала завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
}

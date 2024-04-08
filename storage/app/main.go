package main

import (
	"github.com/mentalisit/logger"
	"os"
	"os/signal"
	"storage/config"
	"storage/mongodb"
	"storage/server"
	"syscall"
)

func main() {
	cfg := config.InitConfig()
	log := logger.LoggerZapDiscord(cfg.Logger.Webhook)
	db, err := mongodb.InitMongoDB(log, cfg.Mongo)
	if err != nil {
		log.Panic(err.Error())
	}

	server.NewServer(db, log)

	log.Info("Service storage load")

	//ожидаем сигнала завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
}

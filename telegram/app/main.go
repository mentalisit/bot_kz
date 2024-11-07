package main

import (
	"github.com/mentalisit/logger"
	"os"
	"os/signal"
	"syscall"
	"telegram/config"
	"telegram/grpc_server"
	"telegram/server"
	"telegram/storage"
	"telegram/telegram"
)

func main() {
	cfg := config.InitConfig()

	log := logger.LoggerZap(cfg.Logger.Token, cfg.Logger.ChatId, cfg.Logger.Webhook, "TG")

	st := storage.NewStorage(log, cfg)

	tg := telegram.NewTelegram(log, cfg.Token.TokenTelegram, st)

	grpc_server.GrpcMain(tg, log)

	server.NewServer(tg, log)

	log.Info("Service telegram load")

	//ожидаем сигнала завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
	tg.Close()

}

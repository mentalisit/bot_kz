package main

import (
	"bridge/config"
	grpc_server "bridge/grpc-server"
	"bridge/server"
	"bridge/storage"
	"github.com/mentalisit/logger"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.InitConfig()

	log := logger.LoggerZap(cfg.Logger.Token, cfg.Logger.ChatId, cfg.Logger.Webhook, "bridge")

	st := storage.NewStorage(log, cfg)
	b := server.NewBridge(log, st)
	grpc_server.GrpcMain(b, log)

	log.Info("Service bridge load")
	//ожидаем сигнала завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
}

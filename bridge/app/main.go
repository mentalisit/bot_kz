package main

import (
	"bridge/config"
	grpc_server "bridge/grpc-server"
	"bridge/logic"
	"bridge/storage"
	"os"
	"os/signal"
	"syscall"

	"github.com/mentalisit/logger"
)

func main() {
	cfg := config.InitConfig()

	log := logger.LoggerZap(cfg.Logger.Token, cfg.Logger.ChatId, cfg.Logger.Webhook, "bridge")

	st := storage.NewStorage(log, cfg)
	b := logic.NewBridge(log, st, cfg)
	grpc_server.GrpcMain(b, log)

	log.Info("Service bridge load")
	//ожидаем сигнала завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	// Корректное завершение всех модулей
	log.Info("Received shutdown signal, shutting down...")
	b.Shutdown()
	log.Info("Service bridge stopped")
}

//go:generate protoc --go_out=. --go-grpc_out=. bridge.proto

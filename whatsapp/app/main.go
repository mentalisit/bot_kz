package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"whatsapp/config"
	"whatsapp/grpc_server"
	"whatsapp/storage"
	wa "whatsapp/whatsapp"

	"github.com/mentalisit/logger"
)

func main() {
	cfg := config.InitConfig()

	log := logger.LoggerZap(cfg.Logger.Token, cfg.Logger.ChatId, cfg.Logger.Webhook, "WA")

	st := storage.NewStorage(log, cfg)
	fmt.Printf("Loading Number:%s SesionFile:%s\n", cfg.Whatsapp.Number, cfg.Whatsapp.SessionFile)
	w := wa.NewWhatsapp(log, cfg, st)

	grpc_server.GrpcMain(w, log)

	//log.Info("Service whatsapp load")

	//ожидаем сигнала завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	//whatsapp.Disconnect()
}

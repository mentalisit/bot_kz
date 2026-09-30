package main

import (
	DiscordClient "discord/discord"
	"discord/storage"
	"os"
	"os/signal"
	"syscall"

	"discord/grpc_server"

	"github.com/mentalisit/conf"
)

func main() {
	cfg, log, db := conf.InitConf("DISCORD")
	st := storage.NewStorage(log, db)

	ds := DiscordClient.NewDiscord(log, st, cfg)
	grpc_server.GrpcMain(ds, log)

	//server.NewServer(ds, log)

	log.Info("Service discord load")

	//ожидаем сигнала завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	ds.Shutdown()

}

//go:generate protoc --go_out=. --go-grpc_out=. discord.proto

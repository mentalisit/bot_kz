package main

import (
	"crypto/tls"
	"io"
	"log"
	"net"
	"time"
)

func main() {
	// Определение конфигураций для нескольких экземпляров прокси
	configs := []ProxyConfig{
		{PrimaryService: "ws:4443", BackupService: "ws2:4443", ListenPort: ":14443"},
		{PrimaryService: "compendium_s:8443", BackupService: "compendium_s2:8443", ListenPort: ":18443"},
		{PrimaryService: "queue:9443", BackupService: "queue2:9443", ListenPort: ":19443"},
	}

	//configsGRPC := []ProxyConfig{
	//	{PrimaryService: "telegram:50051", BackupService: "telegram2:50051", ListenPort: "5051"}, //telegram:5051
	//	{PrimaryService: "discord:50051", BackupService: "discord2:50051", ListenPort: "5052"},   //discord:5052
	//	{PrimaryService: "rs_bot:50051", BackupService: "rs_bot2:50051", ListenPort: "5053"},     //rs_bot:5053
	//	{PrimaryService: "bridge:50051", BackupService: "bridge2:50051", ListenPort: "5054"},     //bridge:5054
	//	//{PrimaryService: "compendiumnew:50051",BackupService: "compendiumnew2:50051",ListenPort: "5055"},//compendiumnew:5055
	//}

	// Запуск каждого экземпляра прокси в отдельной горутине
	for _, config := range configs {
		time.Sleep(1 * time.Second)
		log.Printf("Loading ProxyTLS %s ListenPort %s\n", config.PrimaryService, config.ListenPort)
		go startProxyTls(config)
	}

	// Блокировка основного потока для ожидания завершения горутин
	select {}
}

func startProxyTls(config ProxyConfig) {
	primaryAvailable := checkServiceTls(config.PrimaryService)

	go func() {
		for {
			primaryAvailable = checkServiceTls(config.PrimaryService)
			time.Sleep(checkInterval)
		}
	}()

	// Запуск прокси-сервера
	listener, err := net.Listen("tcp", config.ListenPort)
	if err != nil {
		log.Fatalf("Failed to start proxy on port %s: %v", config.ListenPort, err)
	}
	defer listener.Close()
	log.Printf("TCP Proxy started on %s, Primary: %s, Backup: %s", config.ListenPort, config.PrimaryService, config.BackupService)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}

		go handleConnection(conn, primaryAvailable, config)
	}
}
func checkServiceTls(service string) bool {
	// Создаем настроенный Dialer для установки таймаута
	dialer := &net.Dialer{Timeout: timeoutDuration}
	// Используем TLS-соединение для проверки, если сервис требует его
	conn, err := tls.DialWithDialer(dialer, "tcp", service, &tls.Config{
		InsecureSkipVerify: true, // Пропуск проверки сертификата для тестирования
	})
	if err != nil {
		log.Printf("Service %s is not available: %v", service, err)
		return false
	}
	_ = conn.Close()
	//log.Printf("Service %s is available", service)
	return true
}
func handleConnection(clientConn net.Conn, primaryAvailable bool, config ProxyConfig) {
	defer clientConn.Close()

	// Определение целевого сервиса
	targetService := config.BackupService
	if primaryAvailable {
		targetService = config.PrimaryService
	}

	// Установка соединения с целевым сервисом
	serverConn, err := net.Dial("tcp", targetService)
	if err != nil {
		log.Printf("Failed to connect to %s: %v", targetService, err)
		return
	}
	defer serverConn.Close()

	// Дуплексное копирование данных между клиентом и сервером
	go io.Copy(serverConn, clientConn)
	io.Copy(clientConn, serverConn)
}

const (
	checkInterval   = 5 * time.Second
	timeoutDuration = 3 * time.Second
)

// Структура конфигурации прокси
type ProxyConfig struct {
	PrimaryService string
	BackupService  string
	ListenPort     string
}

//// Проверка доступности gRPC-сервиса
//func checkServiceGrpc(service string) bool {
//	conn, err := grpc.Dial(service, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock(), grpc.WithTimeout(3*time.Second))
//	if err != nil {
//		log.Printf("Service %s is not available: %v", service, err)
//		return false
//	}
//	_ = conn.Close()
//	log.Printf("Service %s is available", service)
//	return true
//}
//func startGrpcProxy(config ProxyConfig) {
//	primaryAvailable := checkServiceGrpc(config.PrimaryService)
//
//	go func() {
//		for {
//			primaryAvailable = checkServiceGrpc(config.PrimaryService)
//			time.Sleep(checkInterval)
//		}
//	}()
//
//	listener, err := net.Listen("tcp", config.ListenPort)
//	if err != nil {
//		log.Fatalf("Failed to start gRPC proxy on port %s: %v", config.ListenPort, err)
//	}
//	defer listener.Close()
//	log.Printf("gRPC Proxy started on %s, Primary: %s, Backup: %s", config.ListenPort, config.PrimaryService, config.BackupService)
//
//	for {
//		conn, err := listener.Accept()
//		if err != nil {
//			log.Printf("Failed to accept connection: %v", err)
//			continue
//		}
//
//		go func(conn net.Conn) {
//			defer conn.Close()
//			targetService := config.PrimaryService
//			if !primaryAvailable {
//				targetService = config.BackupService
//			}
//
//			backendConn, err := grpc.Dial(targetService, grpc.WithTransportCredentials(insecure.NewCredentials()))
//			if err != nil {
//				log.Printf("Failed to connect to backend service %s: %v", targetService, err)
//				return
//			}
//			defer backendConn.Close()
//
//			// Прокси-трафик между conn и backendConn
//			// Например, используйте проксирование gRPC-запросов
//		}(conn)
//	}
//}

package TgApi

import (
	"fmt"
	"github.com/mentalisit/logger"
	"google.golang.org/grpc"
	"sync"
	"time"
)

type Client struct {
	conn      *grpc.ClientConn
	client    TelegramServiceClient
	log       *logger.Logger
	primary   string
	backup    string
	mutex     sync.Mutex // Защита для повторного подключения
	isPrimary bool       // Флаг для отслеживания текущего подключения
}

func NewClient(log *logger.Logger) *Client {
	client := &Client{
		log:       log,
		primary:   "telegram:50051",
		backup:    "telegram2:50051",
		isPrimary: true,
	}

	client.connect(client.primary)

	// Запуск фоновой проверки для переподключения к основному сервису
	go client.monitorPrimary()

	return client

	//conn, err := grpc.Dial("telegram:50051", grpc.WithInsecure())
	//if err != nil {
	//	log.ErrorErr(err)
	//	return nil
	//}
	//fmt.Println("connect to grpc ok")
	//return &Client{
	//	conn:   conn,
	//	client: NewTelegramServiceClient(conn),
	//	log:    log,
	//}
}

func (c *Client) connect(address string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// Close the existing connection, if it exists
	if c.conn != nil {
		c.conn.Close()
	}

	// Try to connect to the new address
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(3*time.Second))
	if err != nil {
		c.log.Error(fmt.Sprintf("Failed to connect to %s", address))
		// If the primary server is not available, try to connect to the backup server
		if address == c.primary {
			c.connect(c.backup)
		} else {
			c.connect(c.primary)
		}
		return
	}

	// Successful connection
	c.conn = conn
	c.client = NewTelegramServiceClient(conn)
	fmt.Printf("Connected to %s\n", address)
	c.isPrimary = (address == c.primary)
}

func (c *Client) monitorPrimary() {
	for {
		time.Sleep(10 * time.Second) // Check interval

		if !c.isPrimary {
			conn, err := grpc.Dial(c.primary, grpc.WithInsecure(), grpc.WithTimeout(3*time.Second))
			if err == nil {
				conn.Close()
				fmt.Println("Primary service is back online, switching to primary")
				c.connect(c.primary)
			}
		}
	}
}

func (c *Client) MonitorPrimary() {
	if !c.isPrimary {
		conn, err := grpc.Dial(c.primary, grpc.WithInsecure(), grpc.WithTimeout(3*time.Second))
		if err == nil {
			conn.Close()
			fmt.Println("Primary service is back online, switching to primary")
			c.connect(c.primary)
		}
	}
}

// Close закрывает соединение
func (c *Client) Close() error {
	return c.conn.Close()
}

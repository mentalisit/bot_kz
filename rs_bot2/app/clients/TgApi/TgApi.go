package TgApi

import (
	"fmt"
	"github.com/mentalisit/logger"
	"google.golang.org/grpc"
)

type Client struct {
	conn   *grpc.ClientConn
	client TelegramServiceClient
	log    *logger.Logger
}

func NewClient(log *logger.Logger) *Client {
	target := "telegram:50051"
	conn, err := grpc.Dial(target, grpc.WithInsecure())
	if err != nil {
		log.ErrorErr(err)
		return nil
	}
	fmt.Printf("connect to grpc %s ok\n", target)
	return &Client{
		conn:   conn,
		client: NewTelegramServiceClient(conn),
		log:    log,
	}
}

// Close закрывает соединение
func (c *Client) Close() error {
	return c.conn.Close()
}

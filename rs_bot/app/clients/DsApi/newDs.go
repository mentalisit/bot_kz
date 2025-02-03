package ds

import (
	"fmt"
	"github.com/mentalisit/logger"
	"google.golang.org/grpc"
)

type Client struct {
	conn     *grpc.ClientConn
	client   BotServiceClient
	log      *logger.Logger
	rolePing map[string]map[string]string
}

func NewClient(log *logger.Logger) *Client {
	target := "discord:50051"
	conn, err := grpc.Dial(target, grpc.WithInsecure())
	if err != nil {
		log.ErrorErr(err)
		return nil
	}
	fmt.Printf("connect to grpc %s ok\n", target)
	return &Client{
		conn:     conn,
		client:   NewBotServiceClient(conn),
		log:      log,
		rolePing: make(map[string]map[string]string),
	}
}

// Close закрывает соединение
func (c *Client) Close() error {
	return c.conn.Close()
}

package compendium

import (
	"fmt"
	"github.com/mentalisit/logger"
	"google.golang.org/grpc"
)

type Client struct {
	conn   *grpc.ClientConn
	client LogicServiceClient
	log    *logger.Logger
}

func NewClient(log *logger.Logger) *Client {
	conn, err := grpc.Dial("compendiumnew:50051", grpc.WithInsecure())
	if err != nil {
		log.ErrorErr(err)
		return nil
	}
	fmt.Println("connect to rs grpc ok")
	return &Client{
		conn:   conn,
		client: NewLogicServiceClient(conn),
		log:    log,
	}
}

// Close закрывает соединение
func (c *Client) Close() error {
	return c.conn.Close()
}

package dbredis

import (
	"context"
	"fmt"
	"github.com/mentalisit/logger"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()
var log *logger.Logger

type Db struct {
	c *redis.Client
}

func NewDb(myLogger *logger.Logger) *Db {
	log = myLogger
	client := redis.NewClient(&redis.Options{
		Addr:     "192.168.100.131:6379", // Replace with your Redis server address
		Password: "",                     // No password for local development
		DB:       0,                      // Default DB
	})

	// Ping the Redis server to check the connection
	pong, err := client.Ping(ctx).Result()
	if err != nil {
		log.ErrorErr(err)
	}
	fmt.Println("Connected to Redis:", pong)
	return &Db{c: client}
}

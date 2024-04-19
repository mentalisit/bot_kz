package mongo

import (
	"context"
	"github.com/mentalisit/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"kz_bot/pkg/clientDB/mongodb"
)

type DB struct {
	s   *mongo.Client
	log *logger.Logger
}

func InitMongoDB(log *logger.Logger) *DB {
	client, err := mongodb.NewMongoClient()
	if err != nil {
		log.ErrorErr(err)
		return nil
	}

	d := &DB{
		s:   client,
		log: log,
	}
	return d
}
func (d *DB) Shutdown() {
	err := d.s.Disconnect(context.Background())
	if err != nil {
		d.log.ErrorErr(err)
		return
	}
}

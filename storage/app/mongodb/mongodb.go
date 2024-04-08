package mongodb

import (
	"context"
	"github.com/mentalisit/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DB struct {
	s   *mongo.Client
	log *logger.Logger
}

func InitMongoDB(log *logger.Logger, uri string) (*DB, error) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		log.ErrorErr(err)
		return nil, err
	}

	d := &DB{
		s:   client,
		log: log,
	}
	return d, nil
}

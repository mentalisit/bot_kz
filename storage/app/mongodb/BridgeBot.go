package mongodb

import (
	"context"
	"storage/models"

	"go.mongodb.org/mongo-driver/bson"
)

func (d *DB) DBReadBridgeConfig() (data []models.BridgeConfig, err error) {
	//collection := d.s.Database("BridgeChat").Collection("Bridge")
	collection := d.s.Database("BridgeChat").Collection("Bridge")
	find, err := collection.Find(context.Background(), bson.D{})
	if err != nil {
		d.log.ErrorErr(err)
		return nil, err
	}
	err = find.All(context.Background(), &data)
	if err != nil {
		d.log.ErrorErr(err)
		return nil, err
	}
	return data, nil
}
func (d *DB) UpdateBridgeChat(br models.BridgeConfig) {
	collection := d.s.Database("BridgeChat").Collection("Bridge")
	filter := bson.M{"namerelay": br.NameRelay}
	collection.FindOneAndDelete(context.Background(), filter)
	d.InsertBridgeChat(br)
}
func (d *DB) InsertBridgeChat(br models.BridgeConfig) {
	collection := d.s.Database("BridgeChat").Collection("Bridge")
	bsonData, _ := bson.Marshal(br)
	_, err := collection.InsertOne(context.Background(), bsonData)
	if err != nil {
		d.log.ErrorErr(err)
	}
}

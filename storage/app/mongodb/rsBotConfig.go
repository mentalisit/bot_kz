package mongodb

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"storage/models"
)

func (d *DB) collectionRsConfig() *mongo.Collection {
	collection := d.s.Database("RsBot").Collection("RsConfig")
	return collection
}
func (d *DB) ReadConfigRs() []models.CorporationConfig {
	cursor, err := d.collectionRsConfig().Find(context.Background(), bson.M{})
	if err != nil {
		d.log.ErrorErr(err)
		return nil
	}
	var m []models.CorporationConfig
	err = cursor.All(context.Background(), &m)
	if err != nil {
		d.log.ErrorErr(err)
		return nil
	}
	return m
}
func (d *DB) InsertConfigRs(c models.CorporationConfig) {
	ins, err := d.collectionRsConfig().InsertOne(context.Background(), c)
	if err != nil {
		d.log.ErrorErr(err)
	}
	fmt.Println(ins.InsertedID)
}
func (d *DB) DeleteConfigRs(c models.CorporationConfig) {
	ins, err := d.collectionRsConfig().DeleteOne(context.Background(), c)
	if err != nil {
		d.log.ErrorErr(err)
	}
	fmt.Println(ins.DeletedCount)
}
func (d *DB) UpdateRsConfig(c models.CorporationConfig) {
	filter := bson.M{"corpname": c.CorpName}
	_, err := d.collectionRsConfig().ReplaceOne(context.Background(), filter, c)
	if err != nil {
		d.log.ErrorErr(err)
	}
}

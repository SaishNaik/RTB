package datastore

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"ssp/lib/mongo"
)

type IMongoStore interface {
	Find(ctx context.Context, mongoClient *mongo.Client, collectionName string,
		criteria bson.M, opt bson.M) ([]bson.M, error)
	UpdateOne(ctx context.Context, mongoClient *mongo.Client, collectionName string,
		filter interface{}, update interface{}) (int64, error)
	InsertOne(ctx context.Context, mongoClient *mongo.Client, collectionName string,
		doc interface{}) (interface{}, error)
}

type mongoStore struct {
}

func InitMongoStore() IMongoStore {
	return &mongoStore{}
}

// TODO: add more functions

func (m *mongoStore) Find(ctx context.Context, mongoClient *mongo.Client, collectionName string,
	filter bson.M, opt bson.M) ([]bson.M, error) {

	results, err := mongoClient.Find(nil, collectionName, filter, opt)
	return results, err
}

func (m *mongoStore) UpdateOne(ctx context.Context, mongoClient *mongo.Client, collectionName string,
	filter interface{}, update interface{}) (int64, error) {

	modifiedCount, err := mongoClient.UpdateOne(nil, collectionName, filter, update)
	return modifiedCount, err
}

func (m *mongoStore) InsertOne(ctx context.Context, mongoClient *mongo.Client, collectionName string,
	doc interface{}) (interface{}, error) {

	createdId, err := mongoClient.InsertOne(nil, collectionName, doc)
	return createdId, err
}

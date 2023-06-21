package model

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	datastore2 "ssp/datastore"
	"ssp/lib/mongo"
	"ssp/model/datastore"
)

type IAdSlotModel interface {
	GetAdSlotDetails(ctx context.Context, mongoClient *mongo.Client, adSlotId string) (slotDetails *datastore.AdSlot, error error)
}

type AdSlotModel struct {
	store      datastore2.IMongoStore
	collection string
}

//todo correct initialisation
func InitAdSlotModel(store datastore2.IMongoStore, collection string) IAdSlotModel {
	return &AdSlotModel{store: store, collection: collection}
}

func (as *AdSlotModel) GetAdSlotDetails(ctx context.Context, mongoClient *mongo.Client, adSlotId string) (*datastore.AdSlot, error) {

	adslotIdHex, err := primitive.ObjectIDFromHex(adSlotId)
	if err != nil {
		return nil, err
	}

	filter := map[string]interface{}{
		"_id": adslotIdHex,
	}
	option := map[string]interface{}{}

	slots, err := as.store.Find(ctx, mongoClient, as.collection, filter, option)
	if err != nil {
		return nil, err
	}
	result := &datastore.AdSlot{}
	for _, slot := range slots {
		marshalResult, err := bson.Marshal(slot)
		if err != nil {
			return nil, err
		}
		err = bson.Unmarshal(marshalResult, &result)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

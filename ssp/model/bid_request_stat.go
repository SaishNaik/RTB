package model

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	datastore2 "ssp/datastore"
	"ssp/lib/mongo"
	"ssp/model/datastore"
)

//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=./mock_bid_request_stat.go -package=model . IBidRequestStatModel

type IBidRequestStatModel interface {
	SaveBidRequestStat(ctx context.Context, mongoClient *mongo.Client, stat *datastore.BidRequestStat) error
	GetBidRequestStat(ctx context.Context, mongoClient *mongo.Client, bidReqId string) (stat *datastore.BidRequestStat, error error)
	UpdateImpression(ctx context.Context, mongoClient *mongo.Client, bidReqId string, stat *datastore.BidRequestStat) error
}

type BidRequestStatModel struct {
	store      datastore2.IMongoStore
	collection string
}

//todo correct initialisation
func InitBidRequestStatModel(store datastore2.IMongoStore, collection string) IBidRequestStatModel {
	return &BidRequestStatModel{store: store, collection: collection}
}

func (b *BidRequestStatModel) SaveBidRequestStat(ctx context.Context, mongoClient *mongo.Client,
	stat *datastore.BidRequestStat) error {
	//todo return from db
	// get relations
	_, err := b.store.InsertOne(ctx, mongoClient, b.collection, stat)
	return err
}

func (b *BidRequestStatModel) GetBidRequestStat(ctx context.Context, mongoClient *mongo.Client, bidReqId string) (*datastore.BidRequestStat, error) {
	filter := map[string]interface{}{
		"bid_req_id": bidReqId,
	}
	option := map[string]interface{}{}

	stats, err := b.store.Find(ctx, mongoClient, b.collection, filter, option)
	if err != nil {
		return nil, err
	}
	bidRequestStat := &datastore.BidRequestStat{}
	for _, stat := range stats {
		marshalResult, err := bson.Marshal(stat)
		if err != nil {
			return nil, err
		}
		err = bson.Unmarshal(marshalResult, &bidRequestStat)
		if err != nil {
			return nil, err
		}
	}
	return bidRequestStat, nil
}

func (b *BidRequestStatModel) UpdateImpression(ctx context.Context, mongoClient *mongo.Client, bidReqId string, stat *datastore.BidRequestStat) error {

	update := map[string]interface{}{
		"$set": map[string]interface{}{
			"impression": 1,
			"revenue":    stat.ExpectedRevenue,
			"profit":     stat.ExpectedProfit,
		},
	}
	filter := map[string]interface{}{
		"bid_req_id": bidReqId,
	}
	_, err := b.store.UpdateOne(ctx, mongoClient, b.collection, filter, update)
	return err
}

package model

import (
	"context"
	_ "github.com/golang/mock/mockgen/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	datastore2 "ssp/datastore"
	"ssp/lib/mongo"
	"ssp/model/datastore"
)

//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=./mock_dsp_pub.go -package=model . IDSPPubModel

// todo correct modelling from db
// also current setting will support only one kind of db, need to make proper interface
type IDSPPubModel interface {
	GetAllDSPForPub(ctx context.Context, mongoClient *mongo.Client,
		pubId string) (dsps []datastore.DSP, err error)
}

type DSPPubModel struct {
	store      datastore2.IMongoStore
	collection string
}

//todo correct initialisation
func InitDSPPubModel(store datastore2.IMongoStore, collection string) IDSPPubModel {
	return &DSPPubModel{store: store, collection: collection}
}

func (d *DSPPubModel) GetAllDSPForPub(ctx context.Context, mongoClient *mongo.Client,
	pubId string) (dsps []datastore.DSP, err error) {

	pubObjectId, err := primitive.ObjectIDFromHex(pubId)
	if err != nil {
		return nil, err
	}

	filter := map[string]interface{}{
		"pub_id": pubObjectId,
	}

	option := map[string]interface{}{}

	// get relations
	relation, err := d.store.Find(ctx, mongoClient, d.collection, filter, option)
	if err != nil {
		return nil, err
	}

	dspPub := &datastore.DSPPub{}
	for _, rel := range relation {
		marshalResult, err := bson.Marshal(rel)
		if err != nil {
			return nil, err
		}
		err = bson.Unmarshal(marshalResult, &dspPub)
		if err != nil {
			return nil, err
		}
	}

	return dspPub.DSPList, nil
}

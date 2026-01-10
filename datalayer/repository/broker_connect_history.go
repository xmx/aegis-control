package repository

import (
	"context"

	"github.com/xmx/aegis-control/datalayer/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type BrokerConnectHistory interface {
	Repository[bson.ObjectID, model.BrokerConnectHistory, model.BrokerConnectHistories]
}

func NewBrokerConnectHistory(db *mongo.Database, opts ...options.Lister[options.CollectionOptions]) BrokerConnectHistory {
	coll := db.Collection("broker_connect_history", opts...)
	repo := NewRepository[bson.ObjectID, model.BrokerConnectHistory, model.BrokerConnectHistories](coll)

	return &brokerConnectHistoryRepo{
		Repository: repo,
	}
}

type brokerConnectHistoryRepo struct {
	Repository[bson.ObjectID, model.BrokerConnectHistory, model.BrokerConnectHistories]
}

func (r *brokerConnectHistoryRepo) CreateIndex(ctx context.Context) error {
	idx := []mongo.IndexModel{
		{Keys: bson.D{{Key: "broker_id", Value: 1}}},
		{Keys: bson.D{{Key: "tunnel_stat.connected_at", Value: -1}}},
		{Keys: bson.D{{Key: "tunnel_stat.disconnected_at", Value: -1}}},
	}
	_, err := r.Indexes().CreateMany(ctx, idx)

	return err
}

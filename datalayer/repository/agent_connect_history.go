package repository

import (
	"context"

	"github.com/xmx/aegis-control/datalayer/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type AgentConnectHistory interface {
	Repository[bson.ObjectID, model.AgentConnectHistory, model.AgentConnectHistories]
}

func NewAgentConnectHistory(db *mongo.Database, opts ...options.Lister[options.CollectionOptions]) AgentConnectHistory {
	coll := db.Collection("agent_connect_history", opts...)
	repo := NewRepository[bson.ObjectID, model.AgentConnectHistory, model.AgentConnectHistories](coll)

	return &agentConnectHistoryRepo{
		Repository: repo,
	}
}

type agentConnectHistoryRepo struct {
	Repository[bson.ObjectID, model.AgentConnectHistory, model.AgentConnectHistories]
}

func (r *agentConnectHistoryRepo) CreateIndex(ctx context.Context) error {
	idx := []mongo.IndexModel{
		{Keys: bson.D{{Key: "agent_id", Value: 1}}},
		{Keys: bson.D{{Key: "tunnel_stat.connected_at", Value: -1}}},
		{Keys: bson.D{{Key: "tunnel_stat.disconnected_at", Value: -1}}},
	}
	_, err := r.Indexes().CreateMany(ctx, idx)

	return err
}

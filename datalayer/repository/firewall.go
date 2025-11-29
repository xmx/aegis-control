package repository

import (
	"context"

	"github.com/xmx/aegis-control/datalayer/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Firewall interface {
	Repository[bson.ObjectID, model.Firewall, []*model.Firewall]
}

func NewFirewall(db *mongo.Database, opts ...options.Lister[options.CollectionOptions]) Firewall {
	coll := db.Collection("firewall", opts...)
	repo := NewRepository[bson.ObjectID, model.Firewall, []*model.Firewall](coll)

	return &firewallRepo{
		Repository: repo,
	}
}

type firewallRepo struct {
	Repository[bson.ObjectID, model.Firewall, []*model.Firewall]
}

func (r *firewallRepo) CreateIndex(ctx context.Context) error {
	idx := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "name", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	}
	_, err := r.Indexes().CreateMany(ctx, idx)

	return err
}

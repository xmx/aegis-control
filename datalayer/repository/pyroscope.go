package repository

import (
	"context"

	"github.com/xmx/aegis-control/datalayer/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Pyroscope interface {
	Repository[bson.ObjectID, model.Pyroscope, []*model.Pyroscope]
}

func NewPyroscope(db *mongo.Database, opts ...options.Lister[options.CollectionOptions]) Pyroscope {
	coll := db.Collection("pyroscope", opts...)
	repo := NewRepository[bson.ObjectID, model.Pyroscope, []*model.Pyroscope](coll)

	return &pyroscopeRepo{
		Repository: repo,
	}
}

type pyroscopeRepo struct {
	Repository[bson.ObjectID, model.Pyroscope, []*model.Pyroscope]
}

func (r *pyroscopeRepo) CreateIndex(ctx context.Context) error {
	idx := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "name", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	}
	_, err := r.Indexes().CreateMany(ctx, idx)

	return err
}

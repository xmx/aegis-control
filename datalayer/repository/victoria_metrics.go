package repository

import (
	"context"

	"github.com/xmx/aegis-control/datalayer/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type VictoriaMetrics interface {
	Repository[bson.ObjectID, model.VictoriaMetrics, []*model.VictoriaMetrics]
}

func NewVictoriaMetrics(db *mongo.Database, opts ...options.Lister[options.CollectionOptions]) VictoriaMetrics {
	coll := db.Collection("victoria_metrics", opts...)
	repo := NewRepository[bson.ObjectID, model.VictoriaMetrics, []*model.VictoriaMetrics](coll)

	return &victoriaMetricsRepo{
		Repository: repo,
	}
}

type victoriaMetricsRepo struct {
	Repository[bson.ObjectID, model.VictoriaMetrics, []*model.VictoriaMetrics]
}

func (r *victoriaMetricsRepo) CreateIndex(ctx context.Context) error {
	idx := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "name", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	}
	_, err := r.Indexes().CreateMany(ctx, idx)

	return err
}

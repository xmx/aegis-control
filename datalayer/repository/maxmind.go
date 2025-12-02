package repository

import (
	"github.com/xmx/aegis-control/datalayer/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Maxmind interface {
	Repository[bson.ObjectID, model.Maxmind, []*model.Maxmind]
}

func NewMaxmind(db *mongo.Database, opts ...options.Lister[options.CollectionOptions]) Maxmind {
	coll := db.Collection("maxmind", opts...)
	repo := NewRepository[bson.ObjectID, model.Maxmind, []*model.Maxmind](coll)

	return &maxmindRepo{
		Repository: repo,
	}
}

type maxmindRepo struct {
	Repository[bson.ObjectID, model.Maxmind, []*model.Maxmind]
}

package repository

import (
	"context"
	"time"

	"github.com/xmx/aegis-control/datalayer/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Setting interface {
	Get(ctx context.Context) (*model.Setting, error)
	Upsert(ctx context.Context, data *model.SettingData) error
}

func NewSetting(db *mongo.Database, opts ...options.Lister[options.CollectionOptions]) Setting {
	coll := db.Collection("setting", opts...)
	repo := NewRepository[bson.ObjectID, model.Setting, []*model.Setting](coll)

	return &settingRepo{
		repo: repo,
	}
}

type settingRepo struct {
	repo Repository[bson.ObjectID, model.Setting, []*model.Setting]
}

func (r *settingRepo) Get(ctx context.Context) (*model.Setting, error) {
	opt := options.FindOne().SetSort(bson.D{{Key: "_id", Value: 1}})
	return r.repo.FindOne(ctx, bson.D{}, opt)
}

func (r *settingRepo) Upsert(ctx context.Context, data *model.SettingData) error {
	var sd model.SettingData
	if data != nil {
		sd = *data
	}

	now := time.Now()
	mod := &model.Setting{SettingData: sd, UpdatedAt: now}
	opt := options.UpdateOne().SetUpsert(true)
	update := bson.M{"$set": mod, "$setOnInsert": bson.M{"created_at": now}}
	_, err := r.repo.UpdateOne(ctx, update, opt)

	return err
}

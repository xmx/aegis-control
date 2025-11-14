package repository

import (
	"context"
	"crypto/tls"

	"github.com/xmx/aegis-control/datalayer/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Certificate interface {
	Repository[bson.ObjectID, model.Certificate, model.Certificates]
	Enables(context.Context) ([]*tls.Certificate, error)
}

func NewCertificate(db *mongo.Database, opts ...options.Lister[options.CollectionOptions]) Certificate {
	coll := db.Collection("certificate", opts...)
	repo := NewRepository[bson.ObjectID, model.Certificate, model.Certificates](coll)

	return &certificateRepo{
		Repository: repo,
	}
}

type certificateRepo struct {
	Repository[bson.ObjectID, model.Certificate, model.Certificates]
}

// Enables 主要是给 tlscert 使用。
func (r *certificateRepo) Enables(ctx context.Context) ([]*tls.Certificate, error) {
	filter := bson.D{{"enabled", true}}
	dats, err := r.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	rets := make([]*tls.Certificate, 0, len(dats))
	for _, dat := range dats {
		pair, err := tls.X509KeyPair([]byte(dat.PublicKey), []byte(dat.PrivateKey))
		if err != nil {
			return nil, err
		}
		rets = append(rets, &pair)
	}

	return rets, nil
}

func (r *certificateRepo) CreateIndex(ctx context.Context) error {
	idx := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "certificate_sha256", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{Key: "name", Value: 1}},
		},
	}
	_, err := r.Indexes().CreateMany(ctx, idx)

	return err
}

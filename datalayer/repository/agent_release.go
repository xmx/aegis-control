package repository

import (
	"context"
	"io"

	"github.com/xmx/aegis-control/datalayer/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type AgentSaveInfo struct {
	FileID   bson.ObjectID  `json:"file_id"  bson:"file_id"`
	Length   int64          `json:"length"   bson:"length"`
	Checksum model.Checksum `json:"checksum" bson:"checksum"`
}

type AgentRelease interface {
	Repository[bson.ObjectID, model.AgentRelease, model.AgentReleases]
	SaveFile(ctx context.Context, rd io.Reader, filename string) (*AgentSaveInfo, error)
	OpenFile(ctx context.Context, fileID bson.ObjectID) (*mongo.GridFSDownloadStream, error)
	DeleteFile(ctx context.Context, fileID bson.ObjectID) error
}

func NewAgentRelease(db *mongo.Database, opts ...options.Lister[options.CollectionOptions]) AgentRelease {
	coll := db.Collection("agent_release", opts...)
	repo := NewRepository[bson.ObjectID, model.AgentRelease, model.AgentReleases](coll)

	return &agentReleaseRepo{
		Repository: repo,
	}
}

type agentReleaseRepo struct {
	Repository[bson.ObjectID, model.AgentRelease, model.AgentReleases]
}

func (r *agentReleaseRepo) Bucket(opts ...options.Lister[options.BucketOptions]) *mongo.GridFSBucket {
	name := r.Name()
	opt := options.GridFSBucket().SetName(name)
	opts = append(opts, opt)

	return r.Database().GridFSBucket(opts...)
}

func (r *agentReleaseRepo) SaveFile(ctx context.Context, rd io.Reader, filename string) (*AgentSaveInfo, error) {
	fileID := bson.NewObjectID()
	bucket := r.Bucket()
	stm, err := bucket.OpenUploadStreamWithID(ctx, fileID, filename)
	if err != nil {
		return nil, err
	}
	defer stm.Close()

	hw := model.NewHashWriter()
	mw := io.MultiWriter(hw, stm)
	length, err := io.Copy(mw, rd)
	if err != nil {
		return nil, err
	}
	chk := hw.Sum()
	ret := &AgentSaveInfo{
		FileID:   fileID,
		Length:   length,
		Checksum: chk,
	}

	return ret, nil
}

func (r *agentReleaseRepo) OpenFile(ctx context.Context, fileID bson.ObjectID) (*mongo.GridFSDownloadStream, error) {
	bucket := r.Bucket()
	return bucket.OpenDownloadStream(ctx, fileID)
}

func (r *agentReleaseRepo) DeleteFile(ctx context.Context, fileID bson.ObjectID) error {
	bucket := r.Bucket()
	return bucket.Delete(ctx, fileID)
}

func (r *agentReleaseRepo) CreateIndex(ctx context.Context) error {
	idx := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "goos", Value: 1},
				{Key: "goarch", Value: 1},
				{Key: "semver", Value: 1},
			},
			Options: options.Index().SetUnique(true),
		},
	}
	_, err := r.Indexes().CreateMany(ctx, idx)

	return err
}

package mongodb

import (
	"errors"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/x/mongo/driver/connstring"
)

func Open(uri string, opts ...*options.ClientOptions) (*mongo.Database, error) {
	pu, err := connstring.ParseAndValidate(uri)
	if err != nil {
		return nil, err
	}
	dbname := pu.Database
	if dbname == "" {
		return nil, errors.New("database name is required")
	}

	opt := options.Client().ApplyURI(uri)
	opts = append(opts, opt)

	cli, err := mongo.Connect(opts...)
	if err != nil {
		return nil, err
	}

	db := cli.Database(dbname)

	return db, nil
}

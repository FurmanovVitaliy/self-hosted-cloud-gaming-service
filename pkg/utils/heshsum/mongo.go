package hashsum

import (
	"cloud/pkg/logger"
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

type Hash struct {
	ID      string `bson:"_id"`
	Hashsum string `bson:"hash"`
}

type db struct {
	collection *mongo.Collection
	logger     *logger.Logger
}

func HashStorage(database *mongo.Database, collection string, logger *logger.Logger) Storage {
	return &db{
		collection: database.Collection(collection),
		logger:     logger,
	}
}

func (d *db) Create(ctx context.Context, id string, hash string) error {
	_, err := d.collection.InsertOne(ctx, bson.M{"_id": id, "hash": hash})
	if err != nil {
		return fmt.Errorf("failed to create user due to error:%v", err)

	}
	return nil
}
func (d *db) FindOne(ctx context.Context, id string) (string, error) {
	var h Hash
	filter := bson.M{"_id": id}
	result := d.collection.FindOne(ctx, filter)
	if result.Err() != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			return "", fmt.Errorf("failed to find hashsum of %s due to error:%v", id, result.Err())
		}
		return "", fmt.Errorf("failed to find hashsum of %s due to error:%v", id, result.Err())
	}
	err := result.Decode(&h)
	if err != nil {
		return "", fmt.Errorf("failed to decode hashsum of %s due to error:%v", id, result.Err())
	}

	return h.Hashsum, nil
}
func (d *db) Update(ctx context.Context, id string, hash string) error {
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"hash": hash}}
	result := d.collection.FindOneAndUpdate(ctx, filter, update)
	if result.Err() != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			d.Create(ctx, id, hash)
		}
		return fmt.Errorf("failed to update hashsum of %s due to error:%v", id, result.Err())
	}
	return nil
}
func (d *db) Delete(ctx context.Context, id string) error {
	return nil
}

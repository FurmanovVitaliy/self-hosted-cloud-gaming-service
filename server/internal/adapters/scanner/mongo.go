package scanner

import (
	"context"

	"github.com/FurmanovVitaliy/pixel-cloud/pkg/logger"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

type scanDAO struct {
	collection *mongo.Collection
	logger     *logger.Logger
}

// NewStorage returns new storage
func NewStorage(database *mongo.Database, collection string, logger *logger.Logger) *scanDAO {
	return &scanDAO{
		collection: database.Collection(collection),
		logger:     logger,
	}
}

func (s *scanDAO) GetHash(ctx context.Context, id string) (string, error) {
	var result HashRecord
	filter := bson.M{"_id": id}
	err := s.collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return "", err
	}
	return result.Hash, nil
}

func (s *scanDAO) UpsertHash(ctx context.Context, id string, hash string) error {
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"hash": hash}}
	opts := options.Update().SetUpsert(true)
	_, err := s.collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return err
	}
	return nil
}

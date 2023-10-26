package user

import (
	"cloud/pkg/logger"
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

// ?? DB wichh implements Storage interface
// ?? We can easy change DB to another one (Postgres, Redis, etc)
type db struct {
	collection *mongo.Collection
	logger     *logger.Logger
}

// NewStorage returns new storage
func NewStorage(database *mongo.Database, collection string, logger *logger.Logger) Storage {
	return &db{
		collection: database.Collection(collection),
		logger:     logger,
	}
}

func (d *db) Update(ctx context.Context, u User) error {
	filter := bson.M{"_id": u.ID}
	update := bson.M{"$set": u}

	result, err := d.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to execute update user query due to error:%v", err)
	}
	if result.MatchedCount == 0 {
		//TODO ErrEntityNotFound
		return fmt.Errorf("not found")
	}
	d.logger.Tracef("Updated %d document(s) ", result.ModifiedCount)
	return nil
}

func (d *db) Create(ctx context.Context, u User) (string, error) {
	d.logger.Debug("create user")
	result, err := d.collection.InsertOne(ctx, u)
	if err != nil {
		return "", fmt.Errorf("failed to create user due to error:%v", err)
	}

	d.logger.Debug("convert InsertedID to ObjectID")
	oid, ok := result.InsertedID.(primitive.ObjectID)
	if ok {
		return oid.Hex(), nil
	}

	d.logger.Trace(u)
	return "", fmt.Errorf("failed to convert objectID to hex. oid:%s", oid)
}
func (s *db) FindByEmail(ctx context.Context, email string) (User, error) {
	var u User
	filter := bson.M{"_email": email}
	result := s.collection.FindOne(ctx, filter)
	if result.Err() != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			return u, fmt.Errorf("failed to find user by email:%s due to error:%v", email, result.Err())
		}
		return u, fmt.Errorf("failed to find user by email:%s due to error:%v", email, result.Err())
	}
	if err := result.Decode(&u); err != nil {
		return u, fmt.Errorf("failed to decode user from Db with email:%s", email)
	}
	return u, nil
}

func (d *db) FindAll(ctx context.Context) ([]User, error) {
	var users []User
	coursor, err := d.collection.Find(ctx, bson.M{})
	if coursor.Err() != nil {
		return users, fmt.Errorf("failed to find all users due to error:%v", err)
	}
	if err := coursor.All(ctx, &users); err != nil {
		return users, fmt.Errorf("failed to read all documents from coursor due to error:%v", err)
	}
	return users, nil
}

func (d *db) FindOne(ctx context.Context, id string) (User, error) {
	var u User
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return u, fmt.Errorf("failed to convert hex to ObjectID. hex:%s", id)
	}
	filter := bson.M{"_id": oid}
	result := d.collection.FindOne(ctx, filter)
	if result.Err() != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			return u, err //TODO:apperror.ErrNotFound ?MIDLEWARE
		}
		return u, fmt.Errorf("failed to find user by id:%s due to error:%v", id, err)
	}
	if err := result.Decode(&u); err != nil {
		return u, fmt.Errorf("failed to decode user from Db with id:%s", id)
	}
	return u, nil
}

func (d *db) Delete(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("failed to decode user from Db with id:%s", id)
	}

	filter := bson.M{"_id": oid}

	result, err := d.collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to execute delete user query due to error:%v", err)
	}
	if result.DeletedCount == 0 {
		//TODO ErrEntityNotFound
		return fmt.Errorf("not found")
	}
	d.logger.Tracef("Dleted %d document(s) ", result.DeletedCount)
	return nil
}

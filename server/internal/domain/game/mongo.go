package game

import (
	"context"
	"errors"
	"fmt"

	"github.com/FurmanovVitaliy/pixel-cloud/pkg/logger"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

type gameDAO struct {
	collection *mongo.Collection
	logger     *logger.Logger
}

func NewStorage(database *mongo.Database, collection string, logger *logger.Logger) *gameDAO {
	return &gameDAO{
		collection: database.Collection(collection),
		logger:     logger,
	}
}

func (d *gameDAO) Create(ctx context.Context, gt Game) (string, error) {
	result, err := d.collection.InsertOne(ctx, gt)
	if err != nil {
		return "", fmt.Errorf("failed to create game title due to error:%v", err)
	}

	oid, ok := result.InsertedID.(primitive.ObjectID)
	if ok {
		return oid.Hex(), nil
	}

	d.logger.Trace(gt)
	return "", fmt.Errorf("failed to convert objectID to hex. oid:%s", oid)
}
func (d *gameDAO) FindAll(ctx context.Context) (games []Game, err error) {
	coursor, err := d.collection.Find(ctx, bson.M{})
	if coursor.Err() != nil {
		return games, fmt.Errorf("failed to find all game titles due to error:%v", err)
	}
	if err := coursor.All(ctx, &games); err != nil {
		return games, fmt.Errorf("failed to read all documents from coursor due to error:%v", err)
	}
	return games, nil
}
func (d *gameDAO) FindOne(ctx context.Context, id string) (g Game, err error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return g, fmt.Errorf("failed to convert hex to ObjectID. hex:%s", id)
	}
	filter := bson.M{"_id": oid}
	result := d.collection.FindOne(ctx, filter)
	if result.Err() != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			return g, err
		}
		return g, fmt.Errorf("failed to find game title by id:%s due to error:%v", id, err)
	}
	if err := result.Decode(&g); err != nil {
		return g, fmt.Errorf("failed to decode game title from Db with id:%s", id)
	}
	return g, nil
}
func (d *gameDAO) FindOneByName(ctx context.Context, name string) (g Game, err error) {
	filter := bson.M{"name": name}
	result := d.collection.FindOne(ctx, filter)
	if result.Err() != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			return g, err
		}
		return g, fmt.Errorf("failed to find game title by name:%s due to error:%v", name, err)
	}
	if err := result.Decode(&g); err != nil {
		return g, fmt.Errorf("failed to decode game title from Db with name:%s", name)
	}
	return g, nil
}
func (d *gameDAO) Update(ctx context.Context, g Game) error {
	oid, err := primitive.ObjectIDFromHex(g.ID)
	if err != nil {
		return fmt.Errorf("failed to decode game title from Db with id:%s", g.ID)
	}
	filter := bson.M{"_id": oid}
	userBytes, err := bson.Marshal(g)
	if err != nil {
		return fmt.Errorf("failed to marshal game title due to err:%v", err)
	}
	var updateUserObject bson.M
	if err := bson.Unmarshal(userBytes, &updateUserObject); err != nil {
		return fmt.Errorf("failed to unmarshal game title bytes due to err:%v", err)
	}

	delete(updateUserObject, "_id")

	update := bson.M{
		"$set": updateUserObject,
	}

	result, err := d.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to execute update game title query due to error:%v", err)
	}
	if result.MatchedCount == 0 {
		return err
	}
	d.logger.Tracef("Matched %d document(s) and modified %d document(s)", result.MatchedCount, result.MatchedCount)
	return nil
}
func (d *gameDAO) FullyUpdate(ctx context.Context, games []Game) error {
	d.Drop(ctx)
	for _, g := range games {
		if _, err := d.Create(ctx, g); err != nil {
			return fmt.Errorf("failed to create game title due to error:%v", err)
		}
	}
	return nil
}
func (d *gameDAO) Delete(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("failed to decode game title from Db with id:%s", id)
	}

	filter := bson.M{"_id": oid}

	result, err := d.collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to execute delete game title query due to error:%v", err)
	}
	if result.DeletedCount == 0 {
		return fmt.Errorf("not found")
	}
	d.logger.Tracef("Dleted %d document(s) ", result.DeletedCount)
	return nil
}

func (d *gameDAO) Drop(ctx context.Context) error {
	if err := d.collection.Drop(ctx); err != nil {
		return fmt.Errorf("failed to drop game title collection due to error:%v", err)
	}
	return nil
}

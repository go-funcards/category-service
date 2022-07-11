package db

import (
	"context"
	"fmt"
	"github.com/go-funcards/category-service/internal/category"
	"github.com/go-funcards/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"time"
)

var _ category.Storage = (*storage)(nil)

const (
	timeout    = 5 * time.Second
	collection = "categories"
)

type storage struct {
	c mongodb.Collection[category.Category]
}

func NewStorage(ctx context.Context, db *mongo.Database, logger *zap.Logger) (*storage, error) {
	s := &storage{c: mongodb.Collection[category.Category]{
		Inner: db.Collection(collection),
		Log:   logger,
	}}

	if err := s.indexes(ctx); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *storage) indexes(ctx context.Context) error {
	name, err := s.c.Inner.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{"owner_id", 1},
			{"board_id", 1},
			{"position", 1},
			{"created_at", 1},
		},
	})
	if err == nil {
		s.c.Log.Info("index created", zap.String("collection", collection), zap.String("name", name))
	}

	return err
}

func (s *storage) Save(ctx context.Context, model category.Category) error {
	return s.SaveMany(ctx, []category.Category{model})
}

func (s *storage) SaveMany(ctx context.Context, models []category.Category) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var write []mongo.WriteModel
	for _, model := range models {
		data, err := s.c.ToM(model)
		if err != nil {
			return err
		}

		delete(data, "_id")
		delete(data, "owner_id")
		delete(data, "created_at")

		write = append(write, mongo.
			NewUpdateOneModel().
			SetUpsert(true).
			SetFilter(bson.M{"_id": model.CategoryID}).
			SetUpdate(bson.M{
				"$set": data,
				"$setOnInsert": bson.M{
					"owner_id":   model.OwnerID,
					"created_at": model.CreatedAt,
				},
			}),
		)
	}

	s.c.Log.Debug("bulk update")
	_, err := s.c.Inner.BulkWrite(ctx, write, options.BulkWrite())
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("bulk update: %s", mongodb.ErrMsgQuery), err)
	}
	return nil
}

func (s *storage) Delete(ctx context.Context, id string) error {
	return s.c.DeleteOne(ctx, bson.M{"_id": id})
}

func (s *storage) Find(ctx context.Context, filter category.Filter, index uint64, size uint32) ([]category.Category, error) {
	return s.c.Find(ctx, s.filter(filter), s.c.FindOptions(index, size).
		SetSort(bson.D{{"position", 1}, {"created_at", 1}}))
}

func (s *storage) Count(ctx context.Context, filter category.Filter) (uint64, error) {
	return s.c.CountDocuments(ctx, s.filter(filter))
}

func (s *storage) filter(filter category.Filter) bson.M {
	f := make(bson.M)
	if len(filter.CategoryIDs) > 0 {
		f["_id"] = bson.M{"$in": filter.CategoryIDs}
	}
	if len(filter.OwnerIDs) > 0 {
		f["owner_id"] = bson.M{"$in": filter.OwnerIDs}
	}
	if len(filter.BoardIDs) > 0 {
		f["board_id"] = bson.M{"$in": filter.BoardIDs}
	}
	return f
}

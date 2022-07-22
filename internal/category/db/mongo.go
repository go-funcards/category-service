package db

import (
	"context"
	"fmt"
	"github.com/go-funcards/category-service/internal/category"
	"github.com/go-funcards/mongodb"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

var _ category.Storage = (*storage)(nil)

const (
	timeout    = 5 * time.Second
	collection = "categories"
)

type storage struct {
	c   *mongo.Collection
	log logrus.FieldLogger
}

func NewStorage(ctx context.Context, db *mongo.Database, log logrus.FieldLogger) *storage {
	s := &storage{
		c:   db.Collection(collection),
		log: log,
	}
	s.indexes(ctx)
	return s
}

func (s *storage) indexes(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	name, err := s.c.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{"owner_id", 1},
			{"board_id", 1},
			{"position", 1},
			{"created_at", 1},
		},
	})
	if err != nil {
		s.log.WithFields(logrus.Fields{
			"collection": collection,
			"error":      err,
		}).Fatal("index not created")
	}

	s.log.WithFields(logrus.Fields{
		"collection": collection,
		"name":       name,
	}).Info("index created")
}

func (s *storage) Save(ctx context.Context, model category.Category) error {
	return s.SaveMany(ctx, []category.Category{model})
}

func (s *storage) SaveMany(ctx context.Context, models []category.Category) error {
	var write []mongo.WriteModel
	for _, model := range models {
		data, err := mongodb.ToBson(model)
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

	s.log.Info("categories save")

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	result, err := s.c.BulkWrite(ctx, write)
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("categories save: %s", mongodb.ErrMsgQuery), err)
	}

	s.log.WithFields(logrus.Fields{"result": result}).Info("categories saved")

	return nil
}

func (s *storage) Delete(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	s.log.WithField("category_id", id).Debug("category delete")
	result, err := s.c.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return fmt.Errorf(mongodb.ErrMsgQuery, err)
	}
	if result.DeletedCount == 0 {
		return fmt.Errorf(mongodb.ErrMsgQuery, mongo.ErrNoDocuments)
	}
	s.log.WithField("category_id", id).Debug("category deleted")

	return nil
}

func (s *storage) Find(ctx context.Context, filter category.Filter, index uint64, size uint32) ([]category.Category, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	opts := mongodb.FindOptions(index, size).SetSort(bson.D{{"position", 1}, {"created_at", 1}})
	cur, err := s.c.Find(ctx, s.build(filter), opts)
	if err != nil {
		return nil, fmt.Errorf(mongodb.ErrMsgQuery, err)
	}
	return mongodb.DecodeAll[category.Category](ctx, cur)
}

func (s *storage) Count(ctx context.Context, filter category.Filter) (uint64, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	total, err := s.c.CountDocuments(ctx, s.build(filter))
	if err != nil {
		return 0, fmt.Errorf(mongodb.ErrMsgQuery, err)
	}
	return uint64(total), nil
}

func (s *storage) build(filter category.Filter) any {
	f := make(mongodb.Filter, 0)
	if len(filter.CategoryIDs) > 0 {
		f = append(f, mongodb.In("_id", filter.CategoryIDs))
	}
	if len(filter.OwnerIDs) > 0 {
		f = append(f, mongodb.In("owner_id", filter.OwnerIDs))
	}
	if len(filter.BoardIDs) > 0 {
		f = append(f, mongodb.In("board_id", filter.BoardIDs))
	}
	return f.Build()
}

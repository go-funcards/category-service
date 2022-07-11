package category

import "context"

type Storage interface {
	Save(ctx context.Context, model Category) error
	SaveMany(ctx context.Context, models []Category) error
	Delete(ctx context.Context, id string) error
	Find(ctx context.Context, filter Filter, index uint64, size uint32) ([]Category, error)
	Count(ctx context.Context, filter Filter) (uint64, error)
}

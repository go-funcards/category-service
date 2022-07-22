package category

import (
	"context"
	"github.com/go-funcards/category-service/proto/v1"
	"github.com/go-funcards/slice"
	"google.golang.org/protobuf/types/known/emptypb"
)

var _ v1.CategoryServer = (*server)(nil)

type server struct {
	v1.UnimplementedCategoryServer
	storage Storage
}

func NewCategoryServer(storage Storage) *server {
	return &server{storage: storage}
}

func (s *server) CreateCategory(ctx context.Context, in *v1.CreateCategoryRequest) (*emptypb.Empty, error) {
	err := s.storage.Save(ctx, CreateCategory(in))

	return s.empty(err)
}

func (s *server) UpdateCategory(ctx context.Context, in *v1.UpdateCategoryRequest) (*emptypb.Empty, error) {
	err := s.storage.Save(ctx, UpdateCategory(in))

	return s.empty(err)
}

func (s *server) UpdateManyCategories(ctx context.Context, in *v1.UpdateManyCategoriesRequest) (*emptypb.Empty, error) {
	err := s.storage.SaveMany(ctx, slice.Map(in.GetCategories(), func(item *v1.UpdateCategoryRequest) Category {
		return UpdateCategory(item)
	}))

	return s.empty(err)
}

func (s *server) DeleteCategory(ctx context.Context, in *v1.DeleteCategoryRequest) (*emptypb.Empty, error) {
	err := s.storage.Delete(ctx, in.GetCategoryId())

	return s.empty(err)
}

func (s *server) GetCategories(ctx context.Context, in *v1.CategoriesRequest) (*v1.CategoriesResponse, error) {
	filter := CreateFilter(in)

	data, err := s.storage.Find(ctx, filter, in.GetPageIndex(), in.GetPageSize())
	if err != nil {
		return nil, err
	}

	total := uint64(len(data))
	if len(in.GetCategoryIds()) == 0 && uint64(in.GetPageSize()) == total {
		if total, err = s.storage.Count(ctx, filter); err != nil {
			return nil, err
		}
	}

	return &v1.CategoriesResponse{
		Categories: slice.Map(data, func(item Category) *v1.CategoriesResponse_Category {
			return item.toProto()
		}),
		Total: total,
	}, nil
}

func (s *server) empty(err error) (*emptypb.Empty, error) {
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

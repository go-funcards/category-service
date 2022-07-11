package category

import (
	"context"
	"github.com/go-funcards/category-service/proto/v1"
	"github.com/go-funcards/slice"
	"google.golang.org/protobuf/types/known/emptypb"
)

var _ v1.CategoryServer = (*categoryService)(nil)

type categoryService struct {
	v1.UnimplementedCategoryServer
	storage Storage
}

func NewCategoryService(storage Storage) *categoryService {
	return &categoryService{storage: storage}
}

func (s *categoryService) CreateCategory(ctx context.Context, in *v1.CreateCategoryRequest) (*emptypb.Empty, error) {
	if err := s.storage.Save(ctx, CreateCategory(in)); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *categoryService) UpdateCategory(ctx context.Context, in *v1.UpdateCategoryRequest) (*emptypb.Empty, error) {
	if err := s.storage.Save(ctx, UpdateCategory(in)); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *categoryService) UpdateManyCategories(ctx context.Context, in *v1.UpdateManyCategoriesRequest) (*emptypb.Empty, error) {
	if err := s.storage.SaveMany(ctx, slice.Map(in.GetCategories(), func(item *v1.UpdateCategoryRequest) Category {
		return UpdateCategory(item)
	})); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *categoryService) DeleteCategory(ctx context.Context, in *v1.DeleteCategoryRequest) (*emptypb.Empty, error) {
	if err := s.storage.Delete(ctx, in.GetCategoryId()); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *categoryService) GetCategories(ctx context.Context, in *v1.CategoriesRequest) (*v1.CategoriesResponse, error) {
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
			return item.toResponse()
		}),
		Total: total,
	}, nil
}
